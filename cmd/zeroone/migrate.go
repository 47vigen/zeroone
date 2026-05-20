package main

import (
	"log/slog"
	"os"
)

// migrateLegacyStateDir relocates the pre-rebrand /var/lib/xray-stack
// directory to /var/lib/zeroone on first run. It only acts when the new
// path is the active state dir and the new directory does not yet
// exist; any failure is logged and ignored so the daemon still starts.
func migrateLegacyStateDir(stateDir string) {
	if stateDir != "/var/lib/zeroone" {
		return
	}
	if _, err := os.Stat(stateDir); err == nil {
		return
	}
	legacy := "/var/lib/xray-stack"
	if _, err := os.Stat(legacy); err != nil {
		return
	}
	if err := os.Rename(legacy, stateDir); err != nil {
		slog.Warn("migrate legacy state dir", "from", legacy, "to", stateDir, "err", err)
		return
	}
	slog.Warn("migrated legacy state dir", "from", legacy, "to", stateDir)
}
