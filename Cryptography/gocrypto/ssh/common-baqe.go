package ssh

const (
	Cipher_aes128_ctr   = "aes128-ctr"
	Cipher_aes192_ctr   = "aes192-ctr"
	Cipher_aes256_ctr   = "aes256-ctr"
	Cipher_aes128_gcm   = "aes128-gcm@openssh.com"
	Cipher_chacha       = chacha20Poly1305ID
	Cipher_arcfour256   = "arcfour256"
	Cipher_arcfour128   = "arcfour128"
	Cipher_arcfour      = "arcfour"
	Cipher_aes128cbc    = aes128cbcID
	Cipher_tripledescbc = tripledescbcID
)

// hash-based message authentication code
const (
	HMAC_sha2_256_etm = "hmac-sha2-256-etm@openssh.com"
	HMAC_sha2_256     = "hmac-sha2-256"
	HMAC_sha1         = "hmac-sha1"
	HMAC_sha1_96      = "hmac-sha1-96"
)
