package bloomfilter

import "github.com/demdxx/gocast"

type LocalBloomService struct {
	m, k, n   int32
	bitmap    []int
	encryptor *Encryptor
}

func NewLocalBloomService(m, k int32, encryptor *Encryptor) *LocalBloomService {
	return &LocalBloomService{
		m:         m,
		k:         k,
		bitmap:    make([]int, m/32+1),
		encryptor: encryptor,
	}
}
func (l *LocalBloomService) Exist(val string) bool {
	for _, offset := range l.getKEncrypted(val) {
		index := offset >> 5     // 等价于 / 32
		bitOffset := offset & 31 // 等价于 % 32

		if l.bitmap[index]&(1<<bitOffset) == 0 {
			return false
		}
	}

	return true
}
func (l *LocalBloomService) getKEncrypted(val string) []int32 {
	encrypteds := make([]int32, 0, l.k)
	origin := val
	for i := 0; int32(i) < l.k; i++ {
		encrypted := l.encryptor.Encrypt(origin)
		encrypteds = append(encrypteds, encrypted%l.m)
		if int32(i) == l.k-1 {
			break
		}
		origin = gocast.ToString(encrypted)
	}
	return encrypteds
}

func (l *LocalBloomService) Set(val string) {
	l.n++
	for _, offset := range l.getKEncrypted(val) {
		index := offset >> 5     // 等价于 / 32
		bitOffset := offset & 31 // 等价于 % 32

		l.bitmap[index] |= (1 << bitOffset)
	}
}
