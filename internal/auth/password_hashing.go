package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/crypto/argon2"
)

const megabyte = 1024

var (
	errInvalidHash         = errors.New("argon2id: hash is not in the correct format")
	errIncompatibleVariant = errors.New("argon2id: incompatible variant of argon2")
	errIncompatibleVersion = errors.New("argon2id: incompatible version of argon2")
)

type argonParams struct {
	saltLength int
	memory     uint32
	time       uint32
	keyLength  uint32
	threads    uint8
}

var cfg = &argonParams{
	saltLength: 16,
	memory:     64 * megabyte,
	time:       1,
	keyLength:  32,
	threads:    uint8(runtime.NumCPU()),
}

func MakePasswordHash(password string) (hash string, err error) {
	salt, err := makeRandomSalt(cfg.saltLength)
	if err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(password), salt, cfg.time, cfg.memory, cfg.threads, cfg.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,t=%d$%s$%s", argon2.Version, cfg.memory, cfg.time, cfg.threads, b64Salt, b64Key)
	return hash, nil
}

func ComparePasswordAndHash(password, hash string) (match bool, err error) {
	params, salt, key, err := decodeHash(hash)
	if err != nil {
		return false, err
	}

	compareKey := argon2.IDKey([]byte(password), salt, params.time, params.memory, params.threads, params.keyLength)

	keyLen := int32(len(key))
	compareKeyLen := int32(len(compareKey))

	if subtle.ConstantTimeEq(keyLen, compareKeyLen) == 0 {
		return false, nil
	}

	if subtle.ConstantTimeCompare(key, compareKey) == 1 {
		return true, nil
	}

	return false, nil
}

func decodeHash(hash string) (params *argonParams, salt, key []byte, err error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errInvalidHash
	}

	if vals[1] != "argon2id" {
		return nil, nil, nil, errIncompatibleVariant
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errIncompatibleVersion
	}

	params = &argonParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,t=%d", &params.memory, &params.time, &params.threads)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.saltLength = len(salt)

	key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.keyLength = uint32(len(key))

	return params, salt, key, nil
}

func makeRandomSalt(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
