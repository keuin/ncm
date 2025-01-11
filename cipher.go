package ncm

import (
	"crypto/cipher"
	"fmt"
	"io"
)

// streamBlockCipherDecrypt converts a block cipher in ECB mode to io.Reader interface.
type streamBlockCipherDecrypt struct {
	R      io.Reader
	Cipher cipher.Block
	j, k   int // 0 <= j <= k <= len(buf)
	buf    [32]byte
}

func (s *streamBlockCipherDecrypt) Read(dst []byte) (n int, err error) {
	bs := s.Cipher.BlockSize()
	if s.k != s.j {
		// buffer has data to consume
		n = copy(dst, s.buf[s.j:s.k])
		dst = dst[n:]
		s.j += n
	}
	if len(dst) == 0 {
		return n, nil
	}
	var n2 int
	// buffer is empty, but dst is not completely filled, read more
	for len(dst) >= bs {
		n2, err = io.ReadFull(s.R, dst[:bs])
		if (err == io.EOF || err == io.ErrUnexpectedEOF) && n2 != bs {
			return n, io.EOF
		}
		if err != nil {
			return 0, fmt.Errorf("insufficient read: %w", err)
		}
		assert(n2 == bs, "bad read")
		s.Cipher.Decrypt(dst, dst)
		dst = dst[bs:]
		n += bs
	}
	if len(dst) == 0 {
		return n, nil
	}
	// left part is smaller than one block, read one block into buffer
	assert(len(dst) < bs, "bad dst size")
	n2, err = io.ReadFull(s.R, s.buf[:bs])
	if (err == io.EOF || err == io.ErrUnexpectedEOF) && n2 == 0 {
		return n, io.EOF
	}
	if err != nil && (err != io.EOF || n2 != bs) {
		return 0, fmt.Errorf("insufficient read: %w", err)
	}
	s.Cipher.Decrypt(s.buf[:bs], s.buf[:bs])
	n2 = copy(dst, s.buf[:bs])
	n += n2
	s.j = n2
	s.k = bs
	assert(s.k > s.j && s.k <= len(s.buf), "bad buffer")
	return n, err
}

type xorReader struct {
	R    io.Reader
	Byte byte
}

func (x xorReader) Read(buf []byte) (n int, err error) {
	n, err = x.R.Read(buf)
	for i := 0; i < n; i++ {
		buf[i] ^= x.Byte
	}
	return n, err
}
