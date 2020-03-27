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
	"path/filepath"
)

// pkg name -> version -> go mod dirs
func pkgVersionDirs(modFiles []modFile, includeIndirect bool) map[string]map[string][]string {
	pkgs := make(map[string]map[string][]string)
	for _, mf := range modFiles {
		for _, req := range mf.parsed.Require {
			if req.Indirect && !includeIndirect {
				continue
			}
			versions, ok := pkgs[req.Mod.Path]
			if !ok {
				versions = make(map[string][]string)
				pkgs[req.Mod.Path] = versions
			}
			dir, _ := filepath.Split(mf.path)
			s := versions[req.Mod.Version]
			s = append(s, filepath.Clean(dir))
			versions[req.Mod.Version] = s
		}
	}
	return pkgs
}
