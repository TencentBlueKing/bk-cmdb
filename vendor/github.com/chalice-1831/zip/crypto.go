// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zip

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"errors"
	"hash"
	"io"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

type EncryptionMethod int

const (
	StandardEncryption EncryptionMethod = 1
	AES128Encryption   EncryptionMethod = 2
	AES192Encryption   EncryptionMethod = 3
	AES256Encryption   EncryptionMethod = 4

	// AES key lengths
	aes128 = 16
	aes192 = 24
	aes256 = 32
)

func aesKeyLen(strength byte) int {
	switch strength {
	case 1:
		return aes128
	case 2:
		return aes192
	case 3:
		return aes256
	default:
		return 0
	}
}

// Encryption/Decryption Errors
var (
	ErrDecryption     = errors.New("zip: decryption error")
	ErrPassword       = errors.New("zip: invalid password")
	ErrAuthentication = errors.New("zip: authentication failed")
)

// Counter (CTR) mode.
// CTR converts a block cipher into a stream cipher by
// repeatedly encrypting an incrementing counter and
// xoring the resulting stream of data with the input.

// This is a re-implementation of Go's CTR mode to allow
// for a little-endian, left-aligned uint32 counter, which
// is required for WinZip AES encryption. Go's cipher.NewCTR
// follows the NIST Standard SP 800-38A, pp 13-15
// which has a big-endian, right-aligned counter.

type ctr struct {
	b       cipher.Block
	ctr     []byte
	out     []byte
	outUsed int
}

const streamBufferSize = 512

// NewWinZipCTR returns a Stream which encrypts/decrypts using the given Block in
// counter mode. The counter is initially set to 1 per WinZip AES.
func newWinZipCTR(block cipher.Block) cipher.Stream {
	bufSize := streamBufferSize
	if bufSize < block.BlockSize() {
		bufSize = block.BlockSize()
	}
	// Set the IV (counter) to 1
	iv := make([]byte, block.BlockSize())
	iv[0] = 1
	return &ctr{
		b:       block,
		ctr:     iv,
		out:     make([]byte, 0, bufSize),
		outUsed: 0,
	}
}

func (x *ctr) refill() {
	remain := len(x.out) - x.outUsed
	if remain > x.outUsed {
		return
	}
	copy(x.out, x.out[x.outUsed:])
	x.out = x.out[:cap(x.out)]
	bs := x.b.BlockSize()
	for remain < len(x.out)-bs {
		x.b.Encrypt(x.out[remain:], x.ctr)
		remain += bs

		// Increment counter
		// for i := len(x.ctr) - 1; i >= 0; i-- {
		// 	x.ctr[i]++
		// 	if x.ctr[i] != 0 {
		// 		break
		// 	}
		// }

		// Change to allow for little-endian,
		// left-aligned counter
		for i := 0; i < len(x.ctr); i++ {
			x.ctr[i]++
			if x.ctr[i] != 0 {
				break
			}
		}
	}
	x.out = x.out[:remain]
	x.outUsed = 0
}

func (x *ctr) XORKeyStream(dst, src []byte) {
	for len(src) > 0 {
		if x.outUsed >= len(x.out)-x.b.BlockSize() {
			x.refill()
		}
		n := xorBytes(dst, src, x.out[x.outUsed:])
		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}
}

func xorBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}

// newAuthReader returns either a buffered or streaming authentication reader.
// Buffered authentication is recommended. Streaming authentication is only
// recommended if: 1. you buffer the data yourself and wait for authentication
// before streaming to another source such as the network, or 2. you just don't
// care about authenticating unknown ciphertext before use :).
func newAuthReader(akey []byte, data, adata io.Reader, streaming bool) io.Reader {
	ar := authReader{
		data:  data,
		adata: adata,
		mac:   hmac.New(sha1.New, akey),
		err:   nil,
		auth:  false,
	}
	if streaming {
		return &ar
	}
	return &bufferedAuthReader{
		ar,
		new(bytes.Buffer),
	}
}

// Streaming authentication
type authReader struct {
	data  io.Reader // data to be authenticated
	adata io.Reader // the authentication code to read
	mac   hash.Hash // hmac hash
	err   error
	auth  bool
}

func (a *authReader) Read(p []byte) (int, error) {
	if a.err != nil {
		return 0, a.err
	}
	end := false
	// read underlying data
	n, err := a.data.Read(p)
	if err != nil && err != io.EOF {
		a.err = err
		return n, a.err
	} else if err == io.EOF {
		// if we are at the end, calculate the mac
		end = true
		a.err = err
	}
	// write any data to mac
	_, err = a.mac.Write(p[:n])
	if err != nil {
		a.err = err
		return n, a.err
	}
	if end {
		ab := new(bytes.Buffer)
		_, err = io.Copy(ab, a.adata)
		if err != nil || ab.Len() != 10 {
			a.err = ErrDecryption
			return n, a.err
		}
		if !a.checkAuthentication(ab.Bytes()) {
			a.err = ErrAuthentication
			return n, a.err
		}
	}
	return n, a.err
}

// buffered authentication
type bufferedAuthReader struct {
	authReader
	buf *bytes.Buffer // buffer to store data to authenticate
}

func (a *bufferedAuthReader) Read(b []byte) (int, error) {
	// check for sticky error
	if a.err != nil {
		return 0, a.err
	}
	// make sure we have auth'ed before we send any data
	if !a.auth {
		_, err := io.Copy(a.buf, a.data)
		if err != nil {
			a.err = err
			return 0, a.err
		}
		ab := new(bytes.Buffer)
		nn, err := io.Copy(ab, a.adata)
		if err != nil {
			a.err = err
			return 0, a.err
		} else if nn != 10 {
			a.err = ErrDecryption
			return 0, a.err
		}
		_, err = a.mac.Write(a.buf.Bytes())
		if err != nil {
			a.err = err
			return 0, a.err
		}
		if !a.checkAuthentication(ab.Bytes()) {
			a.err = ErrAuthentication
			return 0, a.err
		}
	}
	// so we've authenticated the data, now just pass it on.
	n, err := a.buf.Read(b)
	if err != nil {
		a.err = err
	}
	return n, a.err
}

func (a *authReader) checkAuthentication(authcode []byte) bool {
	expectedAuthCode := a.mac.Sum(nil)
	// Truncate at the first 10 bytes
	expectedAuthCode = expectedAuthCode[:10]
	a.auth = subtle.ConstantTimeCompare(expectedAuthCode, authcode) > 0
	return a.auth
}

func checkPasswordVerification(pwvv, pwv []byte) bool {
	b := subtle.ConstantTimeCompare(pwvv, pwv) > 0
	return b
}

func generateKeys(password, salt []byte, keySize int) (encKey, authKey, pwv []byte) {
	totalSize := (keySize * 2) + 2 // enc + auth + pv sizes
	key := pbkdf2.Key(password, salt, 1000, totalSize, sha1.New)
	encKey = key[:keySize]
	authKey = key[keySize : keySize*2]
	pwv = key[keySize*2:]
	return
}

// newDecryptionReader returns an authenticated, decryption reader
func newDecryptionReader(r *io.SectionReader, f *File) (io.Reader, error) {
	keyLen := aesKeyLen(f.aesStrength)
	saltLen := keyLen / 2 // salt is half of key len
	if saltLen == 0 {
		return nil, ErrDecryption
	}
	// grab the salt and pwvv
	saltpwvv := make([]byte, saltLen+2)
	if _, err := r.Read(saltpwvv); err != nil {
		return nil, err
	}
	salt := saltpwvv[:saltLen]
	pwvv := saltpwvv[saltLen : saltLen+2]
	// generate keys only if we have a password
	if f.password == nil {
		return nil, ErrPassword
	}
	decKey, authKey, pwv := generateKeys(f.password(), salt, keyLen)
	if !checkPasswordVerification(pwv, pwvv) {
		return nil, ErrPassword
	}
	dataOff := int64(saltLen + 2)
	dataLen := int64(f.CompressedSize64 - uint64(saltLen) - 2 - 10)
	// // TODO(alex): Should the compressed sizes be fixed?
	// // Not the ideal place to do this.
	// f.CompressedSize64 = uint64(dataLen)
	// f.CompressedSize = uint32(dataLen)
	data := io.NewSectionReader(r, dataOff, dataLen)
	authOff := dataOff + dataLen
	authcode := io.NewSectionReader(r, authOff, 10)
	ar := newAuthReader(authKey, data, authcode, f.DeferAuth)
	dr := decryptStream(decKey, ar)
	if dr == nil {
		return nil, ErrDecryption
	}
	return dr, nil
}

func decryptStream(key []byte, ciphertext io.Reader) io.Reader {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	stream := newWinZipCTR(block)
	reader := &cipher.StreamReader{S: stream, R: ciphertext}
	return reader
}

// writes encrypted data to hmac as it passes through
type authWriter struct {
	hmac hash.Hash // from fw.hmac
	w    io.Writer // this will be the compCount writer
}

func (aw *authWriter) Write(p []byte) (int, error) {
	_, err := aw.hmac.Write(p)
	if err != nil {
		return 0, err
	}
	return aw.w.Write(p)
}

// writes out the salt, pwv, and then the encrypted file data
type encryptionWriter struct {
	pwv   []byte
	salt  []byte
	w     io.Writer // where to write the salt + pwv
	es    io.Writer // where to write plaintext
	first bool      // first write
	err   error
}

func (ew *encryptionWriter) Write(p []byte) (int, error) {
	if ew.err != nil {
		return 0, ew.err
	}
	if ew.first {
		// if our first time writing
		// must write out the salt and pwv first unencrypted
		_, err1 := ew.w.Write(ew.salt)
		_, err2 := ew.w.Write(ew.pwv)
		if err1 != nil || err2 != nil {
			ew.err = errors.New("zip: error writing salt or pwv")
			return 0, ew.err
		}
		ew.first = false
	}
	// now just pass on to the encryption stream
	return ew.es.Write(p)
}

func encryptStream(key []byte, w io.Writer) (io.Writer, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := newWinZipCTR(block)
	writer := &cipher.StreamWriter{S: stream, W: w}
	return writer, nil
}

// newEncryptionWriter returns an io.Writer that when written to, 1. writes
// out the salt, 2. writes out pwv, and 3. writes out authenticated, encrypted
// data. The authcode will be written out in fileWriter.close().
func newEncryptionWriter(w io.Writer, password passwordFn, fw *fileWriter, aesstrength byte) (io.Writer, error) {
	keysize := aesKeyLen(aesstrength)
	salt := make([]byte, keysize/2)
	_, err := rand.Read(salt[:])
	if err != nil {
		return nil, errors.New("zip: unable to generate random salt")
	}
	ekey, akey, pwv := generateKeys(password(), salt[:], keysize)
	fw.hmac = hmac.New(sha1.New, akey)
	aw := &authWriter{
		hmac: fw.hmac,
		w:    w,
	}
	es, err := encryptStream(ekey, aw)
	if err != nil {
		return nil, err
	}
	ew := &encryptionWriter{
		pwv:   pwv,
		salt:  salt[:],
		w:     w,
		es:    es,
		first: true,
	}
	return ew, nil
}

// IsEncrypted indicates whether this file's data is encrypted.
func (h *FileHeader) IsEncrypted() bool {
	return h.Flags&0x1 == 1
}

// WinZip AE-2 specifies that no CRC value is written and
// should be skipped when reading.
func (h *FileHeader) isAE2() bool {
	return h.ae == 2
}

func (h *FileHeader) writeWinZipExtra() {
	// total size is 11 bytes
	var buf [11]byte
	eb := writeBuf(buf[:])
	eb.uint16(winzipAesExtraId) // 0x9901
	eb.uint16(7)                // following data size is 7
	eb.uint16(2)                // ae 2
	eb.uint16(0x4541)           // "AE"
	eb.uint8(h.aesStrength)     // aes256
	eb.uint16(h.Method)         // original compression method
	h.Extra = append(h.Extra, buf[:]...)
}

func (h *FileHeader) setEncryptionMethod(enc EncryptionMethod) {
	h.encryption = enc
	switch enc {
	case AES128Encryption:
		h.aesStrength = 1
	case AES192Encryption:
		h.aesStrength = 2
	case AES256Encryption:
		h.aesStrength = 3
	}
}

func (h *FileHeader) setEncryptionBit() {
	h.Flags |= 0x1
}

// SetPassword sets the password used for encryption/decryption.
func (h *FileHeader) SetPassword(password string) {
	if !h.IsEncrypted() {
		h.setEncryptionBit()
	}
	h.password = func() []byte {
		return []byte(password)
	}
}

// PasswordFn is a function that returns the password
// as a byte slice
type passwordFn func() []byte

// Encrypt adds a file to the zip file using the provided name.
// It returns a Writer to which the file contents should be written. File
// contents will be encrypted with AES-256 using the given password. The
// file's contents must be written to the io.Writer before the next call
// to Create, CreateHeader, or Close.
func (w *Writer) Encrypt(name string, password string, enc EncryptionMethod) (io.Writer, error) {
	fh := &FileHeader{
		Name:   name,
		Method: Deflate,
	}
	fh.SetModTime(time.Now())
	fh.SetPassword(password)
	fh.setEncryptionMethod(enc)
	return w.CreateHeader(fh)
}
