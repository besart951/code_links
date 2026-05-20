package tokenjose

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/besart951/code_links/platform/internal/access"
)

type StaticRegistry struct {
	Signing    SigningKey
	Encryption map[access.ProductKey]AudienceEncryptionKey
	Decryption map[access.ProductKey]AudienceDecryptionKey
}

func (r StaticRegistry) SigningKey(ctx context.Context) (SigningKey, error) {
	if len(r.Signing.PrivateKey) == 0 {
		return SigningKey{}, ErrKeyNotFound
	}
	return r.Signing, nil
}

func (r StaticRegistry) SigningPublicKey(ctx context.Context, keyID string) (ed25519.PublicKey, error) {
	if r.Signing.KeyID != keyID || len(r.Signing.PublicKey) == 0 {
		return nil, ErrKeyNotFound
	}
	return r.Signing.PublicKey, nil
}

func (r StaticRegistry) AudienceEncryptionKey(ctx context.Context, audience access.ProductKey) (AudienceEncryptionKey, error) {
	key, ok := r.Encryption[audience]
	if !ok || key.PublicKey == nil {
		return AudienceEncryptionKey{}, ErrKeyNotFound
	}
	return key, nil
}

func (r StaticRegistry) AudienceDecryptionKey(ctx context.Context, audience access.ProductKey, keyID string) (AudienceDecryptionKey, error) {
	key, ok := r.Decryption[audience]
	if !ok || key.KeyID != keyID || key.PrivateKey == nil {
		return AudienceDecryptionKey{}, ErrKeyNotFound
	}
	return key, nil
}

func LoadEd25519PrivateKey(path string) (ed25519.PrivateKey, error) {
	der, err := readPEM(path)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, err
	}
	privateKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("expected ed25519 private key in %s", path)
	}
	return privateKey, nil
}

func LoadECDSAPublicKey(path string) (*ecdsa.PublicKey, error) {
	der, err := readPEM(path)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, err
	}
	publicKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("expected ecdsa public key in %s", path)
	}
	return publicKey, nil
}

func LoadECDSAPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	der, err := readPEM(path)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, err
	}
	privateKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("expected ecdsa private key in %s", path)
	}
	return privateKey, nil
}

func readPEM(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("decode pem %s", path)
	}
	return block.Bytes, nil
}
