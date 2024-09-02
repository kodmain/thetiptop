package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/password"
)

type JWT struct {
	TZ       string        `yaml:"tz"`
	Secret   string        `yaml:"secret"`
	Expire   int           `yaml:"expire"`
	Refresh  int           `yaml:"refresh"`
	Duration time.Duration `yaml:"duration"`
}

var (
	instance *JWT
	duration time.Duration = time.Minute
)

func New(t *JWT) error {
	location := time.Now().Location().String()
	instance = t

	if instance == nil {
		pass, err := password.GeneratePassword(32, password.All)
		if err != nil {
			return err
		}

		instance = &JWT{
			TZ:       location,
			Secret:   pass,
			Expire:   15,
			Refresh:  30,
			Duration: duration,
		}
	}

	if instance.Duration <= 0 {
		instance.Duration = duration
	}

	if instance.Expire < 1 {
		instance.Expire = 15
	}

	if instance.Refresh < 1 {
		instance.Refresh = 30
	}

	if instance.TZ == "" {
		instance.TZ = location
	}

	if instance.Secret == "" {
		pass, err := password.GeneratePassword(32, password.All)
		if err != nil {
			return err
		}

		instance.Secret = pass
	}

	return nil
}

func FromID(id string) (string, string, error) {
	location, err := time.LoadLocation(instance.TZ)
	if err != nil {
		return "", "", err
	}

	now := time.Now().In(location)
	_, offset := now.Zone()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Token{
		ID:     id,
		Exp:    now.Add(instance.Duration * time.Duration(instance.Refresh)).Unix(),
		TZ:     location.String(),
		Type:   REFRESH,
		Offset: offset,
	}.Claims())

	refresh, err := token.SignedString([]byte(instance.Secret))
	if err != nil {
		return "", "", err
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, Token{
		ID:     id,
		Exp:    now.Add(instance.Duration * time.Duration(instance.Expire)).Unix(),
		TZ:     location.String(),
		Offset: offset,
		Type:   ACCESS,
		//Refresh: &refresh,
	}.Claims())

	access, err := token.SignedString([]byte(instance.Secret))
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func TokenToClaims(tokenString string) (*Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(instance.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return fromClaims(claims), nil
}
