package token_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
)

func TestLuhn_Validate(t *testing.T) {
	t.Run("when the number is valid", func(t *testing.T) {
		l := token.NewLuhn("79927398713")

		err := l.Validate()

		assert.NoError(t, err)
	})

	t.Run("when the number is invalid", func(t *testing.T) {
		l := token.NewLuhn("79927398710")

		err := l.Validate()

		assert.Error(t, err)
	})
}

func TestLuhn_Calculate(t *testing.T) {
	t.Run("when the number is odd", func(t *testing.T) {
		l := token.NewLuhn("7992739871")

		expected := "3"
		expectedNumber := "79927398713"

		result, number, err := l.Calculate()

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, expectedNumber, number)
	})

	t.Run("when the number is even", func(t *testing.T) {
		l := token.NewLuhn("7992739871")

		expected := "3"
		expectedNumber := "79927398713"

		result, number, err := l.Calculate()

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, expectedNumber, number)
	})
}

func TestLuhn_String(t *testing.T) {
	l := token.NewLuhn("79927398713")

	result := l.String()

	assert.Equal(t, "79927398713", result)
}

func TestNewLuhnPointer(t *testing.T) {
	l := token.NewLuhnP(nil)
	result := l.String()
	assert.Equal(t, "", result)

	l = token.NewLuhnP(aws.String("79927398713"))
	result = l.String()
	assert.Equal(t, "79927398713", result)
}

func TestLuhn_PointerString(t *testing.T) {
	l := token.NewLuhn("79927398713")
	Pointer := l.PointerString()
	assert.Equal(t, "79927398713", *Pointer)
}

func TestGenerate(t *testing.T) {

	t.Run("when the number is valid", func(t *testing.T) {
		number := token.Generate(10)
		assert.Len(t, number, 10)
	})
}
