package paramValidator

import (
	"github.com/go-playground/validator/v10"
	"math"
	"strings"
)

const THRESHOLD float64 = 2.50

func calcStringEntropy(data string) (entropy float64) {
	if data == "" {
		return 0
	}
	for i := 0; i < 256; i++ {
		px := float64(strings.Count(data, string(byte(i)))) / float64(len(data))
		if px > 0 {
			entropy += -px * math.Log2(px)
		}
	}
	return entropy
}

var PasswordStrengthValidator validator.Func = func(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	entropy := calcStringEntropy(password)
	return entropy > THRESHOLD
}
