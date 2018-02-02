//
// Copyright © 2018-present Keith Irwin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published
// by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func walkTree(root string) error {

	walker := func(path string, info os.FileInfo, err error) error {
		timeStr := info.ModTime().Format("2006-01-02 15:04:05")
		fmt.Printf("%v - %v\n", timeStr, path)
		return nil
	}

	return filepath.Walk(root, walker)
}

func main() {
	path := "."
	args := os.Args
	if len(args) > 1 {
		path = args[1]
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	walkTree(path)
	os.Exit(0)
}
