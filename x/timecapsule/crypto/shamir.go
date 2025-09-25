package crypto

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// ShamirSecretSharing implements Shamir's Secret Sharing algorithm
type ShamirSecretSharing struct {
	prime *big.Int
}

// NewShamirSecretSharing creates a new instance with a suitable prime
func NewShamirSecretSharing() *ShamirSecretSharing {
	// Using a 256-bit prime for security
	prime, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639747", 10)
	return &ShamirSecretSharing{
		prime: prime,
	}
}

// Share represents a single share of the secret
type Share struct {
	X *big.Int `json:"x"`
	Y *big.Int `json:"y"`
}

// SplitSecret splits a secret into n shares with threshold t
func (sss *ShamirSecretSharing) SplitSecret(secret []byte, threshold, totalShares int) ([]*Share, error) {
	if threshold > totalShares {
		return nil, fmt.Errorf("threshold (%d) cannot be greater than total shares (%d)", threshold, totalShares)
	}
	
	if threshold < 1 {
		return nil, fmt.Errorf("threshold must be at least 1")
	}
	
	if totalShares < 1 {
		return nil, fmt.Errorf("total shares must be at least 1")
	}
	
	// Convert secret to big integer
	secretInt := new(big.Int).SetBytes(secret)
	
	// Ensure secret is within the field
	if secretInt.Cmp(sss.prime) >= 0 {
		return nil, fmt.Errorf("secret is too large for the prime field")
	}
	
	// Generate random coefficients for polynomial of degree (threshold - 1)
	coefficients := make([]*big.Int, threshold)
	coefficients[0] = secretInt // The constant term is the secret
	
	for i := 1; i < threshold; i++ {
		coeff, err := rand.Int(rand.Reader, sss.prime)
		if err != nil {
			return nil, fmt.Errorf("failed to generate coefficient: %w", err)
		}
		coefficients[i] = coeff
	}
	
	// Generate shares by evaluating polynomial at different points
	shares := make([]*Share, totalShares)
	for i := 0; i < totalShares; i++ {
		x := big.NewInt(int64(i + 1)) // x coordinates start from 1
		y := sss.evaluatePolynomial(coefficients, x)
		
		shares[i] = &Share{
			X: new(big.Int).Set(x),
			Y: new(big.Int).Set(y),
		}
	}
	
	return shares, nil
}

// CombineShares reconstructs the secret from shares using Lagrange interpolation
func (sss *ShamirSecretSharing) CombineShares(shares []*Share) ([]byte, error) {
	if len(shares) < 1 {
		return nil, fmt.Errorf("need at least 1 share")
	}
	
	// Check for duplicate x coordinates
	xCoords := make(map[string]bool)
	for _, share := range shares {
		xStr := share.X.String()
		if xCoords[xStr] {
			return nil, fmt.Errorf("duplicate x coordinate found: %s", xStr)
		}
		xCoords[xStr] = true
	}
	
	// Perform Lagrange interpolation to find f(0) = secret
	secret := big.NewInt(0)
	
	for i, share := range shares {
		// Calculate Lagrange basis polynomial value at x=0
		numerator := big.NewInt(1)
		denominator := big.NewInt(1)
		
		for j, otherShare := range shares {
			if i != j {
				// numerator *= (0 - otherShare.X) = -otherShare.X
				temp := new(big.Int).Neg(otherShare.X)
				numerator.Mul(numerator, temp)
				numerator.Mod(numerator, sss.prime)
				
				// denominator *= (share.X - otherShare.X)
				temp = new(big.Int).Sub(share.X, otherShare.X)
				denominator.Mul(denominator, temp)
				denominator.Mod(denominator, sss.prime)
			}
		}
		
		// Calculate modular inverse of denominator
		denomInv := new(big.Int).ModInverse(denominator, sss.prime)
		if denomInv == nil {
			return nil, fmt.Errorf("failed to compute modular inverse")
		}
		
		// lagrange_i = numerator * denomInv
		lagrange := new(big.Int).Mul(numerator, denomInv)
		lagrange.Mod(lagrange, sss.prime)
		
		// secret += share.Y * lagrange_i
		term := new(big.Int).Mul(share.Y, lagrange)
		secret.Add(secret, term)
		secret.Mod(secret, sss.prime)
	}
	
	return secret.Bytes(), nil
}

// evaluatePolynomial evaluates the polynomial at a given x coordinate
func (sss *ShamirSecretSharing) evaluatePolynomial(coefficients []*big.Int, x *big.Int) *big.Int {
	result := big.NewInt(0)
	xPower := big.NewInt(1)
	
	for _, coeff := range coefficients {
		// result += coeff * x^i
		term := new(big.Int).Mul(coeff, xPower)
		result.Add(result, term)
		result.Mod(result, sss.prime)
		
		// x^i *= x for next iteration
		xPower.Mul(xPower, x)
		xPower.Mod(xPower, sss.prime)
	}
	
	return result
}

// ValidateShares validates that shares are properly formatted
func (sss *ShamirSecretSharing) ValidateShares(shares []*Share) error {
	if len(shares) == 0 {
		return fmt.Errorf("no shares provided")
	}
	
	for i, share := range shares {
		if share.X == nil || share.Y == nil {
			return fmt.Errorf("share %d has nil coordinates", i)
		}
		
		if share.X.Sign() <= 0 {
			return fmt.Errorf("share %d has invalid x coordinate (must be positive)", i)
		}
		
		if share.X.Cmp(sss.prime) >= 0 || share.Y.Cmp(sss.prime) >= 0 {
			return fmt.Errorf("share %d coordinates are outside the prime field", i)
		}
	}
	
	return nil
}

// EncryptShares encrypts each share with individual encryption keys
func (sss *ShamirSecretSharing) EncryptShares(shares []*Share, keys [][]byte) ([]*EncryptedShare, error) {
	if len(shares) != len(keys) {
		return nil, fmt.Errorf("number of shares (%d) must match number of keys (%d)", len(shares), len(keys))
	}
	
	encManager := NewEncryptionManager()
	encryptedShares := make([]*EncryptedShare, len(shares))
	
	for i, share := range shares {
		// Convert share to bytes for encryption
		shareData := sss.shareToBytes(share)
		
		// Encrypt the share
		encData, err := encManager.Encrypt(shareData, keys[i])
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt share %d: %w", i, err)
		}
		
		encryptedShares[i] = &EncryptedShare{
			Index:         i,
			EncryptedData: encData,
		}
	}
	
	return encryptedShares, nil
}

// EncryptedShare represents an encrypted Shamir share
type EncryptedShare struct {
	Index         int            `json:"index"`
	EncryptedData *EncryptedData `json:"encrypted_data"`
}

// shareToBytes converts a share to byte representation
func (sss *ShamirSecretSharing) shareToBytes(share *Share) []byte {
	xBytes := share.X.Bytes()
	yBytes := share.Y.Bytes()
	
	// Format: [x_len][x_bytes][y_len][y_bytes]
	result := make([]byte, 0, 4+len(xBytes)+len(yBytes))
	
	// X coordinate length and data
	result = append(result, byte(len(xBytes)>>8), byte(len(xBytes)))
	result = append(result, xBytes...)
	
	// Y coordinate length and data
	result = append(result, byte(len(yBytes)>>8), byte(len(yBytes)))
	result = append(result, yBytes...)
	
	return result
}

// ShareToBytes is a public wrapper for shareToBytes
func ShareToBytes(share *Share) []byte {
	sss := NewShamirSecretSharing()
	return sss.shareToBytes(share)
}

// BytesToShare is a public wrapper for bytesToShare
func BytesToShare(data []byte) (*Share, error) {
	sss := NewShamirSecretSharing()
	return sss.bytesToShare(data)
}

// bytesToShare converts byte representation back to a share
func (sss *ShamirSecretSharing) bytesToShare(data []byte) (*Share, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid share data: too short")
	}
	
	// Read X coordinate
	xLen := int(data[0])<<8 | int(data[1])
	if len(data) < 2+xLen+2 {
		return nil, fmt.Errorf("invalid share data: insufficient data for x coordinate")
	}
	
	xBytes := data[2 : 2+xLen]
	x := new(big.Int).SetBytes(xBytes)
	
	// Read Y coordinate
	yLen := int(data[2+xLen])<<8 | int(data[2+xLen+1])
	if len(data) < 2+xLen+2+yLen {
		return nil, fmt.Errorf("invalid share data: insufficient data for y coordinate")
	}
	
	yBytes := data[2+xLen+2 : 2+xLen+2+yLen]
	y := new(big.Int).SetBytes(yBytes)
	
	return &Share{X: x, Y: y}, nil
}