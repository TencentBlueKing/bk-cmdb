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

// Verify checks if the TLS client configuration is valid
func (cfg *TLSClientConfig) Verify() bool {
	if cfg == nil {
		return false
	}

	// If CAFile is configured, it's considered valid
	if len(cfg.CAFile) > 0 {
		return true
	}

	// If both CertFile and KeyFile are configured, it's considered valid
	if len(cfg.CertFile) > 0 && len(cfg.KeyFile) > 0 {
		return true
	}

	return false
}

// NewTLSConfigFromConf creates a new TLS configuration from TLSClientConfig
func NewTLSConfigFromConf(cfg *TLSClientConfig) (*tls.Config, error) {
	// createTLSConfig creates tls.Config based on TLSConfig.
	// It handles one-way and mutual TLS authentication, and TLS disabling.
	tlsConf := &tls.Config{}

	if cfg != nil && len(cfg.CAFile) != 0 { // if CAFile is configured, then enable TLS
		var err error
		if len(cfg.CertFile) != 0 && len(cfg.KeyFile) != 0 {
			// if CertFile and KeyFile are both configured, then use mutual TLS authentication
			tlsConf, err = ClientTLSConfVerity(cfg.CAFile, cfg.CertFile, cfg.KeyFile, "")
		} else {
			// otherwise, only CAFile is configured, use one-way TLS authentication, only verify server certificate
			tlsConf, err = ClientTslConfVerityServer(cfg.CAFile)
		}
		if err != nil {
			return nil, err
		}
	}
	tlsConf.InsecureSkipVerify = cfg.InsecureSkipVerify
	return tlsConf, nil
}
