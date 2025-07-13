package common

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

var letters = []rune("asbsdfasdfLKJLJKJDSAJLKJKPOUPUKJLKJKLJDKLAJSK")

func RandomSequence(n int) string {
	b := make([]rune, n)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for i := range b {
		b[i] = letters[r1.Intn(999999)%len(letters)]
	}
	return string(b)
}

func GenSalt(length int) string {
	if length < 0 {
		length = 50
	}
	return RandomSequence(length)
}

// Ma hoa
type bcryptHasher struct {
	cost int
}

// contructor
func NewBcryptHasher(cost int) *bcryptHasher {
	if cost <= 0 {
		cost = bcrypt.DefaultCost // 10
	}
	return &bcryptHasher{cost: cost}
}

func (h *bcryptHasher) Hash(data string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(data), h.cost)
	if err != nil {
		return "" // nil de tang biz xu ly
	}
	return string(hashed)
}
func (h *bcryptHasher) Compare(hashedData, plainData string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedData), []byte(plainData))
	return err == nil
}
