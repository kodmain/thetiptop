package services

import "fmt"

type FizzBuzzService struct {
}

func NewFizzBuzzService() *FizzBuzzService {
	return &FizzBuzzService{}
}

func (s *FizzBuzzService) RunFizzBuzz(int1, int2, limit int, str1, str2 string) []string {
	var result []string
	for i := 1; i <= limit; i++ {
		switch {
		case i%int1 == 0 && i%int2 == 0:
			result = append(result, str1+str2)
		case i%int1 == 0:
			result = append(result, str1)
		case i%int2 == 0:
			result = append(result, str2)
		default:
			result = append(result, fmt.Sprintf("%d", i))
		}
	}
	return result
}
