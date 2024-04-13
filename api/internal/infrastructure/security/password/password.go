package password

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	Lowercase    = 1 << iota // 1 (0b0001) pour les minuscules
	Uppercase                // 2 (0b0010) pour les majuscules
	Digits                   // 4 (0b0100) pour les chiffres
	SpecialChars             // 8 (0b1000) pour les caractères spéciaux

	All = Lowercase | Uppercase | Digits | SpecialChars
)

func GeneratePassword(length int, charsetOptions int) (password string, err error) {
	var charsetBuilder strings.Builder

	if charsetOptions&Lowercase > 0 {
		charsetBuilder.WriteString("abcdefghijklmnopqrstuvwxyz")
	}
	if charsetOptions&Uppercase > 0 {
		charsetBuilder.WriteString("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	if charsetOptions&Digits > 0 {
		charsetBuilder.WriteString("0123456789")
	}
	if charsetOptions&SpecialChars > 0 {
		charsetBuilder.WriteString("!@#$%^&*")
	}

	charset := charsetBuilder.String()

	if len(charset) == 0 {
		return "", err
	}

	for i := 0; i < length; i++ {
		charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password += string(charset[charIndex.Int64()])
	}

	return password, nil
}
