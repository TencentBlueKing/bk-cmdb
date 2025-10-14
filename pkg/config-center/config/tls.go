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

package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// TLSConfig is the common TLS configuration.
type TLSConfig struct {
	// InsecureSkipVerify defines whether the tls certificate should be verified.
	InsecureSkipVerify bool
	// CAFile is the trusted root CA certificate bundle file path.
	CAFile string
	// CertFile is the certificate file path.
	CertFile string
	// KeyFile is the private key file path.
	KeyFile string
}

// ToClientConf converts the TLSConfig to a tls.Config for client-side connections.
func (c *TLSConfig) ToClientConf() (*tls.Config, bool, error) {
	// tls required config is not set, returns a default skip verify config and disable TLS flag
	if c == nil || c.CAFile == "" {
		return &tls.Config{InsecureSkipVerify: true}, false, nil // nolint:gosec
	}

	// generate tls config
	tlsConf := &tls.Config{
		InsecureSkipVerify: c.InsecureSkipVerify, // nolint:gosec
	}

	// load ca file and set Root CAs for server verification
	caPool, err := loadCa(c.CAFile)
	if err != nil {
		return nil, false, fmt.Errorf("load CA file %s failed, err: %v", c.CAFile, err)
	}
	tlsConf.RootCAs = caPool

	// if mutual TLS is requested, load client certificate and key
	if c.CertFile != "" && c.KeyFile != "" {
		cert, err := loadCertificates(c.CertFile, c.KeyFile)
		if err != nil {
			return nil, false, fmt.Errorf("load certificate: %s, key: %s failed, err: %v", c.CertFile, c.KeyFile, err)
		}
		tlsConf.Certificates = []tls.Certificate{*cert}
	}

	return tlsConf, true, nil
}

// ToServerConf converts the TLSConfig to a tls.Config for server-side connections.
func (c *TLSConfig) ToServerConf() (*tls.Config, bool, error) {
	// tls required config is not set, returns a default skip verify config and disable TLS flag
	if c == nil || c.CAFile == "" || c.CertFile == "" || c.KeyFile == "" {
		return &tls.Config{InsecureSkipVerify: true}, false, nil // nolint:gosec
	}

	// load ca file
	caPool, err := loadCa(c.CAFile)
	if err != nil {
		return nil, false, fmt.Errorf("load CA file %s failed, err: %v", c.CAFile, err)
	}

	// load certificate and key file
	cert, err := loadCertificates(c.CertFile, c.KeyFile)
	if err != nil {
		return nil, false, fmt.Errorf("load certificate: %s, key: %s failed, err: %v", c.CertFile, c.KeyFile, err)
	}

	return &tls.Config{
		InsecureSkipVerify: c.InsecureSkipVerify, // nolint:gosec
		ClientCAs:          caPool,
		Certificates:       []tls.Certificate{*cert},
		ClientAuth:         tls.RequireAndVerifyClientCert,
	}, true, nil
}

// loadCa loads a CA bundle from the given file.
func loadCa(caFile string) (*x509.CertPool, error) {
	ca, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("read CA file failed, err: %v", err)
	}

	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("append CA certs failed")
	}

	return caPool, nil
}

// loadCertificates loads a TLS certificate and its corresponding private key from PEM files.
func loadCertificates(certFile, keyFile string) (*tls.Certificate, error) {
	// read key file
	priKey, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("read key file: %w", err)
	}

	// read certificate file
	certData, err := os.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("read cert file: %w", err)
	}

	// parse certificate into X509 key pair
	tlsCert, err := tls.X509KeyPair(certData, priKey)
	if err != nil {
		return nil, fmt.Errorf("parse X509 key pair failed: %w", err)
	}

	return &tlsCert, nil
}
