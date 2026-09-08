package models

import (
	"crypto/md5"
	"encoding/hex"
)

// CalculateContentHash computes an MD5 hash of the given content.
// This hash is used to validate cached embeddings against current content
// and detect when posts have been modified.
//
// The function returns a hex-encoded string representation of the MD5 hash.
//
// Requirements:
//   - 3.3: Content hash validation for cache invalidation
//   - Design: Content hash calculation algorithm
func CalculateContentHash(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}
