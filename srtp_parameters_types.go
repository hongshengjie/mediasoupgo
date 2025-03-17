package mediasoupgo

// SrtpParameters represents SRTP parameters.
type SrtpParameters struct {
	// Encryption and authentication transforms to be used.
	CryptoSuite SrtpCryptoSuite `json:"cryptoSuite"`

	// SRTP keying material (master key and salt) in Base64.
	KeyBase64 string `json:"keyBase64"`
}

// SrtpCryptoSuite represents SRTP crypto suite options.
type SrtpCryptoSuite string

const (
	AEAD_AES_256_GCM                SrtpCryptoSuite = "AEAD_AES_256_GCM"
	AEAD_AES_128_GCM                SrtpCryptoSuite = "AEAD_AES_128_GCM"
	AES_CM_128_HMAC_SHA1_80         SrtpCryptoSuite = "AES_CM_128_HMAC_SHA1_80"
	AES_CM_128_HMAC_SHA1_32         SrtpCryptoSuite = "AES_CM_128_HMAC_SHA1_32"
)
