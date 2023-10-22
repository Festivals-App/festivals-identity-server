package festivalspki

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/acme/autocert"
)

// LoadServerCertificates will either load the local server certificates or will try to load valid certificates from Lets Encrypt if there are no local certificate files.
func LoadServerCertificates(serverCert string, serverKey string, rootCACert string, certManager *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {

		certificate, err := tls.LoadX509KeyPair(serverCert, serverKey)
		if err != nil {
			return certManager.GetCertificate(hello)
		}
		rootCACert, err := LoadX509Certificate(rootCACert)
		if err != nil {
			log.Panic().Err(err).Str("type", "server").Msg("Failed to load FestivalsApp Root CA certificate")
		}
		certificate.Certificate = append(certificate.Certificate, rootCACert.Raw)
		return &certificate, err
	}
}

// LoadX509Certificate reads and parses a certificate from a .crt file.
// The file must contain PEM encoded data. The certificate file may only contain one certificate.
func LoadX509Certificate(certFile string) (*x509.Certificate, error) {

	certContent, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(certContent)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}
