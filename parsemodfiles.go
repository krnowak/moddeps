// Copyright Krzesimir Nowak
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/mod/modfile"
)

type modFile struct {
	parsed *modfile.File
	path   string
}

func parseModFiles(paths []string) ([]modFile, error) {
	if len(paths) == 0 {
		return nil, nil
	}
	modFiles := make([]modFile, 0, len(paths))
	for _, path := range paths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", path, err)
		}
		file, err := modfile.Parse(path, data, nil)
		if err != nil {
			return nil, fmt.Errorf("error parsing go mod file %s: %w", path, err)
		}
		for _, req := range file.Require {
			if req.Mod.Version == "" {
				return nil, fmt.Errorf("go mod file %s has a requirement on unversioned package %s, build the project first", path, req.Mod.Path)
			}
		}
		modFiles = append(modFiles, modFile{
			parsed: file,
			path:   path,
		})
	}
	return modFiles, nil
}
