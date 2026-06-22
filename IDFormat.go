package snowflake

import (
	"encoding/base64"
	"encoding/binary"
	"strconv"
)

// ID 是用于雪花 ID 的自定义类型，以便我们可以在 ID 上附加方法。
type ID int64

// Int64 返回雪花 ID 的 int64 值
func (f ID) Int64() int64 {
	return int64(f)
}

// Base64 返回雪花 ID 的 base64 字符串
func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

// Base58 返回雪花 ID 的 base58 字符串
func (f ID) Base58() string {

	if f < 58 {
		return string(encodeBase58Map[f])
	}

	b := make([]byte, 0, 11)
	for f >= 58 {
		b = append(b, encodeBase58Map[f%58])
		f /= 58
	}
	b = append(b, encodeBase58Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// Base36 返回雪花 ID 的 bas36 字符串
func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

// Base32 使用 z-base-32 字符集，但编码和解码方式类似于 base58，从而可以生成更短的结果字符串。
// 注意：存在许多不同的 base32 实现，进行互操作时请务必小心。
func (f ID) Base32() string {

	if f < 32 {
		return string(encodeBase32Map[f])
	}

	b := make([]byte, 0, 12)
	for f >= 32 {
		b = append(b, encodeBase32Map[f%32])
		f /= 32
	}
	b = append(b, encodeBase32Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// String 返回雪花 ID 的字符串表示
func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

// Base2 返回雪花 ID 的 base2 字符串
func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

// Bytes 返回雪花 ID 的字节切片
func (f ID) Bytes() []byte {
	return []byte(f.String())
}

// IntBytes 返回雪花 ID 的字节数组，以大端序整数编码。
func (f ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(f))
	return b
}

// Time 返回雪花 ID 时间的 int64 毫秒级 Unix 时间戳
func (f ID) Time() int64 {
	return (int64(f) & timeMask) + Epoch
}

// Node 返回雪花 ID 节点编号的 int64 值
func (f ID) Node() int64 {
	return (int64(f) & nodeMask) >> nodeShift
}

// Step 返回雪花 ID 步长（或序列号）的 int64 值
func (f ID) Step() int64 {
	return (int64(f) >> stepShift) & stepMask
}
