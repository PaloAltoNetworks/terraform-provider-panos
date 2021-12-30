package decryption

// Valid Action values.
//
// Decrypt and forward is PAN-OS 8.1+.
const (
	ActionNoDecrypt         = "no-decrypt"
	ActionDecrypt           = "decrypt"
	ActionDecryptAndForward = "decrypt-and-forward"
)

// Valid DecryptionType values.
const (
	DecryptionTypeSslForwardProxy      = "ssl-forward-proxy"
	DecryptionTypeSshProxy             = "ssh-proxy"
	DecryptionTypeSslInboundInspection = "ssl-inbound-inspection"
)

const (
	singular = "decryption rule"
	plural   = "decryption rules"
)
