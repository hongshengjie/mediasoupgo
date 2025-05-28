package mediasoupgo

import (
	"bytes"
	"encoding/gob"
)

func DeepCopyGob(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(&buf).Decode(dst)
}
