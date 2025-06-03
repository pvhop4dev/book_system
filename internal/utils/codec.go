package utils

import (
	"book_system/internal/config"
	"fmt"
	"strings"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var secretKey uint32 = config.MustGet().Codec.SecretKey

func encodeID(id uint32) string {
	obfuscated := id ^ secretKey
	if obfuscated == 0 {
		return string(base62Chars[0])
	}
	var sb strings.Builder
	for obfuscated > 0 {
		remainder := obfuscated % 62
		sb.WriteByte(base62Chars[remainder])
		obfuscated /= 62
	}

	runes := []rune(sb.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func encodeBase62(num uint64) string {
	if num == 0 {
		return string(base62Chars[0])
	}
	var sb strings.Builder
	for num > 0 {
		remainder := num % 62
		sb.WriteByte(base62Chars[remainder])
		num /= 62
	}

	runes := []rune(sb.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func decodeBase62(s string) (uint64, error) {
	var num uint64
	for _, c := range s {
		index := strings.IndexRune(base62Chars, c)
		if index == -1 {
			return 0, fmt.Errorf("invalid character: %c", c)
		}
		num = num*62 + uint64(index)
	}
	return num, nil
}

func combineIDs(a, b uint32) uint64 {

	var x, y uint32
	if a < b {
		x, y = uint32(a), uint32(b)
	} else {
		x, y = uint32(b), uint32(a)
	}
	combined := (uint64(x) << 32) | uint64(y)
	return combined ^ uint64(secretKey)
}

func splitIDs(encoded uint64) (uint32, uint32) {
	decoded := encoded ^ uint64(secretKey)
	a := uint32(decoded >> 32)
	b := uint32(decoded & 0xFFFFFFFF)
	return a, b
}

func generateRoomId(a, b uint32) string {
	encoded := combineIDs(a, b)
	return encodeBase62(encoded)
}

func decodeRoomId(roomId string) (uint32, uint32, error) {
	num, err := decodeBase62(roomId)
	if err != nil {
		return 0, 0, err
	}
	a, b := splitIDs(num)
	return a, b, nil
}

func decodeID(s string) (uint32, error) {
	var num uint32
	for _, c := range s {
		index := strings.IndexRune(base62Chars, c)
		if index == -1 {
			return 0, fmt.Errorf("invalid character in ID")
		}
		num = num*62 + uint32(index)
	}
	return num ^ secretKey, nil
}
