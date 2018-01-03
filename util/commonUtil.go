package util

import (
	"crypto/md5"
	"encoding/hex"
	"crypto/sha1"
	"fmt"
)

func GetMd5(orignStr string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(orignStr))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func GetSha1(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs);
}
