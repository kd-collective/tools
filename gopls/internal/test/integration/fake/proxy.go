// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fake

import (
	"fmt"
	"strings"

	"golang.org/x/tools/internal/proxydir"
)

// WriteProxy creates a new proxy file tree using the txtar-encoded content,
// and returns its URL.
func WriteProxy(tmpdir string, files map[string][]byte) (string, error) {
	type moduleVersion struct {
		modulePath, version string
	}
	// Transform into the format expected by the proxydir package.
	filesByModule := make(map[moduleVersion]map[string][]byte)
	for name, data := range files {
		modulePath, version, suffix := splitModuleVersionPath(name)
		mv := moduleVersion{modulePath, version}
		if _, ok := filesByModule[mv]; !ok {
			filesByModule[mv] = make(map[string][]byte)
		}
		filesByModule[mv][suffix] = data
	}
	for mv, files := range filesByModule {
		// Don't hoist this check out of the loop:
		// the problem is benign if filesByModule is empty.
		if strings.Contains(tmpdir, "#") {
			return "", fmt.Errorf("WriteProxy's tmpdir contains '#', which is unsuitable for GOPROXY. (If tmpdir was derived from testing.T.Name, use t.Run to ensure that each subtest has a unique name.)")
		}
		if err := proxydir.WriteModuleVersion(tmpdir, mv.modulePath, mv.version, files); err != nil {
			return "", fmt.Errorf("error writing %s@%s: %v", mv.modulePath, mv.version, err)
		}
	}
	return proxydir.ToURL(tmpdir), nil
}
