package public

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
)

func EncryptPassword(salt, password string) string {
	s1 := sha256.New()
	s1.Write([]byte(password))
	str := fmt.Sprintf("%x", s1.Sum(nil))
	s2 := sha256.New()
	s2.Write([]byte(str + salt))
	return fmt.Sprintf("%x", s2.Sum(nil))
}

func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
