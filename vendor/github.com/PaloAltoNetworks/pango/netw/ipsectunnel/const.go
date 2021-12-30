package ipsectunnel

const (
	TypeAutoKey                = "auto-key"
	TypeManualKey              = "manual-key"
	TypeGlobalProtectSatellite = "global-protect-satellite"
)

const (
	MkEspEncryptionDes    = "des"
	MkEspEncryption3des   = "3des"
	MkEspEncryptionAes128 = "aes-128-cbc"
	MkEspEncryptionAes192 = "aes-192-cbc"
	MkEspEncryptionAes256 = "aes-256-cbc"
	MkEspEncryptionNull   = "null"
)

const (
	MkProtocolEsp = "esp"
	MkProtocolAh  = "ah"
)

const (
	MkAuthTypeMd5    = "md5"
	MkAuthTypeSha1   = "sha1"
	MkAuthTypeSha256 = "sha256"
	MkAuthTypeSha384 = "sha384"
	MkAuthTypeSha512 = "sha512"
	MkAuthTypeNone   = "none"
)

const (
	singular = "ipsec tunnel"
	plural   = "ipsec tunnels"
)
