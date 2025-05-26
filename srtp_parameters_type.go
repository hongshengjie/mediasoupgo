package mediasoupgo

type SrtpParameters struct {
	CryptoSuite SrtpCryptoSuite
	KeyBase64   string
}

type SrtpCryptoSuite string

const (
	AEADAES256GCMSrtpCryptoSuite      SrtpCryptoSuite = "AEAD_AES_256_GCM"
	AEADAES128GCMSrtpCryptoSuite      SrtpCryptoSuite = "AEAD_AES_128_GCM"
	AESCM128HMACSHA180SrtpCryptoSuite SrtpCryptoSuite = "AES_CM_128_HMAC_SHA1_80"
	AESCM128HMACSHA132SrtpCryptoSuite SrtpCryptoSuite = "AES_CM_128_HMAC_SHA1_32"
)
