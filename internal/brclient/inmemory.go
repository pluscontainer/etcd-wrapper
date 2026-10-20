// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package brclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gardener/etcd-wrapper/internal/types"
	"github.com/gardener/etcd-wrapper/internal/util"
)

// TriggerDeltaSnapshot requests an out-of-schedule delta snapshot from backup-restore.
func TriggerDeltaSnapshot(ctx context.Context, brConfig types.BackupRestoreConfig) error {
	client, err := createClient(brConfig)
	if err != nil {
		return err
	}
	c := &brClient{client: client, backupRestoreBaseAddress: brConfig.GetBaseAddress()}

	response, err := c.createAndExecuteHTTPRequest(ctx, http.MethodGet, c.backupRestoreBaseAddress+"/snapshot/delta")
	if err != nil {
		return err
	}
	defer util.CloseResponseBody(response)

	if !util.ResponseHasOKCode(response) {
		return fmt.Errorf("server returned error response code when attempting to trigger a delta snapshot: %s", response.Status)
	}
	return nil
}
