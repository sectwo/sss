package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func lagrangeInterpolation(points []*big.Int, values []*big.Int, threshold int, prime *big.Int) *big.Int {
	var result *big.Int = new(big.Int).SetInt64(0)
	for i := 0; i < threshold; i++ {
		numerator := new(big.Int).SetInt64(1)
		denominator := new(big.Int).SetInt64(1)
		for j := 0; j < threshold; j++ {
			if i != j {
				numerator.Mul(numerator, new(big.Int).Neg(points[j])).Mod(numerator, prime)
				denominator.Mul(denominator, new(big.Int).Sub(points[i], points[j])).Mod(denominator, prime)
			}
		}
		denominator.ModInverse(denominator, prime)
		term := new(big.Int).Mul(values[i], numerator)
		term.Mul(term, denominator).Mod(term, prime)
		result.Add(result, term).Mod(result, prime)
	}
	return result
}

func evaluatePolynomial(coefficients []*big.Int, x *big.Int, prime *big.Int) *big.Int {
	result := new(big.Int).SetInt64(0)
	xi := new(big.Int).SetInt64(1)
	for _, coefficient := range coefficients {
		term := new(big.Int).Mul(xi, coefficient)
		result.Add(result, term).Mod(result, prime)
		xi.Mul(xi, x).Mod(xi, prime)
	}
	return result
}

func shareSecret(secret *big.Int, threshold int, numShares int, prime *big.Int) ([]*big.Int, []*big.Int) {
	coefficients := make([]*big.Int, threshold)
	coefficients[0] = secret
	for i := 1; i < threshold; i++ {
		coefficients[i], _ = rand.Int(rand.Reader, prime)
	}

	points := make([]*big.Int, numShares)
	values := make([]*big.Int, numShares)
	for i := 0; i < numShares; i++ {
		points[i] = big.NewInt(int64(i + 1))
		values[i] = evaluatePolynomial(coefficients, points[i], prime)
	}
	return points, values
}

func recoverSecret(points []*big.Int, values []*big.Int, threshold int, prime *big.Int) *big.Int {
	return lagrangeInterpolation(points, values, threshold, prime)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// 사용자로부터 비밀 입력 받기
	fmt.Print("Enter the secret (integer): ")
	secretStr, _ := reader.ReadString('\n')
	secretStr = strings.TrimSpace(secretStr)
	secret, _ := new(big.Int).SetString(secretStr, 10)

	// 사용자로부터 복구에 필요한 부분의 개수 입력 받기
	fmt.Print("Enter the threshold (minimum number of shares required to recover the secret): ")
	thresholdStr, _ := reader.ReadString('\n')
	thresholdStr = strings.TrimSpace(thresholdStr)
	threshold, _ := strconv.Atoi(thresholdStr)

	// 사용자로부터 생성할 부분의 총 개수 입력 받기
	fmt.Print("Enter the total number of shares to generate: ")
	numSharesStr, _ := reader.ReadString('\n')
	numSharesStr = strings.TrimSpace(numSharesStr)
	numShares, _ := strconv.Atoi(numSharesStr)

	prime, _ := rand.Prime(rand.Reader, 128)
	fmt.Printf("Prime: %s\n", prime.String())

	points, shares := shareSecret(secret, threshold, numShares, prime)
	fmt.Printf("Secret: %s\n", secret.String())
	fmt.Println("Shares:")
	for i := 0; i < numShares; i++ {
		fmt.Printf("  (%s, %s)\n", points[i].String(), shares[i].String())
	}

	// 사용자로부터 복구에 사용할 부분 선택
	selectedPoints := make([]*big.Int, threshold)
	selectedShares := make([]*big.Int, threshold)
	for i := 0; i < threshold; i++ {
		fmt.Printf("Enter the %dth share (format: x,y): ", i+1)
		shareStr, _ := reader.ReadString('\n')
		shareStr = strings.TrimSpace(shareStr)
		shareParts := strings.Split(shareStr, ",")

		x, _ := new(big.Int).SetString(shareParts[0], 10)
		y, _ := new(big.Int).SetString(shareParts[1], 10)

		selectedPoints[i] = x
		selectedShares[i] = y
	}

	recoveredSecret := recoverSecret(selectedPoints, selectedShares, threshold, prime)
	fmt.Printf("Recovered secret: %s\n", recoveredSecret.String())
}
