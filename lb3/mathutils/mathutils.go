package mathutils

import "errors"

func Factorial(n int) (uint64, error) {
	if n < 0 {
		return 0, errors.New("факториал отрицательного числа не определён")
	}
	if n > 20 {
		return 0, errors.New("число слишком большое, результат превысит uint64")
	}

	result := uint64(1)
	for i := 2; i <= n; i++ {
		result *= uint64(i)
	}
	return result, nil
}
