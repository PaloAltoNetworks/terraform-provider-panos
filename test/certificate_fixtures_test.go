package provider_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"time"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// The certificate and private-key material exercised by the certificate_import,
// ssl_decrypt and ssl_tls_service_profile acceptance tests is generated fresh at
// process startup instead of being committed to the repository. Committing
// private keys - even for throwaway self-signed test certificates - trips secret
// scanners and invites accidental reuse of the material in production. These
// values are ephemeral: a new pair is minted every time the test binary runs.
//
// Only the two formats PAN-OS actually needs are produced here: unencrypted
// PKCS#8 PEM and a passphrase-protected PKCS#12 bundle. The passphrase-protected
// PEM (PKCS#8 encrypted) path is not generated because Go's standard library
// cannot emit that format; the passphrase code path remains covered by the
// PKCS#12 test cases below.
const testCertPassphrase = "paloalto"

// These package-level variables replace the former hardcoded PEM/PKCS#12
// constants. They are consumed across the certificate_import, ssl_decrypt,
// ssl_tls_service_profile and saml_idp_profile acceptance tests (all in
// package provider_test).
var (
	certPemInitial       string
	privateKeyPemInitial string
	certPemUpdated       string
	privateKeyPemUpdated string
	certKeyPkcs12Initial string
	certKeyPkcs12Updated string
)

func init() {
	initial := mustGenerateSelfSignedFixture("initial.example.test")
	updated := mustGenerateSelfSignedFixture("updated.example.test")

	certPemInitial = initial.certPEM
	privateKeyPemInitial = initial.keyPEM
	certKeyPkcs12Initial = initial.pkcs12Base64

	certPemUpdated = updated.certPEM
	privateKeyPemUpdated = updated.keyPEM
	certKeyPkcs12Updated = updated.pkcs12Base64
}

type certFixture struct {
	certPEM      string
	keyPEM       string
	pkcs12Base64 string
}

func mustGenerateSelfSignedFixture(commonName string) certFixture {
	fixture, err := generateSelfSignedFixture(commonName)
	if err != nil {
		panic(fmt.Sprintf("generating certificate fixture for %q: %v", commonName, err))
	}
	return fixture
}

// generateSelfSignedFixture mints a fresh RSA-4096 key and a matching self-signed
// CA certificate, then returns them as unencrypted PKCS#8 PEM plus a
// passphrase-protected PKCS#12 bundle. The trailing newline emitted by
// pem.EncodeToMemory is trimmed so the values match the exact form PAN-OS echoes
// back into Terraform state.
func generateSelfSignedFixture(commonName string) (certFixture, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return certFixture{}, fmt.Errorf("generating RSA key: %w", err)
	}

	serialLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, err := rand.Int(rand.Reader, serialLimit)
	if err != nil {
		return certFixture{}, fmt.Errorf("generating serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Example, Inc."},
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return certFixture{}, fmt.Errorf("creating certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return certFixture{}, fmt.Errorf("parsing generated certificate: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})

	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return certFixture{}, fmt.Errorf("marshaling private key: %w", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})

	pfx, err := pkcs12.Modern.Encode(key, cert, nil, testCertPassphrase)
	if err != nil {
		return certFixture{}, fmt.Errorf("encoding PKCS#12 bundle: %w", err)
	}

	return certFixture{
		certPEM:      strings.TrimRight(string(certPEM), "\n"),
		keyPEM:       strings.TrimRight(string(keyPEM), "\n"),
		pkcs12Base64: base64.StdEncoding.EncodeToString(pfx),
	}, nil
}
