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
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func findModFiles(excludes, only []string) ([]string, error) {
	checker, err := newChecker(excludes, only)
	if err != nil {
		return nil, err
	}
	var goModPaths sortPaths
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}
		if info.IsDir() {
			if checker.check(path) == ignore {
				return filepath.SkipDir
			}
			return nil
		}
		if checker.check(path) != allow {
			return nil
		}
		_, file := filepath.Split(path)
		if file == "go.mod" {
			goModPaths = append(goModPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Sort(&goModPaths)
	return goModPaths, nil
}

func newChecker(excludes, only []string) (pathChecker, error) {
	var cleanedExcludes map[string]struct{}
	var cleanedOnly map[string]bool
	if len(excludes) > 0 && len(only) > 0 {
		return pathChecker{}, fmt.Errorf("error finding mod files, can only provide either excludes or only paths")
	}
	if len(excludes) > 0 {
		cleanedExcludes = make(map[string]struct{}, len(excludes))
		for _, path := range excludes {
			cleanedExcludes[filepath.Clean(path)] = struct{}{}
		}
	}
	if len(only) > 0 {
		cleanedOnly = make(map[string]bool, len(only))
		for _, path := range only {
			var parts []string
			var last string
			for _, part := range strings.Split(filepath.Clean(path), string(filepath.Separator)) {
				parts = append(parts, part)
				partial := strings.Join(parts, string(filepath.Separator))
				partial = filepath.Clean(partial)
				cleanedOnly[partial] = false
				last = partial
			}
			cleanedOnly[last] = true
		}
		if _, ok := cleanedOnly["."]; !ok {
			cleanedOnly["."] = false
		}
	}
	return pathChecker{
		excludes: cleanedExcludes,
		only:     cleanedOnly,
	}, nil
}

type pathChecker struct {
	excludes map[string]struct{}
	only     map[string]bool
}

type checkResult int

const (
	ignore checkResult = iota
	ignoreFile
	allow
)

func (c pathChecker) check(path string) checkResult {
	if len(c.excludes) > 0 {
		if _, ok := c.excludes[filepath.Clean(path)]; ok {
			return ignore
		}
		return allow
	}
	if len(c.only) > 0 {
		checkPath := filepath.Clean(path)
		if top, ok := c.only[checkPath]; ok {
			if top {
				return allow
			}
			return ignoreFile
		}
		for {
			dir, _ := filepath.Split(checkPath)
			checkPath = filepath.Clean(dir)
			if top, ok := c.only[checkPath]; ok {
				if top {
					return allow
				}
				return ignore
			}
		}
	}
	return allow
}

type sortPaths []string

var _ sort.Interface = &sortPaths{}

func (s *sortPaths) Len() int {
	return len(*s)
}

func (s *sortPaths) Less(i, j int) bool {
	ci := strings.Count((*s)[i], string(filepath.Separator))
	cj := strings.Count((*s)[j], string(filepath.Separator))
	if ci == cj {
		return (*s)[i] < (*s)[j]
	}
	return ci < cj
}

func (s *sortPaths) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}
