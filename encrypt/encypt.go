package encrypt

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func Md5s(s string) string {
	r := md5.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

func Md5f(fName string) string {
	f, e := os.Open(fName)
	if e != nil {
		log.Fatal(e)
	}
	h := md5.New()
	_, e = io.Copy(h, f)
	if e != nil {
		log.Fatal(e)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1s(s string) string {
	r := sha1.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

func Sha1f(fName string) string {
	f, e := os.Open(fName)
	if e != nil {
		log.Fatal(e)
	}
	h := sha1.New()
	_, e = io.Copy(h, f)
	if e != nil {
		log.Fatal(e)
	}
	return hex.EncodeToString(h.Sum(nil))
}
