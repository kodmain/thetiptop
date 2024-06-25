package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	asciiZero = 48
	asciiTen  = 57
)

// Validate returns an error if the provided string does not pass the luhn check.
func Validate(number string) error {
	p := len(number) % 2
	sum, err := calculateLuhnSum(number, p)
	if err != nil {
		return err
	}

	// If the total modulo 10 is not equal to 0, then the number is invalid.
	if sum%10 != 0 {
		return errors.New("invalid number")
	}

	return nil
}

// Calculate returns luhn check digit and the provided string number with its luhn check digit appended.
func Calculate(number string) (string, string, error) {
	p := (len(number) + 1) % 2
	sum, err := calculateLuhnSum(number, p)
	if err != nil {
		return "", "", nil
	}

	luhn := sum % 10
	if luhn != 0 {
		luhn = 10 - luhn
	}

	// If the total modulo 10 is not equal to 0, then the number is invalid.
	return strconv.FormatInt(luhn, 10), fmt.Sprintf("%s%d", number, luhn), nil
}

// Generate will generate a valid luhn number of the provided length
func Generate(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	var s strings.Builder
	for i := 0; i < length-1; i++ {
		s.WriteString(strconv.Itoa(rand.Intn(9)))
	}

	_, res, _ := Calculate(s.String()) //ignore error because this will always be valid
	return res
}

// GenerateWithExpiry will generate a valid luhn number with an expiration date embedded
func GenerateWithExpiry(length int, expiry time.Time) string {
	rand.Seed(time.Now().UTC().UnixNano())

	expiryStr := expiry.Format("0601021504") // Format expiry date as YYMMDDHHMM
	if length <= len(expiryStr) {
		return ""
	}

	var s strings.Builder
	for i := 0; i < length-len(expiryStr)-1; i++ {
		s.WriteString(strconv.Itoa(rand.Intn(9)))
	}
	s.WriteString(expiryStr)

	_, res, _ := Calculate(s.String())
	return res
}

func calculateLuhnSum(number string, parity int) (int64, error) {
	var sum int64
	for i, d := range number {
		if d < asciiZero || d > asciiTen {
			return 0, errors.New("invalid digit")
		}

		d = d - asciiZero
		// Double the value of every second digit.
		if i%2 == parity {
			d *= 2
			// If the result of this doubling operation is greater than 9.
			if d > 9 {
				// The same final result can be found by subtracting 9 from that result.
				d -= 9
			}
		}

		// Take the sum of all the digits.
		sum += int64(d)
	}

	return sum, nil
}

// GenerateWithPrefix will generate a valid luhn number of the provided length with prefix
func GenerateWithPrefix(prefix string, length int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	var s strings.Builder
	s.WriteString(prefix)
	length -= len(prefix)
	for i := 0; i < length-1; i++ {
		s.WriteString(strconv.Itoa(rand.Intn(9)))
	}

	_, res, _ := Calculate(s.String())
	return res
}

func main3() {
	expiry := time.Now().Add(1 * time.Second)
	tokenWithExpiry := GenerateWithExpiry(18, expiry)
	fmt.Println(tokenWithExpiry, Validate(tokenWithExpiry))
	time.Sleep(2 * time.Second)
	fmt.Println(tokenWithExpiry, Validate(tokenWithExpiry))
}

// calculateSecondsSinceReference Calculer le nombre de secondes depuis le 1er janvier 2000
//
// Parameters:
// - t: time.Time L'objet de temps à convertir
//
// Returns:
// - int: Le nombre de secondes depuis le 1er janvier 2000
func calculateSecondsSinceReference(t time.Time) int {
	referenceDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	duration := t.Sub(referenceDate)
	return int(duration.Seconds())
}

// generateLuhnDigit Génère le chiffre de contrôle selon l'algorithme de Luhn
//
// Parameters:
// - digits: []int La liste des chiffres sans le chiffre de contrôle
//
// Returns:
// - int: Le chiffre de contrôle
func generateLuhnDigit(digits []int) int {
	sum := 0
	for i := len(digits) - 1; i >= 0; i-- {
		digit := digits[i]
		if (len(digits)-i)%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return (10 - (sum % 10)) % 10
}

func main() {
	// Exemple d'utilisation avec la date et heure actuelles
	currentTime := time.Now().UTC()
	seconds := calculateSecondsSinceReference(currentTime)
	// Convertir les secondes en une séquence de chiffres
	digits := make([]int, 7)
	for i := 6; i >= 0; i-- {
		digits[i] = seconds % 10
		seconds /= 10
	}

	// Générer le chiffre de contrôle
	checkDigit := generateLuhnDigit(digits)
	// Ajouter le chiffre de contrôle à la séquence
	digits = append(digits, checkDigit)

	// Afficher le résultat
	fmt.Println("Séquence de chiffres avec contrôle Luhn:", digits)
}
