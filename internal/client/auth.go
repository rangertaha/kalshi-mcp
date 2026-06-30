// SPDX-License-Identifier: MIT

package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// KalshiAuthorizer signs requests with the Kalshi API key scheme: an RSA-PSS
// (SHA-256) signature over "<timestamp><METHOD><path>", sent alongside the key
// ID and timestamp in KALSHI-ACCESS-* headers.
type KalshiAuthorizer struct {
	keyID string
	key   *rsa.PrivateKey
}

// NewKalshiAuthorizer parses a PEM-encoded RSA private key (PKCS#1 or PKCS#8)
// and builds an authorizer for the given access key ID.
func NewKalshiAuthorizer(keyID, pemKey string) (*KalshiAuthorizer, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, fmt.Errorf("private key is not valid PEM")
	}
	var key *rsa.PrivateKey
	if k, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		key = k
	} else {
		parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing RSA private key: %w", err)
		}
		rk, ok := parsed.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not an RSA key")
		}
		key = rk
	}
	return &KalshiAuthorizer{keyID: keyID, key: key}, nil
}

// Authorize signs the request and sets the KALSHI-ACCESS-* headers.
func (a *KalshiAuthorizer) Authorize(r *http.Request) {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	msg := ts + r.Method + r.URL.Path
	hashed := sha256.Sum256([]byte(msg))
	sig, err := rsa.SignPSS(rand.Reader, a.key, crypto.SHA256, hashed[:], &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA256,
	})
	if err != nil {
		return // leave headers unset; the request will fail auth and surface a clear API error
	}
	r.Header.Set("KALSHI-ACCESS-KEY", a.keyID)
	r.Header.Set("KALSHI-ACCESS-SIGNATURE", base64.StdEncoding.EncodeToString(sig))
	r.Header.Set("KALSHI-ACCESS-TIMESTAMP", ts)
}

// BearerAuthorizer authenticates using an OAuth-style bearer token.
type BearerAuthorizer struct {
	header string
}

// NewBearerAuthorizer builds a BearerAuthorizer for the given token.
func NewBearerAuthorizer(token string) *BearerAuthorizer {
	return &BearerAuthorizer{header: "Bearer " + token}
}

// Authorize sets the Authorization header for bearer auth.
func (a *BearerAuthorizer) Authorize(r *http.Request) {
	r.Header.Set("Authorization", a.header)
}
