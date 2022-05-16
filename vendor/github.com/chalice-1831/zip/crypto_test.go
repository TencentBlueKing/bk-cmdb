package zip

import (
	"bytes"
	"io"
	"path/filepath"
	"testing"
)

// Test simple password reading.
func TestPasswordReadSimple(t *testing.T) {
	file := "hello-aes.zip"
	var buf bytes.Buffer
	r, err := OpenReader(filepath.Join("testdata", file))
	if err != nil {
		t.Errorf("Expected %s to open: %v.", file, err)
	}
	defer r.Close()
	if len(r.File) != 1 {
		t.Errorf("Expected %s to contain one file.", file)
	}
	f := r.File[0]
	if f.FileInfo().Name() != "hello.txt" {
		t.Errorf("Expected %s to have a file named hello.txt", file)
	}
	if f.Method != 0 {
		t.Errorf("Expected %s to have its Method set to 0.", file)
	}
	f.SetPassword("golang")
	rc, err := f.Open()
	if err != nil {
		t.Errorf("Expected to open the readcloser: %v.", err)
	}
	_, err = io.Copy(&buf, rc)
	if err != nil {
		t.Errorf("Expected to copy bytes: %v.", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("Hello World\r\n")) {
		t.Errorf("Expected contents were not found.")
	}
}

// Test for multi-file password protected zip.
// Each file can be protected with a different password.
func TestPasswordHelloWorldAes(t *testing.T) {
	file := "world-aes.zip"
	expecting := "helloworld"
	r, err := OpenReader(filepath.Join("testdata", file))
	if err != nil {
		t.Errorf("Expected %s to open: %v", file, err)
	}
	defer r.Close()
	if len(r.File) != 2 {
		t.Errorf("Expected %s to contain two files.", file)
	}
	var b bytes.Buffer
	for _, f := range r.File {
		if !f.IsEncrypted() {
			t.Errorf("Expected %s to be encrypted.", f.FileInfo().Name)
		}
		f.SetPassword("golang")
		rc, err := f.Open()
		if err != nil {
			t.Errorf("Expected to open readcloser: %v", err)
		}
		defer rc.Close()
		if _, err := io.Copy(&b, rc); err != nil {
			t.Errorf("Expected to copy bytes to buffer: %v", err)
		}
	}
	if !bytes.Equal([]byte(expecting), b.Bytes()) {
		t.Errorf("Expected ending content to be %s instead of %s", expecting, b.Bytes())
	}
}

// Test for password protected file that is larger than a single
// AES block size to check CTR implementation.
func TestPasswordMacbethAct1(t *testing.T) {
	file := "macbeth-act1.zip"
	expecting := "Exeunt"
	var b bytes.Buffer
	r, err := OpenReader(filepath.Join("testdata", file))
	if err != nil {
		t.Errorf("Expected %s to open: %v", file, err)
	}
	defer r.Close()
	for _, f := range r.File {
		if !f.IsEncrypted() {
			t.Errorf("Expected %s to be encrypted.", f.Name)
		}
		f.SetPassword("golang")
		rc, err := f.Open()
		if err != nil {
			t.Errorf("Expected to open readcloser: %v", err)
		}
		defer rc.Close()
		if _, err := io.Copy(&b, rc); err != nil {
			t.Errorf("Expected to copy bytes to buffer: %v", err)
		}
	}
	if !bytes.Contains(b.Bytes(), []byte(expecting)) {
		t.Errorf("Expected to find %s in the buffer %v", expecting, b.Bytes())
	}
}

// Change to AE-1 and change CRC value to fail check.
// Must be != 0 due to zip package already skipping if == 0.
func returnAE1BadCRC() (io.ReaderAt, int64) {
	return messWith("hello-aes.zip", func(b []byte) {
		// Change version to AE-1(1)
		b[0x2B] = 1 // file
		b[0xBA] = 1 // TOC
		// Change CRC to bad value
		b[0x11]++ // file
		b[0x6B]++ // TOC
	})
}

// Test for AE-1 Corrupt CRC
func TestPasswordAE1BadCRC(t *testing.T) {
	buf := new(bytes.Buffer)
	file, s := returnAE1BadCRC()
	r, err := NewReader(file, s)
	if err != nil {
		t.Errorf("Expected hello-aes.zip to open: %v", err)
	}
	for _, f := range r.File {
		if !f.IsEncrypted() {
			t.Errorf("Expected zip to be encrypted")
		}
		f.SetPassword("golang")
		rc, err := f.Open()
		if err != nil {
			t.Errorf("Expected the readcloser to open.")
		}
		defer rc.Close()
		if _, err := io.Copy(buf, rc); err != ErrChecksum {
			t.Errorf("Expected the checksum to fail")
		}
	}
}

// Corrupt the last byte of ciphertext to fail authentication
func returnTamperedData() (io.ReaderAt, int64) {
	return messWith("hello-aes.zip", func(b []byte) {
		b[0x50]++
	})
}

// Test for tampered file data payload.
func TestPasswordTamperedData(t *testing.T) {
	buf := new(bytes.Buffer)
	file, s := returnTamperedData()
	r, err := NewReader(file, s)
	if err != nil {
		t.Errorf("Expected hello-aes.zip to open: %v", err)
	}
	for _, f := range r.File {
		if !f.IsEncrypted() {
			t.Errorf("Expected zip to be encrypted")
		}
		f.SetPassword("golang")
		rc, err := f.Open()
		if err != nil {
			t.Errorf("Expected the readcloser to open.")
		}
		defer rc.Close()
		if _, err := io.Copy(buf, rc); err != ErrAuthentication {
			t.Errorf("Expected the checksum to fail")
		}
	}
}

func TestPasswordWriteSimple(t *testing.T) {
	contents := []byte("Hello World")
	conLen := len(contents)

	for _, enc := range []EncryptionMethod{StandardEncryption, AES128Encryption, AES192Encryption, AES256Encryption} {
		raw := new(bytes.Buffer)
		zipw := NewWriter(raw)
		w, err := zipw.Encrypt("hello.txt", "golang", enc)
		if err != nil {
			t.Errorf("Expected to create a new FileHeader")
		}
		n, err := io.Copy(w, bytes.NewReader(contents))
		if err != nil || n != int64(conLen) {
			t.Errorf("Expected to write the full contents to the writer.")
		}
		zipw.Close()

		// Read the zip
		buf := new(bytes.Buffer)
		zipr, err := NewReader(bytes.NewReader(raw.Bytes()), int64(raw.Len()))
		if err != nil {
			t.Errorf("Expected to open a new zip reader: %v", err)
		}
		nn := len(zipr.File)
		if nn != 1 {
			t.Errorf("Expected to have one file in the zip archive, but has %d files", nn)
		}
		z := zipr.File[0]
		z.SetPassword("golang")
		rr, err := z.Open()
		if err != nil {
			t.Errorf("Expected to open the readcloser: %v", err)
		}
		n, err = io.Copy(buf, rr)
		if err != nil {
			t.Errorf("Expected to write to temporary buffer: %v", err)
		}
		if n != int64(conLen) {
			t.Errorf("Expected to copy %d bytes to temp buffer, but copied %d bytes instead", conLen, n)
		}
		if !bytes.Equal(contents, buf.Bytes()) {
			t.Errorf("Expected the unzipped contents to equal '%s', but was '%s' instead", contents, buf.Bytes())
		}
	}
}

func TestZipCrypto(t *testing.T) {
	contents := []byte("Hello World")
	conLen := len(contents)

	raw := new(bytes.Buffer)
	zipw := NewWriter(raw)
	w, err := zipw.Encrypt("hello.txt", "golang", StandardEncryption)
	if err != nil {
		t.Errorf("Expected to create a new FileHeader")
	}
	n, err := io.Copy(w, bytes.NewReader(contents))
	if err != nil || n != int64(conLen) {
		t.Errorf("Expected to write the full contents to the writer.")
	}
	zipw.Close()

	zipr, _ := NewReader(bytes.NewReader(raw.Bytes()), int64(raw.Len()))
	zipr.File[0].SetPassword("golang")
	r, _ := zipr.File[0].Open()
	res := new(bytes.Buffer)
	io.Copy(res, r)
	r.Close()

	if !bytes.Equal(contents, res.Bytes()) {
		t.Errorf("Expected the unzipped contents to equal '%s', but was '%s' instead", contents, res.Bytes())
	}
}
