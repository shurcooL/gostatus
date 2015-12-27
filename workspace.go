package main

import (
	"go/build"
	"log"
	"sync"

	"github.com/bradfitz/iter"
	"github.com/shurcooL/vcsstate"
	"golang.org/x/tools/go/vcs"
)

// workspace is a Go workspace environment; each repo has local and remote components.
type workspace struct {
	shouldShow RepoFilter
	presenter  RepoPresenter

	reposMu sync.Mutex
	repos   map[string]*Repo // Map key is repoRoot.

	in     chan string
	phase2 chan *Repo
	phase3 chan *Repo  // Output is processed repos (complete with local and remote information), filtered with shouldShow.
	Out    chan string // Out contains results of running presenter on processed repos.
}

func NewWorkspace(shouldShow RepoFilter, presenter RepoPresenter) *workspace {
	w := &workspace{
		shouldShow: shouldShow,
		presenter:  presenter,

		repos: make(map[string]*Repo),

		in:     make(chan string, 64),
		phase2: make(chan *Repo, 64),
		phase3: make(chan *Repo, 64),
		Out:    make(chan string, 64),
	}

	{
		var wg sync.WaitGroup
		for range iter.N(parallelism) {
			wg.Add(1)
			go w.phase12Worker(&wg)
		}
		go func() {
			wg.Wait()
			close(w.phase2)
		}()
	}
	{
		var wg sync.WaitGroup
		for range iter.N(parallelism) {
			wg.Add(1)
			go w.phase23Worker(&wg)
		}
		go func() {
			wg.Wait()
			close(w.phase3)
		}()
	}
	{
		var wg sync.WaitGroup
		for range iter.N(parallelism) {
			wg.Add(1)
			go w.phase34Worker(&wg)
		}
		go func() {
			wg.Wait()
			close(w.Out)
		}()
	}

	return w
}

// Add adds a package with specified import path for processing.
func (w *workspace) Add(importPath string) {
	w.in <- importPath
}

// Done should be called after the workspace is finished being populated.
func (w *workspace) Done() {
	close(w.in)
}

// worker for phase 1, sends unique repos to phase 2.
func (w *workspace) phase12Worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for importPath := range w.in {
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
		vcs, root, err := vcs.FromDir(bpkg.Dir, bpkg.SrcRoot)
		if err != nil {
			// TODO: Include for "????" output in gostatus.
			log.Println("not in VCS:", bpkg.Dir)
			continue
		}
		vcsstate, err := vcsstate.NewVCS(vcs)
		if err != nil {
			// TODO: Include for "????" output in gostatus.
			log.Println("repo not supported by vcsstate:", err)
			continue
		}

		var repo *Repo
		w.reposMu.Lock()
		if _, ok := w.repos[root]; !ok {
			repo = &Repo{
				Path: bpkg.Dir,
				Root: root,
				VCS:  vcsstate,
				// TODO: Maybe keep track of import paths inside, etc.
			}
			w.repos[root] = repo
		} else {
			// TODO: Maybe keep track of import paths inside, etc.
		}
		w.reposMu.Unlock()

		// If new repo, send off to phase 2 channel.
		if repo != nil {
			w.phase2 <- repo
		}
	}
}

// Phase 2 to 3 figures out repo local and remote information.
func (w *workspace) phase23Worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for repo := range w.phase2 {
		if s, err := repo.VCS.Status(repo.Path); err == nil {
			repo.Local.Status = s
		}
		if b, err := repo.VCS.Branch(repo.Path); err == nil {
			repo.Local.Branch = b
		}
		if r, err := repo.VCS.LocalRevision(repo.Path); err == nil {
			repo.Local.Revision = r
		}
		if s, err := repo.VCS.Stash(repo.Path); err == nil {
			repo.Local.Stash = s
		}
		if r, err := repo.VCS.RemoteURL(repo.Path); err == nil {
			repo.Local.RemoteURL = r
		}
		if r, err := repo.VCS.RemoteRevision(repo.Path); err == nil {
			repo.Remote.Revision = r
		}
		if repo.Remote.Revision != "" {
			if c, err := repo.VCS.Contains(repo.Path, repo.Remote.Revision); err == nil {
				repo.LocalContainsRemoteRevision = c
			}
		}
		if rr, err := vcs.RepoRootForImportPath(repo.Root, false); err == nil {
			repo.Remote.RepoURL = rr.Repo
		}

		if !w.shouldShow(repo) {
			continue
		}

		w.phase3 <- repo
	}
}

// Phase 3 to 4 ...
func (w *workspace) phase34Worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for repo := range w.phase3 {
		w.Out <- w.presenter(repo)
	}
}
