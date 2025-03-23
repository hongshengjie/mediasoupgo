package mediasoupgo

// SrtpCryptoSuite represents the SRTP crypto suite.
type SrtpCryptoSuite string

// Constants for SrtpCryptoSuite values.
const (
	SrtpCryptoSuiteAEADAES256GCM      SrtpCryptoSuite = "AEAD_AES_256_GCM"
	SrtpCryptoSuiteAEADAES128GCM      SrtpCryptoSuite = "AEAD_AES_128_GCM"
	SrtpCryptoSuiteAESCM128HMACSHA180 SrtpCryptoSuite = "AES_CM_128_HMAC_SHA1_80"
	SrtpCryptoSuiteAESCM128HMACSHA132 SrtpCryptoSuite = "AES_CM_128_HMAC_SHA1_32"
)

// SrtpParameters represents SRTP parameters.
type SrtpParameters struct {
	// Encryption and authentication transforms to be used.
	CryptoSuite SrtpCryptoSuite
	// SRTP keying material (master key and salt) in Base64.
	KeyBase64 string
}
