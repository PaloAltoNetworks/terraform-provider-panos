package ike

// Valid Entryption values.
const (
	EncryptionDes       = "des"
	Encryption3des      = "3des"
	EncryptionAes128    = "aes-128-cbc"
	EncryptionAes192    = "aes-192-cbc"
	EncryptionAes256    = "aes-256-cbc"
	EncryptionAes128Gcm = "aes-128-gcm"
	EncryptionAes256Gcm = "aes-256-gcm"
)

// Valid Time values.
const (
	TimeSeconds = "seconds"
	TimeMinutes = "minutes"
	TimeHours   = "hours"
	TimeDays    = "days"
)

const (
	singular = "ike crypto profile"
	plural   = "ike crypto profiles"
)
