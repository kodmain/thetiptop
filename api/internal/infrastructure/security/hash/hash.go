package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
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
func Hash(data *string, algo HashAlgo) (*string, errors.ErrorInterface) {
	var hashedData []byte

	if data == nil {
		return nil, errors.ErrNoData
	}

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
		return nil, errors.ErrInternalServer
	}

	hashed := hex.EncodeToString(hashedData)

	return &hashed, nil
}

func CompareHash(hashedData, data *string, algo HashAlgo) errors.ErrorInterface {
	if data == nil || hashedData == nil {
		return errors.ErrNoData
	}

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
		return compareHashBcrypt(hashedData, data)
	default:
		return errors.ErrHashAlgoUnknown
	}
}

func compareHashBcrypt(hashedData, data *string) errors.ErrorInterface {
	if bcrypt.CompareHashAndPassword([]byte(*hashedData), []byte(*data)) != nil {
		return errors.ErrUnauthorized
	}

	return nil
}

func compareHash(hashedData, data *string, h hash.Hash) errors.ErrorInterface {
	if hex.EncodeToString(hashWithAlgo(h, data)) != *hashedData {
		return errors.ErrUnauthorized
	}

	return nil
}

// hashWithAlgo hache les données avec un algorithme spécifique
func hashWithAlgo(h hash.Hash, data *string) []byte {
	h.Write([]byte(*data))
	return h.Sum(nil)
}

// hashWithBcrypt utilise bcrypt pour hacher les données
func hashWithBcrypt(data *string) (*string, errors.ErrorInterface) {
	hashedData, err := bcrypt.GenerateFromPassword([]byte(*data), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	hashed := string(hashedData)

	return &hashed, nil
}
