// goiex: Golang interface to IEX Cloud API
// Copyright (C) 2019 Brian Hazeltine

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package helper

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
)

// RepoBaseDir contains the absolute path to the root of the repository on disk.
// This is determined by searching up the directory tree until a file called
// go.mod appears in the current directory.
var RepoBaseDir string

// RepoBasePkg contains the name of the module (eg github.com/myuser/someproj)
var RepoBasePkg string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("error getting runtime information")
	}

	RepoBaseDir = func(fpath string) string {
		for fpath != "." {
			fpath = filepath.Dir(fpath)
			if _, err := os.Stat(filepath.Join(fpath, "go.mod")); err == nil {
				return fpath
			}
		}
		panic("could not find go.mod")
	}(file)

	RepoBasePkg = func(fpath string) string {
		type foo struct{}
		pkgPath := reflect.TypeOf(foo{}).PkgPath()
		fpath = filepath.Dir(fpath)
		for path.Base(pkgPath) == filepath.Base(fpath) && filepath.Base(fpath) != filepath.Base(RepoBaseDir) {
			pkgPath = path.Dir(pkgPath)
			fpath = filepath.Dir(fpath)
		}
		return pkgPath
	}(file)
}
