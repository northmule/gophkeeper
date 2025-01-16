package util

import (
	"testing"
)

func TestPasswordHashSha256(t *testing.T) {
	tests := []struct {
		password string
		expected string
	}{
		{"password123", "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"averylongpasswordthatshouldtestthelengthoftheinputandhowthehashfunctionhandlesit", "750fef5b6088007348d0af24d60d8efad3e1a37c442e25a4e1738e302362e732"},
	}

	for _, test := range tests {
		result := PasswordHashSha256(test.password)
		if result != test.expected {
			t.Errorf("PasswordHashSha256(%q) = %q, want %q", test.password, result, test.expected)
		}
	}
}

func TestPasswordHashSha512(t *testing.T) {
	tests := []struct {
		password string
		expected string
	}{
		{"password123", "bed4efa1d4fdbd954bd3705d6a2a78270ec9a52ecfbfb010c61862af5c76af1761ffeb1aef6aca1bf5d02b3781aa854fabd2b69c790de74e17ecfec3cb6ac4bf"},
		{"", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
		{"averylongpasswordthatshouldtestthelengthoftheinputandhowthehashfunctionhandlesit", "5b114577bd6bc2fe35d14865eb82fb58a81a146eadaa0548d56eb46026414f65fcecfe9d08034ebd75b88b7e79cba39979534f0bbf5e9e7bc7b85ddbe5bfe104"},
	}

	for _, test := range tests {
		result := PasswordHashSha512(test.password)
		if result != test.expected {
			t.Errorf("PasswordHashSha512(%q) = %q, want %q", test.password, result, test.expected)
		}
	}
}
