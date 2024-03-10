package bloomfilter

import (
	"SQL/logs"
	"github.com/demdxx/gocast"
)

type LocalBloomService struct {
	m, k, n   int32
	bitmap    []int
	encryptor *Encryptor
}

// NewLocalBloomService 新建布隆过滤器
func NewLocalBloomService(m, k int32, encryptor *Encryptor) *LocalBloomService {
	return &LocalBloomService{
		m:         m,
		k:         k,
		bitmap:    make([]int, m/32+1),
		encryptor: encryptor,
	}
}

// Exist 判断是否存在于布隆过滤器里面
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

// getKEncrypted 获取一个元素 val 对应 k 个 bit 位偏移量 offset
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

// Set 添加进布隆过滤器
func (l *LocalBloomService) Set(val string) {
	logs.SugarLogger.Info("LocalBloomService is setting")
	l.n++
	for _, offset := range l.getKEncrypted(val) {
		index := offset >> 5     // 等价于 / 32
		bitOffset := offset & 31 // 等价于 % 32

		l.bitmap[index] |= (1 << bitOffset)
	}
}
