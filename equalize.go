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

	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type equalizeCmd struct {
	includeIndirect bool
	excludes        []string
	only            []string
	goBinary        string
}

func getEqualizeCommand() *cobra.Command {
	ec := &equalizeCmd{}
	cmd := &cobra.Command{
		Use:  "equalize",
		RunE: ec.equalize,
		Args: cobra.ArbitraryArgs,
	}
	f := cmd.Flags()
	f.BoolVarP(&ec.includeIndirect, "include-indirect", "i", false, "TODO")
	f.StringArrayVarP(&ec.excludes, "exclude", "e", nil, "TODO")
	f.StringArrayVarP(&ec.only, "only", "o", nil, "TODO")
	f.StringVarP(&ec.goBinary, "go-binary", "g", "", "TODO")
	return cmd
}

func (c *equalizeCmd) equalize(cmd *cobra.Command, args []string) error {
	if c.goBinary == "" {
		c.goBinary = "go"
	}
	paths, err := findModFiles(c.excludes, c.only)
	if err != nil {
		return err
	}
	modFiles, err := parseModFiles(paths)
	if err != nil {
		return err
	}
	modVersions, err := argsToModuleMap(args)
	if err != nil {
		return err
	}
	pkgs := pkgVersionDirs(modFiles, c.includeIndirect)
	dirModVersionsToUpdate := make(map[string]map[string]string)
	for pkg, versionDirs := range pkgs {
		if len(versionDirs) > 1 {
			// invalid version, will always
			// compare as lower
			maxVersion := ""
			// find max version to equalize to
			for version := range versionDirs {
				if semver.Compare(maxVersion, version) < 0 {
					maxVersion = version
				}
			}
			for version, dirs := range versionDirs {
				if semver.Compare(version, maxVersion) >= 0 {
					continue
				}
				for _, dir := range dirs {
					dirModVersions, ok := dirModVersionsToUpdate[dir]
					if !ok {
						dirModVersions = make(map[string]string)
						dirModVersionsToUpdate[dir] = dirModVersions
					}
					dirModVersions[pkg] = maxVersion
				}
			}
		} else {
			// all modules have the same version
			// of the module, so try updating it
			//
			// this will iterate once
			for _, dirs := range versionDirs {
				for _, dir := range dirs {
					dirModVersions, ok := dirModVersionsToUpdate[dir]
					if !ok {
						dirModVersions = make(map[string]string)
						dirModVersionsToUpdate[dir] = dirModVersions
					}
					// we want to update,
					// hence empty version
					dirModVersions[pkg] = ""
				}
			}
		}
	}
	if len(modVersions) > 0 {
		for dir, dirModVersions := range dirModVersionsToUpdate {
			for pkg := range dirModVersions {
				if versionOverride, ok := modVersions[pkg]; ok {
					if versionOverride != "" {
						dirModVersions[pkg] = versionOverride
					}
				} else {
					delete(dirModVersions, pkg)
				}
			}
			if len(dirModVersions) == 0 {
				delete(dirModVersionsToUpdate, dir)
			}
		}
		if len(dirModVersionsToUpdate) == 0 {
			dirModVersionsToUpdate = nil
		}
	} else {
		for dir, dirModVersions := range dirModVersionsToUpdate {
			for pkg, version := range dirModVersions {
				if version == "" {
					delete(dirModVersions, pkg)
				}
			}
			if len(dirModVersions) == 0 {
				delete(dirModVersionsToUpdate, dir)
			}
		}
		if len(dirModVersionsToUpdate) == 0 {
			dirModVersionsToUpdate = nil
		}
	}
	for dir, dirModVersions := range dirModVersionsToUpdate {
		var toUpdate []string
		for pkg, version := range dirModVersions {
			if version == "" {
				toUpdate = append(toUpdate, pkg)
			} else {
				toUpdate = append(toUpdate, fmt.Sprintf("%s@%s", pkg, version))
			}
		}
		if err := c.callGoGet(dir, toUpdate); err != nil {
			return err
		}
	}
	return nil
}

func (c *equalizeCmd) callGoGet(dir string, toUpdate []string) error {
	return callGoGet(c.goBinary, dir, toUpdate)
}
