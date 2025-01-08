package ncm

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func unsafeBytes(s string) []byte {
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return unsafe.Slice((*byte)(unsafe.Pointer(h.Data)), h.Len)
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

func unpad(b []byte) ([]byte, error) {
	if len(b) == 0 {
		return nil, errors.New("empty ciphertext")
	}
	n := len(b) - int(b[len(b)-1])
	if n < 0 {
		return nil, fmt.Errorf("invalid padding: %v", int(b[len(b)-1]))
	}
	return b[:n], nil
}

func decryptAll(c cipher.Block, buf []byte) error {
	bs := c.BlockSize()
	if len(buf)%bs != 0 {
		return fmt.Errorf("invalid ciphertext length: %v", len(buf))
	}
	for len(buf) > 0 {
		c.Decrypt(buf, buf)
		buf = buf[bs:]
	}
	return nil
}
