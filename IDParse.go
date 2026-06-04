package snowflake

import (
	"encoding/base64"
	"encoding/binary"
	"strconv"
)

// ParseInt64 将 int64 转换为雪花 ID
func ParseInt64(id int64) ID {
	return ID(id)
}

// ParseString 将字符串转换为雪花 ID
func ParseString(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	return ID(i), err

}

// ParseBase2 将 Base2 字符串转换为雪花 ID
func ParseBase2(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 2, 64)
	return ID(i), err
}

// ParseBase32 将 base32 []byte 解析为雪花 ID
// 注意：存在许多不同的 base32 实现，进行互操作时请务必小心。
func ParseBase32(s string) (ID, error) {
	var id int64
	b := []byte(s)

	for i := range b {
		if decodeBase32Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase32
		}
		id = id*32 + int64(decodeBase32Map[b[i]])
	}

	return ID(id), nil
}

// ParseBase36 将 Base36 字符串转换为雪花 ID
func ParseBase36(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 36, 64)
	return ID(i), err
}

// ParseBase58 将 base58 []byte 解析为雪花 ID
func ParseBase58(s string) (ID, error) {

	var id int64
	b := []byte(s)

	for i := range b {
		if decodeBase58Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase58
		}
		id = id*58 + int64(decodeBase58Map[b[i]])
	}

	return ID(id), nil
}

// ParseBase64 将 base64 字符串转换为雪花 ID
func ParseBase64(id string) (ID, error) {
	b, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return -1, err
	}
	return ParseBytes(b)

}

// ParseBytes 将字节切片转换为雪花 ID
func ParseBytes(id []byte) (ID, error) {
	i, err := strconv.ParseInt(string(id), 10, 64)
	return ID(i), err
}

// ParseIntBytes 将以大端序整数编码的字节数组转换为雪花 ID
func ParseIntBytes(id [8]byte) ID {
	return ID(int64(binary.BigEndian.Uint64(id[:])))
}
