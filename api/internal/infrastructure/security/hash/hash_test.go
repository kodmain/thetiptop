package hash_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"github.com/stretchr/testify/assert"
)

// TestHash tests different hashing algorithms
// This test verifies the correctness of hashing algorithms against known expected values.
//
// Parameters:
// - t: *testing.T the test framework
//
// Returns:
// - None
func TestHash(t *testing.T) {
	testCases := []struct {
		data      *string
		algo      hash.HashAlgo
		expected  *string
		expectErr bool
	}{
		{
			data:      aws.String("password123"),
			algo:      hash.SHA1,
			expected:  aws.String("cbfdac6008f9cab4083784cbd1874f76618d2a97"),
			expectErr: false,
		},
		{
			data:      aws.String("password123"),
			algo:      hash.SHA256,
			expected:  aws.String("ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f"),
			expectErr: false,
		},
		{
			data:      aws.String("password123"),
			algo:      hash.SHA512,
			expected:  aws.String("bed4efa1d4fdbd954bd3705d6a2a78270ec9a52ecfbfb010c61862af5c76af1761ffeb1aef6aca1bf5d02b3781aa854fabd2b69c790de74e17ecfec3cb6ac4bf"),
			expectErr: false,
		},
		{
			data:      aws.String("password123"),
			algo:      hash.MD5,
			expected:  aws.String("482c811da5d5b4bc6d497ffa98491e38"),
			expectErr: false,
		},
		{
			data:      aws.String("password123"),
			algo:      hash.BCRYPT,
			expected:  nil, // Update with the expected bcrypt hash value
			expectErr: false,
		},
		{
			data:      aws.String("password123"),
			algo:      99,
			expected:  aws.String("cbfdac6008f9cab4083784cbd1874f76618d2a97"),
			expectErr: true,
		},
		{
			data:      nil,
			algo:      99,
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		hashedData, err := hash.Hash(tc.data, tc.algo)
		if tc.expectErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			if tc.algo != hash.BCRYPT {
				assert.Equal(t, tc.expected, hashedData, "Expected hash: %s but got: %s for data: %s, algo: %d", *tc.expected, *hashedData, tc.data, tc.algo)
			}
		}

		err = hash.CompareHash(hashedData, tc.data, tc.algo)
		if tc.expectErr {
			assert.Error(t, err)
		}
	}
}
