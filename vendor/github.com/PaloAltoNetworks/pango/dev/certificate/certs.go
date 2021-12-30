package certificate

// Pem is a PEM certificate suitable to be imported into PAN-OS.
//
// Importing the certificate and the private key are two separate API
// calls.  If the PrivateKey is left unspecified, then the 2nd API call
// will not be made.
type Pem struct {
	Name                string
	Certificate         string
	CertificateFilename string
	PrivateKey          string
	PrivateKeyFilename  string
	Passphrase          string
}

// Pkcs12 is a PKCS12 certificate suitable to be imported into PAN-OS.
type Pkcs12 struct {
	Name                string
	Certificate         string
	CertificateFilename string
	Passphrase          string
}
