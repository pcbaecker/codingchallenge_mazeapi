package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	ARGON2_PARAMETER_TIME    = 1
	ARGON2_PARAMETER_MEMORY  = 64 * 1024
	ARGON2_PARAMETER_THREADS = 4
	ARGON2_PARAMETER_KEYLEN  = 32
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func hashPassword(password string) (string, error) {
	salt, err := generateRandomBytes(16)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, ARGON2_PARAMETER_TIME, ARGON2_PARAMETER_MEMORY, ARGON2_PARAMETER_THREADS, ARGON2_PARAMETER_KEYLEN)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, ARGON2_PARAMETER_MEMORY, ARGON2_PARAMETER_TIME, ARGON2_PARAMETER_THREADS, b64Salt, b64Hash), nil
}

func comparePasswordAndHash(password string, encodedHash string) (bool, error) {
	memory, iterations, threads, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, iterations, memory, uint8(threads), ARGON2_PARAMETER_KEYLEN)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func decodeHash(encodedHash string) (memory uint32, iterations uint32, threads uint32, salt []byte, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return 0, 0, 0, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}
	if version != argon2.Version {
		return 0, 0, 0, nil, nil, ErrIncompatibleVersion
	}

	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &threads)
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}

	return memory, iterations, threads, salt, hash, nil
}
