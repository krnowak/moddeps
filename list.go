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
	"sort"

	"github.com/spf13/cobra"
)

type listCmd struct {
	includeIndirect bool
	excludes        []string
	only            []string
	verbose         bool
}

func getListCommand() *cobra.Command {
	lc := &listCmd{}
	cmd := &cobra.Command{
		Use:  "list",
		RunE: lc.list,
	}
	f := cmd.Flags()
	f.BoolVarP(&lc.includeIndirect, "include-indirect", "i", false, "TODO")
	f.StringArrayVarP(&lc.excludes, "exclude", "e", nil, "TODO")
	f.StringArrayVarP(&lc.only, "only", "o", nil, "TODO")
	f.BoolVarP(&lc.verbose, "verbose", "v", false, "TODO")
	return cmd
}

func (c *listCmd) list(cmd *cobra.Command, args []string) error {
	paths, err := findModFiles(c.excludes, c.only)
	if err != nil {
		return err
	}
	modFiles, err := parseModFiles(paths)
	if err != nil {
		return err
	}
	pkgs := pkgVersionDirs(modFiles, c.includeIndirect)
	sortedPkgs := make([]string, 0, len(pkgs))
	for pkg := range pkgs {
		sortedPkgs = append(sortedPkgs, pkg)
	}
	sort.Strings(sortedPkgs)
	for _, pkg := range sortedPkgs {
		versionMap := pkgs[pkg]
		if c.verbose || len(versionMap) > 1 {
			fmt.Printf("%s\n", pkg)
			for v, dirs := range versionMap {
				fmt.Printf("- %s:\n", v)
				for _, dir := range dirs {
					fmt.Printf("  %s\n", dir)
				}
			}
		} else {
			// this will iterate just once
			for v := range versionMap {
				fmt.Printf("%s: %s\n", pkg, v)
			}
		}
	}
	return nil
}
