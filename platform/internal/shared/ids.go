package shared

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

type RandomIDGenerator struct{}

func (RandomIDGenerator) NewID(prefix string) (string, error) {
	secret, err := randomBase64(16)
	if err != nil {
		return "", err
	}
	if prefix == "" {
		return secret, nil
	}
	return fmt.Sprintf("%s_%s", prefix, secret), nil
}

func (RandomIDGenerator) NewSecret(bytes int) (string, error) {
	return randomBase64(bytes)
}

func randomBase64(bytes int) (string, error) {
	buf := make([]byte, bytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}
