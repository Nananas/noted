package security

import (
	"crypto/sha256"
	"encoding/hex"
)

var _SALT string

func SetSalt(s string) {
	_SALT = s
}

func SaltHash(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s + _SALT))
	return hex.EncodeToString(hasher.Sum([]byte(nil)))
}
