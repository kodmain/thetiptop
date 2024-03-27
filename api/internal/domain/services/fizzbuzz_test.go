package services_test

import (
	"reflect"
	"testing"

	"github.com/kodmain/thetiptop/api/internal/domain/services"
	"github.com/stretchr/testify/assert"
)

func TestNewFizzBuzz(t *testing.T) {
	fizzBuzzService := services.NewFizzBuzzService()
	assert.NotNil(t, fizzBuzzService, "FizzBuzzService should not be nil")
}

func TestRunFizzBuzz(t *testing.T) {
	tests := []struct {
		name      string
		int1      int
		int2      int
		limit     int
		str1      string
		str2      string
		want      []string
		wantError bool
	}{
		{
			name:  "Basic FizzBuzz",
			int1:  3,
			int2:  5,
			limit: 15,
			str1:  "Fizz",
			str2:  "Buzz",
			want:  []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := services.NewFizzBuzzService()
			got := s.RunFizzBuzz(tt.int1, tt.int2, tt.limit, tt.str1, tt.str2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunFizzBuzz() got = %v, want %v", got, tt.want)
			}
		})
	}
}
