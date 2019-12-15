// Copyright 2014 The go-pbfcoin Authors
// This file is part of the go-pbfcoin library.
//
// The go-pbfcoin library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-pbfcoin library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-pbfcoin library. If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// MakeName creates a node name that follows the pbfcoin convention
// for such names. It adds the operation system name and Go runtime version
// the name.
func MakeName(name, version string) string {
	return fmt.Sprintf("%s/v%s/%s/%s", name, version, runtime.GOOS, runtime.Version())
}

func ExpandHomePath(p string) (path string) {
	path = p
	sep := fmt.Sprintf("%s", os.PathSeparator)

	// Check in case of paths like "/sompbfing/~/sompbfing/"
	if len(p) > 1 && p[:1+len(sep)] == "~"+sep {
		usr, _ := user.Current()
		dir := usr.HomeDir

		path = strings.Replace(p, "~", dir, 1)
	}

	return
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func AbsolutePath(Datadir string, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(Datadir, filename)
}

func HomeDir() (home string) {
	if usr, err := user.Current(); err == nil {
		home = usr.HomeDir
	} else {
		home = os.Getenv("HOME")
	}
	return
}

func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := HomeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "pbfcoin")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "pbfcoin")
		} else {
			return filepath.Join(home, ".pbfcoin")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func DefaultIpcPath() string {
	if runtime.GOOS == "windows" {
		return `\\.\pipe\gpbf.ipc`
	}
	return filepath.Join(DefaultDataDir(), "gpbf.ipc")
}
