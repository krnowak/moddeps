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
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortPaths(t *testing.T) {
	tcs := []struct{
		inPaths  []string
		outPaths []string
	}{
		{
			inPaths: []string{"a/c", "a/b"},
			outPaths: []string{"a/b", "a/c"},
		},
		{
			inPaths: []string{"a/c", "x"},
			outPaths: []string{"x", "a/c"},
		},
		{
			inPaths: []string{"a/c", "a/b", "a/b/c", "x"},
			outPaths: []string{"x", "a/b", "a/c", "a/b/c"},
		},
	}

	for _, tc := range tcs {
		inSort := make([]string, len(tc.inPaths))
		copy(inSort, tc.inPaths)
		sort.Sort((*sortPaths)(&inSort))
		require.Equal(t, tc.outPaths, inSort)
	}
}
