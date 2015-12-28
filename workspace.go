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
	ImportPaths       chan string // ImportPaths is the input for Go packages to be processed.
	unique            chan *Repo  // unique repos.
	processedFiltered chan *Repo  // processed repos, populated with local and remote info, filtered with shouldShow.
	Statuses          chan string // Statuses has results of running presenter on processed repos.

	shouldShow RepoFilter
	presenter  RepoPresenter

	reposMu sync.Mutex
	repos   map[string]*Repo // Map key is repoRoot.
}

func NewWorkspace(shouldShow RepoFilter, presenter RepoPresenter) *workspace {
	w := &workspace{
		ImportPaths:       make(chan string, 64),
		unique:            make(chan *Repo, 64),
		processedFiltered: make(chan *Repo, 64),
		Statuses:          make(chan string, 64),

		shouldShow: shouldShow,
		presenter:  presenter,

		repos: make(map[string]*Repo),
	}

	{
		var wg sync.WaitGroup
		for range iter.N(parallelism) {
			wg.Add(1)
			go w.uniqueWorker(&wg)
		}
		go func() {
			wg.Wait()
			close(w.unique)
		}()
	}
	{
		var wg sync.WaitGroup
		for range iter.N(parallelism) {
			wg.Add(1)
			go w.processFilterWorker(&wg)
		}
		go func() {
			wg.Wait()
			close(w.processedFiltered)
		}()
	}
	{
		var wg sync.WaitGroup
		for range iter.N(parallelism) {
			wg.Add(1)
			go w.presenterWorker(&wg)
		}
		go func() {
			wg.Wait()
			close(w.Statuses)
		}()
	}

	return w
}

// uniqueWorker finds unique repos out of all input Go packages.
func (w *workspace) uniqueWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	for importPath := range w.ImportPaths {
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
			w.unique <- repo
		}
	}
}

// processFilterWorker figures out repo local and remote info, and filters with shouldShow.
func (w *workspace) processFilterWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	for repo := range w.unique {
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

		w.processedFiltered <- repo
	}
}

// presenterWorker runs presenter on processed and filtered repos.
func (w *workspace) presenterWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	for repo := range w.processedFiltered {
		w.Statuses <- w.presenter(repo)
	}
}
