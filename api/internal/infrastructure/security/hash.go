package security

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"

	"golang.org/x/crypto/bcrypt"
)

// HashAlgo représente les algorithmes de hachage disponibles
type HashAlgo int

const (
	SHA1 HashAlgo = iota
	SHA256
	SHA512
	MD5
	BCRYPT
)

// hash crée un hachage du mot de passe en fonction de l'algorithme spécifié
func Hash(data string, algo HashAlgo) (string, error) {
	var hashedData []byte
	var err error

	switch algo {
	case SHA1:
		hashedData = hashWithAlgo(sha1.New(), data)
	case SHA256:
		hashedData = hashWithAlgo(sha256.New(), data)
	case SHA512:
		hashedData = hashWithAlgo(sha512.New(), data)
	case MD5:
		hashedData = hashWithAlgo(md5.New(), data)
	case BCRYPT:
		return hashWithBcrypt(data)
	default:
		return "", errors.New("unknown hash algorithm")
	}

	return hex.EncodeToString(hashedData), err
}

func CompareHash(hashedData, data string, algo HashAlgo) error {
	switch algo {
	case SHA1:
		return compareHash(hashedData, data, sha1.New())
	case SHA256:
		return compareHash(hashedData, data, sha256.New())
	case SHA512:
		return compareHash(hashedData, data, sha512.New())
	case MD5:
		return compareHash(hashedData, data, md5.New())
	case BCRYPT:
		return bcrypt.CompareHashAndPassword([]byte(hashedData), []byte(data))
	default:
		return errors.New("unknown hash algorithm")
	}
}

func compareHash(hashedData, data string, h hash.Hash) error {
	if hex.EncodeToString(hashWithAlgo(h, data)) != hashedData {
		return errors.New("hashes do not match")
	}

	return nil
}

// hashWithAlgo hache les données avec un algorithme spécifique
func hashWithAlgo(h hash.Hash, data string) []byte {
	h.Write([]byte(data))
	return h.Sum(nil)
}

// hashWithBcrypt utilise bcrypt pour hacher les données
func hashWithBcrypt(data string) (string, error) {
	hashedData, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	return string(hashedData), err
}
