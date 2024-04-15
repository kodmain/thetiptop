package hash_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	testCases := []struct {
		data      string
		algo      hash.HashAlgo
		expected  string
		expectErr bool
	}{
		{
			data:      "password123",
			algo:      hash.SHA1,
			expected:  "cbfdac6008f9cab4083784cbd1874f76618d2a97",
			expectErr: false,
		},
		{
			data:      "password123",
			algo:      hash.SHA256,
			expected:  "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
			expectErr: false,
		},
		{
			data:      "password123",
			algo:      hash.SHA512,
			expected:  "bed4efa1d4fdbd954bd3705d6a2a78270ec9a52ecfbfb010c61862af5c76af1761ffeb1aef6aca1bf5d02b3781aa854fabd2b69c790de74e17ecfec3cb6ac4bf",
			expectErr: false,
		},
		{
			data:      "password123",
			algo:      hash.MD5,
			expected:  "482c811da5d5b4bc6d497ffa98491e38",
			expectErr: false,
		},
		{
			data:      "password123",
			algo:      hash.BCRYPT,
			expected:  "", // Update with the expected bcrypt hash value
			expectErr: false,
		},
		{
			data:      "password123",
			algo:      99,
			expected:  "", // Update with the expected bcrypt hash value
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		hashedData, err := hash.Hash(tc.data, tc.algo)
		if tc.expectErr && err == nil {
			t.Errorf("Expected an error but got nil for data: %s, algo: %d", tc.data, tc.algo)
		}

		if !tc.expectErr && err != nil {
			t.Errorf("Expected no error but got: %v for data: %s, algo: %d", err, tc.data, tc.algo)
		}

		if hashedData != tc.expected && tc.algo != hash.BCRYPT {
			t.Errorf("Expected hash: %s but got: %s for data: %s, algo: %d", tc.expected, hashedData, tc.data, tc.algo)
		} else {
			err := hash.CompareHash(hashedData, tc.data, tc.algo)
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		}
	}
}
