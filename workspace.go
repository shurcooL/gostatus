package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bradfitz/iter"
	vcs2 "github.com/shurcooL/go/vcs"
	"github.com/shurcooL/gostatus/pkg"
	"golang.org/x/tools/go/vcs"
)

// workspace is a Go workspace environment; each repo has local and remote components.
type workspace struct {
	shouldShow pkg.RepoFilter
	presenter  pkg.RepoStringer

	reposMu sync.Mutex
	repos   map[string]*pkg.Repo // Map key is repoRoot.

	in     chan string
	phase2 chan *pkg.Repo
	phase3 chan *pkg.Repo // Output is processed repos (complete with local and remote information), filtered with shouldShow.
	Out    chan string    // Out contains results of running presenter on processed repos.
}

func NewWorkspace(shouldShow pkg.RepoFilter, presenter pkg.RepoStringer) *workspace {
	u := &workspace{
		shouldShow: shouldShow,
		presenter:  presenter,

		repos: make(map[string]*pkg.Repo),

		in:     make(chan string, 64),
		phase2: make(chan *pkg.Repo, 64),
		phase3: make(chan *pkg.Repo, 64),
		Out:    make(chan string, 64),
	}

	var wg1, wg2, wg3 sync.WaitGroup

	for range iter.N(parallelism) {
		wg1.Add(1)
		go u.phase12Worker(&wg1)
	}
	go func() {
		wg1.Wait()
		close(u.phase2)
	}()

	for range iter.N(parallelism) {
		wg2.Add(1)
		go u.phase23Worker(&wg2)
	}
	go func() {
		wg2.Wait()
		close(u.phase3)
	}()

	for range iter.N(parallelism) {
		wg3.Add(1)
		go u.phase34Worker(&wg3)
	}
	go func() {
		wg3.Wait()
		close(u.Out)
	}()

	return u
}

// Add adds a package with specified import path for processing.
func (u *workspace) Add(importPath string) {
	u.in <- importPath
}

// Done should be called after the workspace is finished being populated.
func (u *workspace) Done() {
	close(u.in)
}

// worker for phase 1, sends unique repos to phase 2.
func (u *workspace) phase12Worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for importPath := range u.in {
		//started := time.Now()
		// Determine repo root and local revision.
		// This is potentially somewhat slow.
		bpkg, err := build.Import(importPath, wd, build.FindOnly)
		if err != nil {
			log.Println("build.Import:", err)
			continue
		}
		if bpkg.Goroot {
			continue
		}
		vcs2 := vcs2.New(bpkg.Dir)
		if vcs2 == nil {
			// TODO: Include for "????" output in gostatus.
			log.Println("not in VCS:", bpkg.Dir)
			continue
		}
		root := repoRoot(vcs2.RootPath(), bpkg.SrcRoot)
		//fmt.Printf("build + vcs: %v ms.\n", time.Since(started).Seconds()*1000)

		var repo *pkg.Repo
		u.reposMu.Lock()
		if _, ok := u.repos[root]; !ok {
			repo = &pkg.Repo{
				Root: root,
				VCS:  vcs2,
				// TODO: Maybe keep track of import paths inside, etc.
			}
			u.repos[root] = repo
		} else {
			// TODO: Maybe keep track of import paths inside, etc.
		}
		u.reposMu.Unlock()

		// If new repo, send off to phase 2 channel.
		if repo != nil {
			u.phase2 <- repo
		}
	}
}

// Phase 2 to 3 figures out repo local and remote information.
func (u *workspace) phase23Worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range u.phase2 {
		//started := time.Now()
		// Determine remote revision.
		// This is slow because it requires a network operation.
		var remoteVCS vcs2.Remote = p.VCS
		var localVCS vcs2.Vcs = p.VCS
		remoteRevision := remoteVCS.GetRemoteRev()
		//fmt.Printf("remoteVCS.GetRemoteRev: %v ms.\n", time.Since(started).Seconds()*1000)

		p.Remote.Revision = remoteRevision

		if rr, err := vcs.RepoRootForImportPath(p.Root, false); err == nil {
			p.Remote.RepoURL = rr.Repo
		}

		// TODO: Organize.
		p.Local.Revision = localVCS.GetLocalRev()

		// TODO: Organize.
		p.Local.RemoteURL = localVCS.GetRemote()

		// TODO: Organize.
		if remoteRevision != "" {
			p.Remote.IsContained = localVCS.IsContained(remoteRevision)
		}

		// TODO: Organize and maybe do at a later stage, after checking shouldShow?
		//       Actually, probably need this for shouldShow, etc.
		p.Local.Status = localVCS.GetStatus()
		p.Local.Stash = localVCS.GetStash()
		p.Local.LocalBranch = localVCS.GetLocalBranch()

		if !u.shouldShow(p) {
			continue
		}

		u.phase3 <- p
	}
}

// Phase 3 to 4 ...
func (u *workspace) phase34Worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for repo := range u.phase3 {
		u.Out <- u.presenter(repo)
	}
}

// repoRoot figures out the repo root import path given repoPath and srcRoot.
// It handles symlinks that may be involved in the paths.
// It also handles a possible case mismatch in the prefix, printing a warning to stderr if detected.
func repoRoot(repoPath, srcRoot string) string {
	if s, err := filepath.EvalSymlinks(repoPath); err == nil {
		repoPath = s
	} else {
		fmt.Fprintln(os.Stderr, "warning: repoRoot: can't resolve symlink:", err)
	}
	if s, err := filepath.EvalSymlinks(srcRoot); err == nil {
		srcRoot = s
	} else {
		fmt.Fprintln(os.Stderr, "warning: repoRoot: can't resolve symlink:", err)
	}

	sep := string(filepath.Separator)

	// Detect and handle case mismatch in prefix.
	if prefixLen := len(srcRoot + sep); len(repoPath) >= prefixLen && srcRoot+sep != repoPath[:prefixLen] && strings.EqualFold(srcRoot+sep, repoPath[:prefixLen]) {
		fmt.Fprintln(os.Stderr, "warning: repoRoot: prefix case doesn't match:", srcRoot+sep, repoPath[:prefixLen])
		return filepath.ToSlash(repoPath[prefixLen:])
	}

	return filepath.ToSlash(strings.TrimPrefix(repoPath, srcRoot+sep))
}
