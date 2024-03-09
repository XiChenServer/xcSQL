package bloomfilter

import (
	"github.com/spaolacci/murmur3"
	"math"
)

type Encryptor struct {
}

func NewEncryptor() *Encryptor {
	return &Encryptor{}
}

func (e *Encryptor) Encrypt(origin string) int32 {
	hasher := murmur3.New32()
	_, _ = hasher.Write([]byte(origin))
	return int32(hasher.Sum32() % math.MaxInt32)
}
