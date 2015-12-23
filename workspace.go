package main

import (
	"go/build"
	"log"
	"sync"

	"github.com/bradfitz/iter"
	vcs2 "github.com/shurcooL/go/vcs"
	"github.com/shurcooL/gostatus/pkg"
)

// goWorkspace is a workspace environment, meaning each repo has local and remote components.
type goWorkspace struct {
	shouldShow func(*pkg.Repo) bool
	presenter  pkg.RepoStringer

	reposMu sync.Mutex
	repos   map[string]*pkg.Repo // Map key is repoRoot.

	in  chan string
	wg1 sync.WaitGroup

	phase2 chan *pkg.Repo
	wg2    sync.WaitGroup

	phase3 chan *pkg.Repo
	wg3    sync.WaitGroup

	// Out is the output of processed repos (complete with local and remote revisions).
	Out chan string
}

func NewGoWorkspace(shouldShow func(*pkg.Repo) bool, presenter pkg.RepoStringer) *goWorkspace {
	u := &goWorkspace{
		shouldShow: shouldShow,
		presenter:  presenter,

		repos:  make(map[string]*pkg.Repo),
		in:     make(chan string, 64),
		phase2: make(chan *pkg.Repo, 64),
		phase3: make(chan *pkg.Repo, 64),
		Out:    make(chan string, 64),
	}

	for range iter.N(numWorkers) {
		u.wg1.Add(1)
		go u.phase12Worker()
	}
	go func() {
		u.wg1.Wait()
		close(u.phase2)
	}()

	for range iter.N(numWorkers) {
		u.wg2.Add(1)
		go u.phase23Worker()
	}
	go func() {
		u.wg2.Wait()
		close(u.phase3)
	}()

	for range iter.N(numWorkers) {
		u.wg3.Add(1)
		go u.phase34Worker()
	}
	go func() {
		u.wg3.Wait()
		close(u.Out)
	}()

	return u
}

// Add adds a package with specified import path for processing.
func (u *goWorkspace) Add(importPath string) {
	u.in <- importPath
}

// Done should be called after the workspace is finished being populated.
func (u *goWorkspace) Done() {
	close(u.in)
}

// worker for phase 1, sends unique repos to phase 2.
func (u *goWorkspace) phase12Worker() {
	defer u.wg1.Done()
	for importPath := range u.in {
		//started := time.Now()
		// Determine repo root and local revision.
		// This is potentially somewhat slow.
		bpkg, err := build.Import(importPath, wd, build.FindOnly)
		if err != nil {
			log.Println("build.Import:", err)
			continue
		}
		//goon.DumpExpr(bpkg)
		if bpkg.Goroot {
			continue
		}
		vcs2 := vcs2.New(bpkg.Dir)
		if vcs2 == nil {
			// TODO: Include for "????" output in gostatus.
			log.Println("not in VCS:", bpkg.Dir)
			continue
		}
		repoRoot := vcs2.RootPath()[len(bpkg.SrcRoot)+1:] // TODO: Consider sym links, etc.
		//fmt.Printf("build + vcs: %v ms.\n", time.Since(started).Seconds()*1000)

		var repo *pkg.Repo
		u.reposMu.Lock()
		if _, ok := u.repos[repoRoot]; !ok {
			repo = &pkg.Repo{
				Root: repoRoot,
				VCS:  vcs2,
				// TODO: Maybe keep track of import paths inside, etc.
			}
			u.repos[repoRoot] = repo
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
func (u *goWorkspace) phase23Worker() {
	defer u.wg2.Done()
	for p := range u.phase2 {
		//started := time.Now()
		// Determine remote revision.
		// This is slow because it requires a network operation.
		var remoteVCS vcs2.Remote = p.VCS
		var localVCS vcs2.Vcs = p.VCS
		remoteRevision := remoteVCS.GetRemoteRev()
		//fmt.Printf("remoteVCS.GetRemoteRev: %v ms.\n", time.Since(started).Seconds()*1000)

		p.Remote = pkg.Remote{
			Revision: remoteRevision,
		}

		// TODO: Organize.
		p.Local = pkg.Local{
			Revision: localVCS.GetLocalRev(),
		}

		// TODO: Organize.
		p.RemoteURL = localVCS.GetRemote()

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
func (u *goWorkspace) phase34Worker() {
	defer u.wg3.Done()
	for repo := range u.phase3 {
		u.Out <- u.presenter(repo)
	}
}
