package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvxyz"

// init() function will be called automatically when the package is first used.
// we will set the seed value for the random generator by calling rand.Seed()
func intit() {
	// the seed value is often set to the current time.
	// because rand.Seed() take int64 as input
	// we should conver the time to unix nano before passing it to the function.
	// that will make sure every time we run the code
	// the generated values will be different.
	// if we didn't call rand.Seed(), the random generator will behave like
	// it is seeded by 1. So the generated values will be the same for every run.
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min amd max
func RandomInt(min, max int64) int64 {
	// rand.Int63n function returns a random integer between 0 and n-1
	return min + rand.Int63n(max-min+1) // 0 ~ (max - min)
}

// RandomString to generate a random string of length n.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	// use for loop to generate n random characters
	for i := 0; i < n; i++ {
		// use rand.Intn(k) to get a random position from 0 ~ k-1
		// and take the corresponding character at that position in the alphabet
		c := alphabet[rand.Intn(k)]
		// write that character c to the string builder
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "TWD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
