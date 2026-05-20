package main

import (
	"log/slog"
	"os"
	"path/filepath"
)

// migrateLegacyStateDir relocates the pre-rebrand /var/lib/xray-stack
// directory to /var/lib/zeroone on first run. It acts when the new
// path is the active state dir AND the new directory is either absent
// or empty — so an install flow that pre-creates /var/lib/zeroone with
// `mkdir -p` still picks up legacy state. Any failure is logged and
// ignored so the daemon still starts.
func migrateLegacyStateDir(stateDir string) {
	if stateDir != "/var/lib/zeroone" {
		return
	}
	legacy := "/var/lib/xray-stack"
	legacyInfo, err := os.Stat(legacy)
	if err != nil || !legacyInfo.IsDir() {
		return
	}

	switch newEntries, err := os.ReadDir(stateDir); {
	case os.IsNotExist(err):
		// Target absent — rename the whole legacy tree.
		if err := os.Rename(legacy, stateDir); err != nil {
			slog.Warn("migrate legacy state dir", "from", legacy, "to", stateDir, "err", err)
			return
		}
		slog.Warn("migrated legacy state dir", "from", legacy, "to", stateDir)
	case err != nil:
		slog.Warn("migrate legacy state dir: stat target", "path", stateDir, "err", err)
	case len(newEntries) == 0:
		// Target exists but is empty (e.g. pre-created by an installer).
		// Move each legacy child into place, then drop the empty source.
		legacyEntries, lerr := os.ReadDir(legacy)
		if lerr != nil {
			slog.Warn("migrate legacy state dir: read source", "path", legacy, "err", lerr)
			return
		}
		moved := 0
		for _, e := range legacyEntries {
			src := filepath.Join(legacy, e.Name())
			dst := filepath.Join(stateDir, e.Name())
			if err := os.Rename(src, dst); err != nil {
				slog.Warn("migrate legacy state dir: move child", "src", src, "dst", dst, "err", err)
				continue
			}
			moved++
		}
		if moved > 0 {
			slog.Warn("migrated legacy state dir contents", "from", legacy, "to", stateDir, "entries", moved)
		}
		if err := os.Remove(legacy); err != nil && !os.IsNotExist(err) {
			slog.Warn("migrate legacy state dir: remove source", "path", legacy, "err", err)
		}
	}
	// Non-empty target: assume already migrated or operator-owned; skip.
}
