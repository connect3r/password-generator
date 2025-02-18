package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"strings"
	"unicode"
)

var (
	length          int
	digits          bool
	symbols         bool
	upper           bool
	lower           bool
	excludeSimilar  bool
	excludeAmbiguous bool
)

const (
	digitChars    = "23456789"
	symbolChars   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	upperChars    = "ABCDEFGHJKLMNPQRSTUVWXYZ"
	lowerChars    = "abcdefghijkmnpqrstuvwxyz"
	similarChars  = "il1Lo0O"
	ambiguousChars = "{}[]()/\\'\"`~,;:.<>"
)

func init() {
	flag.IntVar(&length, "length", 16, "password length")
	flag.BoolVar(&digits, "digits", true, "include digits")
	flag.BoolVar(&symbols, "symbols", true, "include symbols")
	flag.BoolVar(&upper, "upper", true, "include uppercase letters")
	flag.BoolVar(&lower, "lower", true, "include lowercase letters")
	flag.BoolVar(&excludeSimilar, "exclude-similar", true, "exclude similar characters")
	flag.BoolVar(&excludeAmbiguous, "exclude-ambiguous", true, "exclude ambiguous characters")
}

func main() {
	flag.Parse()
	password, err := GeneratePassword(length, digits, symbols, upper, lower, excludeSimilar, excludeAmbiguous)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Generated password:", password)
}

func GeneratePassword(length int, digits, symbols, upper, lower, excludeSimilar, excludeAmbiguous bool) (string, error) {
	if length < 8 {
		return "", errors.New("password length must be at least 8 characters")
	}

	var charPool string
	var requiredChars []string

	if digits {
		d := filterChars(digitChars, excludeSimilar, excludeAmbiguous)
		charPool += d
		requiredChars = append(requiredChars, randomChar(d))
	}

	if symbols {
		s := filterChars(symbolChars, excludeSimilar, excludeAmbiguous)
		charPool += s
		requiredChars = append(requiredChars, randomChar(s))
	}

	if upper {
		u := filterChars(upperChars, excludeSimilar, excludeAmbiguous)
		charPool += u
		requiredChars = append(requiredChars, randomChar(u))
	}

	if lower {
		l := filterChars(lowerChars, excludeSimilar, excludeAmbiguous)
		charPool += l
		requiredChars = append(requiredChars, randomChar(l))
	}

	if charPool == "" {
		return "", errors.New("at least one character set must be selected")
	}

	if length < len(requiredChars) {
		return "", fmt.Errorf("password length must be at least %d for selected character sets", len(requiredChars))
	}

	password := make([]rune, length)
	for i, c := range requiredChars {
		password[i] = rune(c[0])
	}

	for i := len(requiredChars); i < length; i++ {
		password[i] = rune(randomChar(charPool)[0])
	}

	shuffle(password)

	return string(password), nil
}

func filterChars(chars string, excludeSimilar, excludeAmbiguous bool) string {
	var result strings.Builder
	for _, c := range chars {
		if excludeSimilar && strings.ContainsRune(similarChars, c) {
			continue
		}
		if excludeAmbiguous && strings.ContainsRune(ambiguousChars, c) {
			continue
		}
		result.WriteRune(c)
	}
	return result.String()
}

func randomChar(pool string) string {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))
	if err != nil {
		panic(err)
	}
	return string(pool[n.Int64()])
}

func shuffle(password []rune) {
	for i := len(password) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			panic(err)
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}
}