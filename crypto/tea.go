package crypto

import (
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/math"
	tea "github.com/fumiama/gofastTEA"
)

// GetTEA 从string生成TEA密钥
func GetTEA(key string) tea.TEA {
	if len(key) == 0 {
		return tea.TEA{}
	}
	if len(key) > 16 {
		key = key[:16]
	} else {
		for len(key) < 16 {
			key += key[:math.Min(16-len(key), len(key))]
		}
	}
	return tea.NewTeaCipher(binary.StringToBytes(key))
}
