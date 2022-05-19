package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte { // interface를 인자로 하면 어떤 타입이라도 받을 수 있다.
	var aBuffer bytes.Buffer            // where we put ourt bytes
	encoder := gob.NewEncoder(&aBuffer) // gob is package for encoding/decoding bytes
	HandleErr(encoder.Encode(i))        // encode whole block then put the result in the blockBuffer
	return aBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(decoder.Decode(i))
}

func Hash(i interface{}) string {
	s := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func Splitter(s, sep string, i int) string {
	r := strings.Split(s, sep)
	if len(r)-1 < i {
		return ""
	}
	return r[i]
}
