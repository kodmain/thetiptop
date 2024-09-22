package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TYPE uint8 // Type de jeton

const (
	ACCESS  TYPE = 0 // Jeton d'accès
	REFRESH TYPE = 1 // Jeton de rafraîchissement
)

type Token struct {
	ID     string         `json:"id"`
	Exp    int64          `json:"exp"`
	TZ     string         `json:"tz"`
	Offset int            `json:"offset"`
	Type   TYPE           `json:"type"`
	Data   map[string]any `json:"data"`
}

func (t *Token) IsNotValid() bool {
	return t.Type != ACCESS
}

func (t *Token) HasExpired() bool {
	// Charger le fuseau horaire spécifié ou utiliser UTC si une erreur survient
	location, err := time.LoadLocation(t.TZ)
	if err != nil {
		log.Printf("Erreur lors du chargement du fuseau horaire '%s': %v. Utilisation de UTC.", t.TZ, err)
		location = time.UTC // Utiliser UTC si le fuseau horaire n'est pas valide
	}

	// Convertir le timestamp d'expiration en Time, ajusté par le fuseau horaire et l'offset
	expirationTime := time.Unix(t.Exp, 0).In(location)

	// Obtenir le temps courant dans le même fuseau horaire pour une comparaison cohérente
	currentTime := time.Now().In(location)

	// Vérifier si le temps courant est après le temps d'expiration
	return currentTime.After(expirationTime)
}
func (a Token) Claims() jwt.MapClaims {
	claims := jwt.MapClaims{
		"id":   a.ID,
		"exp":  a.Exp,
		"tz":   a.TZ,
		"off":  a.Offset,
		"type": a.Type,
		"data": a.Data,
	}

	return claims
}

func fromClaims(claims jwt.MapClaims) *Token {
	token := &Token{}

	if id, ok := claims["id"].(string); ok {
		token.ID = id
	}

	if exp, ok := claims["exp"]; ok {
		token.Exp = convertToInt64(exp)
	}

	if tz, ok := claims["tz"].(string); ok {
		token.TZ = tz
	}

	if typ, ok := claims["type"]; ok {
		token.Type = TYPE(convertToInt(typ))
	}

	if offset, ok := claims["off"]; ok {
		token.Offset = convertToInt(offset)
	}

	if data, ok := claims["data"].(map[string]any); ok {
		token.Data = data
	}

	return token
}

func convertToInt64(val interface{}) int64 {
	switch v := val.(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	default:
		log.Fatalf("Invalid type for int64 conversion: %T\n", v)
		return 0
	}
}

func convertToInt(val interface{}) int {
	switch v := val.(type) {
	case float64:
		return int(v)
	case int:
		return v
	default:
		log.Fatalf("Invalid type for int conversion: %T\n", v)
		return 0
	}
}
