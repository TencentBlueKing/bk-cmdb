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
func NewTLSConfigFromConf(cfg *TLSClientConfig) (*tls.Config, error) {
	// createTLSConfig creates tls.Config based on TLSConfig.
	// It handles one-way and mutual TLS authentication, and TLS disabling.
	var tlsConf *tls.Config = nil // initialize tlsConf to nil, which means TLS is disabled by default

	if len(cfg.CAFile) != 0 { // if CAFile is configured, then enable TLS
		var err error
		if cfg.CertFile != "" && cfg.KeyFile != "" {
			// if CertFile and KeyFile are both configured, then use mutual TLS authentication
			tlsConf, err = ClientTLSConfVerity(cfg.CAFile, cfg.CertFile,
				cfg.KeyFile, "")
		} else {
			// otherwise, only CAFile is configured, use one-way TLS authentication, only verify server certificate
			tlsConf, err = ClientTslConfVerityServer(cfg.CAFile)
		}
		if err != nil {
			return nil, err
		}
		if tlsConf != nil {
			tlsConf.InsecureSkipVerify = cfg.InsecureSkipVerify
		}
	}
	// if cfg.TLSConfig.CAFile is empty, then tlsConf remains nil, which means TLS is disabled
	return tlsConf, nil
}
