/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package config_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
)

func prepareTLSFiles(t *testing.T) (string, string, string) {
	// generate a self-signed CA certificate and key
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate CA key: %v", err)
	}
	serialNumber, _ := rand.Int(rand.Reader, big.NewInt(1<<62))
	caTemplate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "Test CA",
			Organization: []string{"Test Org"},
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create CA cert: %v", err)
	}

	caCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})
	tempDir := t.TempDir()
	caPath := writeTempFile(t, tempDir, "ca.pem", caCertPEM)

	// generate a leaf certificate signed by given CA
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate leaf key: %v", err)
	}
	serialNumber, _ = rand.Int(rand.Reader, big.NewInt(1<<62))
	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "Test Cert",
			Organization: []string{"Test Org"},
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, caTemplate, &key.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create leaf cert: %v", err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	certPath := writeTempFile(t, tempDir, "client.crt", certPEM)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	keyPath := writeTempFile(t, tempDir, "client.key", keyPEM)

	return caPath, certPath, keyPath
}

func writeTempFile(t *testing.T, dir, name string, data []byte) string {
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("write file %s failed, err: %v", path, err)
	}
	return path
}

func TestTLSConfig(t *testing.T) {
	caPath, certPath, keyPath := prepareTLSFiles(t)

	// test skip verify tls config
	conf := &config.TLSConfig{
		InsecureSkipVerify: true,
	}
	_, enabled, err := conf.ToClientConf()
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, enabled)

	// test one way tls config
	conf = &config.TLSConfig{
		CAFile: caPath,
	}
	clientConf, enabled, err := conf.ToClientConf()
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, enabled)
	assert.False(t, clientConf.InsecureSkipVerify)
	assert.NotNil(t, clientConf.RootCAs)

	_, enabled, err = conf.ToServerConf()
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, enabled)

	// test mutual tls config
	conf = &config.TLSConfig{
		CAFile:   caPath,
		CertFile: certPath,
		KeyFile:  keyPath,
	}
	clientConf, enabled, err = conf.ToClientConf()
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, enabled)
	assert.False(t, clientConf.InsecureSkipVerify)
	assert.NotNil(t, clientConf.RootCAs)
	assert.Len(t, clientConf.Certificates, 1)

	serverConf, enabled, err := conf.ToServerConf()
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, enabled)
	assert.False(t, serverConf.InsecureSkipVerify)
	assert.NotNil(t, serverConf.ClientCAs)
	assert.Len(t, serverConf.Certificates, 1)
}
