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
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type bumpCmd struct {
	includeIndirect bool
	excludes        []string
	only            []string
	goBinary        string
}

func getBumpCommand() *cobra.Command {
	bc := &bumpCmd{}
	cmd := &cobra.Command{
		Use:  "bump",
		RunE: bc.bump,
		Args: cobra.ArbitraryArgs,
	}
	f := cmd.Flags()
	f.BoolVarP(&bc.includeIndirect, "include-indirect", "i", false, "TODO")
	f.StringArrayVarP(&bc.excludes, "exclude", "e", nil, "TODO")
	f.StringArrayVarP(&bc.only, "only", "o", nil, "TODO")
	f.StringVarP(&bc.goBinary, "go-binary", "g", "", "TODO")
	return cmd
}

func (c *bumpCmd) bump(cmd *cobra.Command, args []string) error {
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
	for _, mf := range modFiles {
		var toUpdate []string
		for _, req := range mf.parsed.Require {
			if req.Indirect && !c.includeIndirect {
				continue
			}
			if version, ok := modVersions[req.Mod.Path]; ok {
				if version == "" {
					toUpdate = append(toUpdate, req.Mod.Path)
				} else if semver.Compare(version, req.Mod.Version) > 0 {
					toUpdate = append(toUpdate, fmt.Sprintf("%s@%s", req.Mod.Path, version))
				}
			} else if len(modVersions) == 0 {
				toUpdate = append(toUpdate, req.Mod.Path)
			}
		}
		dir, _ := filepath.Split(mf.path)
		if err := c.callGoGet(filepath.Clean(dir), toUpdate); err != nil {
			return err
		}
	}
	return nil
}

func (c *bumpCmd) callGoGet(dir string, toUpdate []string) error {
	return callGoGet(c.goBinary, dir, toUpdate)
}
