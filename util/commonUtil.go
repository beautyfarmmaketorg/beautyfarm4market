package util

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMd5(orignStr string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(orignStr))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
