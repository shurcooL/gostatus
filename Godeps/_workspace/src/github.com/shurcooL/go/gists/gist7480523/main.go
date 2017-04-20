package gist7480523

import (
	"go/build"
	"strings"

	"github.com/shurcooL/go/exp/12"

	. "github.com/shurcooL/go/gists/gist5504644"
	. "github.com/shurcooL/go/gists/gist7519227"
	. "github.com/shurcooL/go/gists/gist7802150"
)

type GoPackageStringer func(*GoPackage) string

// A GoPackage describes a single package found in a directory.
// This is partially a copy of "cmd/go".Package, except it can be imported and reused. =.=
// https://code.google.com/p/go/source/browse/src/cmd/go/pkg.go?name=release#24
type GoPackage struct {
	Bpkg     *build.Package
	BpkgErr  error
	Standard bool // is this package part of the standard Go library?

	Dir *exp12.Directory
}

func GoPackageFromImportPathFound(importPathFound ImportPathFound) *GoPackage {
	bpkg, err := BuildPackageFromSrcDir(importPathFound.FullPath())
	return goPackageFromBuildPackage(bpkg, err)
}

func GoPackageFromImportPath(importPath string) *GoPackage {
	bpkg, err := BuildPackageFromImportPath(importPath)
	return goPackageFromBuildPackage(bpkg, err)
}

func goPackageFromBuildPackage(bpkg *build.Package, bpkgErr error) *GoPackage {
	if bpkgErr != nil {
		if _, noGo := bpkgErr.(*build.NoGoError); noGo || bpkg.Dir == "" {
			return nil
		}
	}

	if bpkg.ConflictDir != "" {
		return nil
	}

	goPackage := &GoPackage{
		Bpkg:     bpkg,
		BpkgErr:  bpkgErr,
		Standard: bpkg.Goroot && bpkg.ImportPath != "" && !strings.Contains(bpkg.ImportPath, "."), // https://code.google.com/p/go/source/browse/src/cmd/go/pkg.go?name=release#110

		Dir: exp12.LookupDirectory(bpkg.Dir),
	}

	/*if goPackage.Bpkg.Goroot == false { // Optimization that assume packages under Goroot are not under vcs
		// TODO: markAsNotNeedToUpdate() because of external insight?
	}*/

	return goPackage
}

// This is okay to call concurrently (a mutex is used internally).
// Actually, not completely okay because MakeUpdated technology is not thread-safe.
func (this *GoPackage) UpdateVcs() {
	if this.Bpkg.Goroot == false { // Optimization that assume packages under Goroot are not under vcs
		MakeUpdated(this.Dir)
	}
}

func (this *GoPackage) UpdateVcsFields() {
	if this.Dir.Repo != nil {
		MakeUpdated(this.Dir.Repo.VcsLocal)
		MakeUpdated(this.Dir.Repo.VcsRemote)
	}
}

func GetRepoImportPath(repoPath, srcRoot string) (repoImportPath string) {
	return strings.TrimPrefix(repoPath, srcRoot+"/")
}
func GetRepoImportPathPattern(repoPath, srcRoot string) (repoImportPathPattern string) {
	return GetRepoImportPath(repoPath, srcRoot) + "/..."
}

func (this *GoPackage) String() string {
	return this.Bpkg.ImportPath
}
