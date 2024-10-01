package token

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

const (
	asciiZero = 48
	asciiTen  = 57
)

// Luhn représente une chaîne de caractères avec des méthodes associées pour les opérations de Luhn.
type Luhn string

func NewLuhn(s string) Luhn {
	return Luhn(s)
}

func (l Luhn) String() string {
	return string(l)
}

func (l Luhn) PointerString() *string {
	s := string(l)
	return &s
}

func (l Luhn) Pointer() *Luhn {
	return &l
}

// Validate retourne une erreur si la chaîne ne passe pas le test de Luhn.
func (l Luhn) Validate() errors.ErrorInterface {
	number := string(l)
	p := len(number) % 2
	sum, err := calculateLuhnSum(number, p)
	if err != nil {
		return err
	}

	// Si le total modulo 10 n'est pas égal à 0, alors le nombre est invalide.
	if sum%10 != 0 {
		return errors.ErrValueIsNotLuhn
	}

	return nil
}

// Calculate retourne le chiffre de contrôle Luhn et le nombre fourni avec son chiffre de contrôle Luhn ajouté.
func (l Luhn) Calculate() (string, string, error) {
	number := string(l)
	p := (len(number) + 1) % 2
	sum, err := calculateLuhnSum(number, p)
	if err != nil {
		return "", "", err
	}

	luhn := sum % 10
	if luhn != 0 {
		luhn = 10 - luhn
	}

	return strconv.FormatInt(luhn, 10), fmt.Sprintf("%s%d", number, luhn), nil
}

// Generate génère un numéro valide Luhn de la longueur fournie.
func Generate(length int) Luhn {
	rnd := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	var s strings.Builder
	for i := 0; i < length-1; i++ {
		s.WriteString(strconv.Itoa(rnd.Intn(9)))
	}

	l := Luhn(s.String())
	_, res, _ := l.Calculate() // ignorer l'erreur car cela sera toujours valide
	return Luhn(res)
}

// calculateLuhnSum calcule la somme de Luhn pour un nombre donné avec une parité donnée.
func calculateLuhnSum(number string, parity int) (int64, errors.ErrorInterface) {
	var sum int64
	for i, d := range number {
		if d < asciiZero || d > asciiTen {
			return 0, errors.ErrValueIsNotNumber
		}

		d = d - asciiZero
		// Doubler la valeur de chaque deuxième chiffre.
		if i%2 == parity {
			d *= 2
			// Si le résultat de cette opération de doublement est supérieur à 9.
			if d > 9 {
				// Le même résultat final peut être trouvé en soustrayant 9 de ce résultat.
				d -= 9
			}
		}

		// Prendre la somme de tous les chiffres.
		sum += int64(d)
	}

	return sum, nil
}
