package mediasoupgo

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
)

// Clone clones the given value.
func Clone[T any](value T) T {
    var zero T
    if any(value) == nil {
        return zero
    }

    // Go doesn't have a direct equivalent of Number.isNaN, but we can handle NaN for float types
    // Since Go doesn't allow direct NaN checks for non-float types, we skip this check

    // Try to use JSON serialization as a fallback for deep cloning
    data, err := json.Marshal(value)
    if err != nil {
        return zero
    }

    var cloned T
    err = json.Unmarshal(data, &cloned)
    if err != nil {
        return zero
    }

    return cloned
}

// GenerateUUIDv4 generates a random UUID v4.
func GenerateUUIDv4() string {
    uuid := make([]byte, 16)
    _, err := rand.Read(uuid)
    if err != nil {
        return ""
    }

    // Set version (4) and variant bits
    uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
    uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant RFC 4122

    // Format as a UUID string
    return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// GenerateRandomNumber generates a random positive integer between 100000000 and 999999999.
func GenerateRandomNumber() int {
    const min, max = 100000000, 999999999
    n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
    if err != nil {
        return 0
    }
    return int(n.Int64()) + min
}

// DeepFreeze in Go cannot enforce immutability like JavaScript's Object.freeze.
// This is a no-op in Go since immutability must be handled by convention or runtime checks.
func DeepFreeze[T any](object T) T {
    return object
}
