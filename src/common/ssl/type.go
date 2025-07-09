/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package ssl

import "crypto/tls"

// TLSClientConfig Common TLS client configuration
type TLSClientConfig struct {
	// Server should be accessed without verifying the TLS certificate. For testing only.
	InsecureSkipVerify bool
	// Server requires TLS client certificate authentication
	CertFile string
	// Server requires TLS client certificate authentication
	KeyFile string
	// Trusted root certificates for server
	CAFile string
	// the password to decrypt the certificate
	Password string
}

// NewTLSConfigFromConf creates a new TLS configuration from TLSClientConfig
// Returns:
// - *tls.Config: TLS configuration
// - bool: whether TLS is enabled
// - error: any error occurred during configuration
func NewTLSConfigFromConf(cfg *TLSClientConfig) (*tls.Config, bool, error) {
	// createTLSConfig creates tls.Config based on TLSConfig.
	// It handles one-way and mutual TLS authentication, and TLS disabling.
	tlsConf := &tls.Config{}

	if cfg == nil {
		return tlsConf, false, nil
	}

	if cfg != nil && len(cfg.CAFile) != 0 { // if CAFile is configured, then enable TLS
		var err error
		if len(cfg.CertFile) != 0 && len(cfg.KeyFile) != 0 {
			// if CertFile and KeyFile are both configured, then use mutual TLS authentication
			tlsConf, err = ClientTLSConfVerity(cfg.CAFile, cfg.CertFile, cfg.KeyFile, cfg.Password)
		} else {
			// otherwise, only CAFile is configured, use one-way TLS authentication, only verify server certificate
			tlsConf, err = ClientTslConfVerityServer(cfg.CAFile)
		}
		if err != nil {
			return tlsConf, false, err
		}
		tlsConf.InsecureSkipVerify = cfg.InsecureSkipVerify
		return tlsConf, true, nil
	}

	return tlsConf, false, nil
}
