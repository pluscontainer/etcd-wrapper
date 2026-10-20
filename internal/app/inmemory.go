// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"time"

	"github.com/gardener/etcd-wrapper/internal/brclient"

	"go.uber.org/zap"
)

// shutdownSnapshotTimeout bounds the delta snapshot taken before etcd is stopped.
const shutdownSnapshotTimeout = 10 * time.Second

// takeShutdownSnapshot requests a delta snapshot from backup-restore if the application context was cancelled.
func (a *Application) takeShutdownSnapshot() {
	if a.ctx.Err() == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownSnapshotTimeout)
	defer cancel()

	a.logger.Info("taking delta snapshot before stopping etcd")
	if err := brclient.TriggerDeltaSnapshot(ctx, a.Config.BackupRestore); err != nil {
		a.logger.Error("failed to take delta snapshot before stopping etcd", zap.Error(err))
		return
	}
	a.logger.Info("took delta snapshot before stopping etcd")
}
