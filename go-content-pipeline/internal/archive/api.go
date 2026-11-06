package archive

import (
	"context"
	"os"
	"time"
)

// SaveJSON atomically writes report data. Repair workflows share the same
// storage implementation to avoid truncated files after interruption.
func SaveJSON(path string, value any) error { return saveJSON(path, value) }

func Lock(ctx context.Context, path string, timeout time.Duration) (func(), error) {
	return acquire(ctx, path, timeout)
}

// ReplaceLink applies the same exact-destination, backed-up replacement used by
// the archive verifier. Callers retain a manager for the lifetime of a run.
func ReplaceLink(o Options, e Entry, target string, backups *BackupManager) (string, error) {
	return replace(o, e, target, backups)
}

func (v *Verifier) VerifyURL(ctx context.Context, raw string, e Entry) Verification {
	return v.verify(ctx, raw, e, false)
}

// FindSnapshots returns successful captures newest first. Empty before searches
// all history through the present rather than requiring an existing archive.
func (v *Verifier) FindSnapshots(ctx context.Context, original, before string) ([]string, error) {
	if before == "" {
		before = time.Now().UTC().Add(time.Second).Format("20060102150405")
	}
	return v.snapshots(ctx, original, before)
}

func WriteFile(path string, data []byte, mode os.FileMode) error {
	return atomicWrite(path, data, mode)
}
