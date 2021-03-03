package jwks

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
)

const (

	// RS256 represents a cryptography key generated by an RSA algorithm.
	RS256 = "RS256"
)

var (

	// ErrKIDNotFound indicates that the given key ID was not found in the JWKS.
	ErrKIDNotFound = errors.New("the given key ID was not found in the JWKS")

	// ErrNotExpectedKeyType indicates that the given public key was not of the expected type.
	ErrNotExpectedKeyType = errors.New("the public key was not of the expected type")
)

// JSONKey represents a raw key inside a JWKS.
type JSONKey struct {
	Exponent    string `json:"e"`
	ID          string `json:"kid"`
	Alg         string `json:"alg"`
	Modulus     string `json:"n"`
	precomputed interface{}
}

// Keystore represents a JWKS.
type Keystore map[string]*JSONKey

// rawKeystore represents a JWKS in JSON format.
type rawKeystore struct {
	Keys []JSONKey `json:"keys"`
}

// New creates a new JWKS from a raw JSON message.
func New(keystoreBytes json.RawMessage) (keystore Keystore, err error) {

	// Turn the raw JWKS into the correct Go type.
	var rawKS rawKeystore
	if err = json.Unmarshal(keystoreBytes, &rawKS); err != nil {
		return nil, err
	}

	// Iterate through the keys in the raw keystore. Add them to the JWKS.
	keystore = make(map[string]*JSONKey, len(rawKS.Keys))
	for _, key := range rawKS.Keys {
		keystore[key.ID] = &key
	}

	return keystore, nil
}

// RSA retrieves an RSA public key from the JWKS.
func (k Keystore) RSA(kid string) (publicKey *rsa.PublicKey, err error) {

	// Get the JSONKey from the JWKS.
	key, ok := k[kid]
	if !ok {
		return nil, ErrKIDNotFound
	}

	// Confirm the key is an RSA key.
	if key.Alg != RS256 {
		return nil, ErrNotExpectedKeyType
	}

	// Transform the key from JSON to an RSA key.
	return key.RSA()
}