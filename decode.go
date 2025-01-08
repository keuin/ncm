package ncm

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

const (
	masterKey      = "hzHRAmso5kInbaxW"
	metadataKey    = "#14ljk_!\\]&0U<'("
	fileHeader     = "CTENFDAM"
	metadataHeader = "163 key(Don't modify):"
)

type offsetReader struct {
	R      io.Reader
	Offset int64
}

func (o *offsetReader) Read(buf []byte) (n int, err error) {
	n, err = o.R.Read(buf)
	o.Offset += int64(n)
	return n, err
}

// Decoder decodes NCM files.
// Do not instantiate it directly, use NewDecoder instead.
type Decoder struct {
	or         offsetReader
	keyBox     []byte
	dataOffset int64
	Metadata   struct {
		// Format is music file extension name, e.g. "mp3"
		Format string `json:"format"`
	}
}

func (d *Decoder) readHeader() error {
	if d.or.Offset != 0 {
		return nil
	}
	var buf [8]byte
	_, err := io.ReadFull(&d.or, buf[:])
	if err != nil {
		return fmt.Errorf("read file header: %w", err)
	}
	if unsafeString(buf[:]) != fileHeader {
		return fmt.Errorf("invalid file header: 0x%s", hex.EncodeToString(buf[:]))
	}
	err = skip(&d.or, 2)
	if err != nil {
		return fmt.Errorf("read padding #1: %w", err)
	}
	_, err = io.ReadFull(&d.or, buf[:4])
	if err != nil {
		return fmt.Errorf("read keylen: %w", err)
	}
	keyLen := int(binary.LittleEndian.Uint32(buf[:]))
	key := make([]byte, keyLen)
	_, err = io.ReadFull(&d.or, key)
	if err != nil {
		return fmt.Errorf("read key: %w", err)
	}
	for i := range key {
		key[i] ^= 0x64
	}
	masterCipher := newCipher(masterKey)
	decryptAll(masterCipher, key)
	key = unpad(key)[17:]
	keyLen = len(key)
	keyBox := make([]byte, 256)
	for i := range keyBox {
		keyBox[i] = byte(i)
	}
	d.keyBox = keyBox
	var prev byte
	var offset int
	for i := range keyBox {
		v := keyBox[i]
		c := (v + prev + key[offset]) & 0xff
		offset = (offset + 1) % keyLen
		keyBox[i] = keyBox[c]
		keyBox[c] = v
		prev = c
	}
	_, err = io.ReadFull(&d.or, buf[:4])
	if err != nil {
		return fmt.Errorf("read metadataLen: %w", err)
	}
	metadataLen := int(binary.LittleEndian.Uint32(buf[:]))
	metadata := make([]byte, metadataLen)
	_, err = io.ReadFull(&d.or, metadata)
	if err != nil {
		return err
	}
	for i := range metadata {
		metadata[i] ^= 0x63
	}
	if !bytes.HasPrefix(metadata, unsafeBytes(metadataHeader)) {
		return fmt.Errorf("invalid metadata header: 0x%s", hex.EncodeToString(metadata))
	}
	metadata, err = base64.StdEncoding.DecodeString(unsafeString(metadata[len(metadataHeader):]))
	if err != nil {
		return fmt.Errorf("decode metadata base64 string: %w", err)
	}
	metadataCipher := newCipher(metadataKey)
	decryptAll(metadataCipher, metadata)
	metadata = unpad(metadata)[6:]
	err = json.Unmarshal(metadata, &d.Metadata)
	if err != nil {
		return fmt.Errorf("unmarshal metadata JSON: 0x%s, error: %w",
			hex.Dump(metadata), err)
	}
	err = skip(&d.or, 4+5) // crc32 sum (4byte) + dummy data (5byte)
	if err != nil {
		return fmt.Errorf("read padding #2: %w", err)
	}
	_, err = io.ReadFull(&d.or, buf[:4])
	if err != nil {
		return fmt.Errorf("read imageSizeBytes: %w", err)
	}
	imageSizeBytes := int64(binary.LittleEndian.Uint32(buf[:]))
	err = skip(&d.or, imageSizeBytes)
	if err != nil {
		return fmt.Errorf("skip image data: %w", err)
	}
	d.dataOffset = d.or.Offset
	return nil
}

func (d *Decoder) Read(buf []byte) (n int, err error) {
	offset := int(d.or.Offset - d.dataOffset)
	n, err = d.or.Read(buf)
	for i := 1; i <= n; i++ {
		j := (i + offset) & 0xff
		buf[i-1] ^= d.keyBox[(d.keyBox[j]+d.keyBox[(int(d.keyBox[j])+j)&0xff])&0xff]
	}
	return n, err
}

// NewDecoder accepts an io.Reader to NCM data and creates a Decoder instance
// which yields decrypted music data. Metadata such as original file extension,
// will be available in Decoder.Metadata.
func NewDecoder(r io.Reader) (*Decoder, error) {
	ret := &Decoder{
		or: offsetReader{
			R: r,
		},
	}
	err := ret.readHeader()
	if err != nil {
		return nil, err
	}
	return ret, err
}
