// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package brclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gardener/etcd-wrapper/internal/types"
	. "github.com/onsi/gomega"
)

func TestTriggerDeltaSnapshot(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		expectErr bool
	}{
		{"succeeds on 200", http.StatusOK, false},
		{"fails on 500", http.StatusInternalServerError, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			var path string
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				path = req.URL.Path
				rw.WriteHeader(tc.status)
			}))
			defer server.Close()

			err := TriggerDeltaSnapshot(context.Background(), types.BackupRestoreConfig{HostPort: strings.TrimPrefix(server.URL, "http://")})

			g.Expect(path).To(Equal("/snapshot/delta"))
			if tc.expectErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).ToNot(HaveOccurred())
			}
		})
	}
}
