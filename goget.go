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
	"os/exec"
	"path/filepath"
	"strings"
)

func callGoGet(goBinary, dir string, toUpdate []string) error {
	if len(toUpdate) == 0 {
		return nil
	}
	pkgs := strings.Join(toUpdate, " ")
	fmt.Printf("Updating %s with %s\n", dir, pkgs)
	cmd := exec.Command(goBinary, append([]string{"get", "-v"}, toUpdate...)...)
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("error getting absolute path to %s: %w", dir, err)
	}
	cmd.Dir = absDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running '%s get -v %s' in %s: %w", goBinary, pkgs, dir, err)
	}
	return nil
}
