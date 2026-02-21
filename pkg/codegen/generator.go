package codegen

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

func (g *codeGenerator) Generate() (string, error) {
	var parts []string

	for _, size := range g.groups {	
		randomString, err := g.generateRandomString(size)
		if err != nil {
			return "", err
		}
		parts = append(parts, randomString)
	}

	return strings.Join(parts, g.separator), nil
}

func (g *codeGenerator) generateRandomString(length int) (string, error) {
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(g.charset)))

	for i := range length {
		idx, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = g.charset[idx.Int64()]
	}

	return string(result), nil
}


