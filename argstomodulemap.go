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
	"strings"

	"golang.org/x/mod/module"
)

func argsToModuleMap(args []string) (map[string]string, error) {
	if len(args) == 0 {
		return nil, nil
	}
	modVersions := make(map[string]string, len(args))
	for _, modVersionString := range args {
		path := modVersionString
		version := ""
		if i := strings.Index(modVersionString, "@"); i >= 0 {
			path = modVersionString[:i]
			version = modVersionString[i+1:]
			if err := module.Check(path, version); err != nil {
				return nil, fmt.Errorf("invalid module path %s: %w", modVersionString, err)
			}
		} else {
			if err := module.CheckPath(path); err != nil {
				return nil, fmt.Errorf("invalid module path %s: %w", modVersionString, err)
			}
		}
		if _, ok := modVersions[path]; ok {
			return nil, fmt.Errorf("path %s specified twice", path)
		}
		modVersions[path] = version
	}
	return modVersions, nil
}
