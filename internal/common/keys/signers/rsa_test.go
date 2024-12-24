package signers

import (
	"crypto/rsa"
	"testing"
)

func TestRsaSigner_GenerateKey(t *testing.T) {
	signer := NewRsaSigner()

	key, err := signer.GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}
	if key == nil {
		t.Errorf("Generated key is nil")
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		t.Errorf("Generated key is not of type *rsa.PrivateKey")
	}
	if rsaKey.PublicKey.N.BitLen() != 4096 {
		t.Errorf("Generated key does not have the correct bit length: %d != 4096", rsaKey.PublicKey.N.BitLen())
	}

	keys := make(map[string]struct{})
	for i := 0; i < 5; i++ {
		key, err := signer.GenerateKey()
		if err != nil {
			t.Fatalf("GenerateKey failed: %v", err)
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			t.Errorf("Generated key is not of type *rsa.PrivateKey")
		}
		keyStr := rsaKey.PublicKey.N.String()
		if _, exists := keys[keyStr]; exists {
			t.Errorf("Generated key is not unique")
		}
		keys[keyStr] = struct{}{}
	}
}
