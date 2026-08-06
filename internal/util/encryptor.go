package util

// Encryptor abstracts the PII encrypt/decrypt boundary Java implements via
// EncryptionDecryptionUtil against a proprietary ABAC-aware encryption
// microservice. Per agreed scope, sw-services-go stubs this behind an
// interface with a no-op (plaintext) implementation rather than depending on
// that unavailable service — swapping in a real implementation later touches
// this file only, not call sites.
type Encryptor interface {
	Encrypt(plaintext string, key string) (string, error)
	Decrypt(ciphertext string, key string) (string, error)
}

type NoOpEncryptor struct{}

func NewNoOpEncryptor() *NoOpEncryptor { return &NoOpEncryptor{} }

func (NoOpEncryptor) Encrypt(plaintext string, _ string) (string, error)  { return plaintext, nil }
func (NoOpEncryptor) Decrypt(ciphertext string, _ string) (string, error) { return ciphertext, nil }
