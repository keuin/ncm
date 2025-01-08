package ncm

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"unsafe"
)

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func unsafeBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func skip(r io.Reader, n int64) error {
	if n == 0 {
		return nil
	}
	_, err := io.CopyN(io.Discard, r, n)
	return err
}

func newCipher(key string) cipher.Block {
	b, err := aes.NewCipher(unsafeBytes(key))
	if err != nil {
		panic(err)
	}
	return b
}

func unpad(b []byte) []byte {
	return b[:len(b)-int(b[len(b)-1])]
}

func decryptAll(c cipher.Block, buf []byte) {
	bs := c.BlockSize()
	for len(buf) > 0 {
		c.Decrypt(buf, buf)
		buf = buf[bs:]
	}
}
