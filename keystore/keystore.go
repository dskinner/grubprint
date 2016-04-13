// Package keystore provides storage, retrieval, and verification
// of public keys for oauth2 bearer authorization schemes.
//
// This is to support server-to-server interactions, often referred to as
// two-legged OAuth.
//
// Implementation of this package is guided by the following RFCs:
//
// The OAuth 2.0 Authorization Framework
// https://tools.ietf.org/html/rfc6749
//
// The OAuth 2.0 Authorization Framework: Bearer Token Usage
// https://tools.ietf.org/html/rfc6750
//
// TODO require http.Request references to assure this library
// is only used with https. A global Strict bool could manage
// enforcement and when Strict is false, the library must produce
// verbose warnings. When true, panic if not https.
package keystore

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2/jws"
)

var (
	// Default provides a default in-memory keystore implementation.
	Default = New(nil)

	// ExpiresIn is the default expiry time for access tokens.
	ExpiresIn = 3600 * time.Second

	// ErrInvalidRequest represents an error for when the request parameters are malformed.
	ErrInvalidRequest = Error{Err: "invalid_request", Code: 400}

	// ErrInvalidClient represents an error for when client authentication failed,
	// e.g. unknown client, no client authentication included, or unsupported authentication method.
	ErrInvalidClient = Error{Err: "invalid_client", Code: 401}

	// ErrInvalidToken represents an error for when the access token provided is
	// expired, revoked, malformed, or invalid for other reasons.
	ErrInvalidToken = Error{Err: "invalid_token", Code: 401}

	// ErrUnauthorizedClient represents an error for when the authenticated client is not authorized
	// to use this authorization grant type.
	ErrUnauthorizedClient = Error{Err: "unauthorized_client", Code: 400}

	// ErrUnsupportedGrantType represents an error for when the authorization grant type is not
	// supported by the authorization server.
	ErrUnsupportedGrantType = Error{Err: "unsupported_grant_type", Code: 400}

	// ErrInvalidScope represents an error for when the requested scope is invalid, unknown,
	// malformed, or exceeds the scope granted by the resource owner.
	ErrInvalidScope = Error{Err: "invalid_scope", Code: 400}

	// stubbed for tests
	now = time.Now
)

// Error represents error types as defined in RFC 6749 and RFC 6750.
type Error struct {
	Err  string `json:"error"`
	Desc string `json:"error_description,omitempty"`
	Code int    `json:"-"`
}

func (e Error) Error() string { return e.Desc }

// as returns a copy of error with description.
func (e Error) as(desc string) Error {
	e.Desc = desc
	return e
}

// HandleError is a helper for writing appropriate headers and body to response.
// Writes Content-Type and Cache-Control headers, and writes serialized error as json.
// If error results in 401, WWW-Authenticate header will be set.
func HandleError(w http.ResponseWriter, e Error) {
	switch e {
	case ErrInvalidClient:
		w.Header().Set("WWW-Authenticate", "Bearer")
	case ErrInvalidToken:
		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Bearer error="%s",error_description="%s"`, e.Err, e.Desc))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "private, no-store")
	w.WriteHeader(e.Code)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		log.Printf("encode error failed: %s\n", err)
	}
}

// TokenHandler implements an oauth token response handler using Default keystore.
var TokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	x := r.FormValue("assertion")
	if err := Verify(x); err != nil {
		switch err := err.(type) {
		case Error:
			if err.Err != "invalid_token" {
				HandleError(w, err)
				return
			}
		default:
			http.Error(w, err.Error(), 500)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "private, no-store")
	type tokenRes struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"` // seconds
	}
	if err := json.NewEncoder(w).Encode(tokenRes{AccessToken: x, TokenType: "Bearer", ExpiresIn: int(ExpiresIn / time.Second)}); err != nil {
		log.Printf("encode token failed: %s\n", err)
	}
})

// Keygen generates PEM encoded PKCS1 key pair.
func Keygen() (public []byte, private []byte, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, nil, err
	}

	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(priv),
	}
	private = pem.EncodeToMemory(&privBlock)

	pubDer, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	pubBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubDer,
	}
	public = pem.EncodeToMemory(&pubBlock)

	return
}

// Store is an interface for retrieval and storage of public keys.
type Store interface {
	// Get returns the PEM encoded public key or error.
	Get(string) ([]byte, error)

	// Set stores PEM encoded public key and must assure key is valid.
	Set(string, []byte) error
}

type memstore struct {
	sync.RWMutex
	keymap map[string][]byte
}

func (st *memstore) Get(id string) ([]byte, error) {
	st.RLock()
	defer st.RUnlock()

	key, ok := st.keymap[id]
	if !ok {
		return nil, fmt.Errorf("id %q does not exist", id)
	}
	return key, nil
}

func (st *memstore) Set(id string, key []byte) error {
	st.Lock()
	defer st.Unlock()

	if _, err := parsePublicKey(key); err != nil {
		return err
	}
	st.keymap[id] = key
	return nil
}

// Get is a helper that operates on Default.
func Get(id string) ([]byte, error) { return Default.Get(id) }

// Set is a helper that operates on Default.
func Set(id string, key []byte) error { return Default.Set(id, key) }

// Verify is a helper that operates on Default.
func Verify(token string) error { return Default.Verify(token) }

// VerifyRequest is a helper that operates on Default.
func VerifyRequest(r *http.Request) error { return Default.VerifyRequest(r) }

// Keystore provides methods for token verification.
type Keystore struct {
	Store
}

// New returns a new Keystore. If st is nil, an in-memory Store implementation is used.
func New(st Store) *Keystore {
	if st == nil {
		st = &memstore{keymap: make(map[string][]byte)}
	}
	return &Keystore{st}
}

// Verify decodes token and verifies signature against store.
func (ks *Keystore) Verify(token string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ErrInvalidRequest.as("invalid token")
	}

	// verify claim set
	payload, err := base64Decode(parts[1])
	if err != nil {
		return ErrInvalidRequest.as("base64 decode payload failed")
	}
	var cs jws.ClaimSet
	if err := json.NewDecoder(bytes.NewBuffer(payload)).Decode(&cs); err != nil {
		return ErrInvalidRequest.as("json decode payload failed")
	}

	// verify signature
	bin, err := ks.Get(cs.Iss)
	if err != nil {
		return ErrInvalidClient.as("unknown client id")
	}
	publicKey, err := parsePublicKey(bin)
	if err != nil {
		return err
	}
	sig, err := base64Decode(parts[2])
	if err != nil {
		return ErrInvalidRequest.as("base64 decode signature failed")
	}
	if !crypto.SHA256.Available() {
		return fmt.Errorf("Hash unavailable")
	}
	hash := crypto.SHA256.New()
	hash.Write([]byte(strings.Join(parts[:2], ".")))

	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash.Sum(nil), sig); err != nil {
		return ErrInvalidClient.as("key verification failed")
	}

	// TODO these checks are intentionally last. Need to consider TokenHandler which needs to verify
	// signature but ignore expiry when issuing new tokens.
	if time.Unix(cs.Iat, 0).After(now()) {
		return ErrInvalidToken.as("invalid timestamp")
	}
	if time.Unix(cs.Exp, 0).Before(now()) {
		return ErrInvalidToken.as("token expired")
	}

	return nil
}

// VerifyRequest inspects the Authorization request header and verifies
// the header token signature against store.
func (ks *Keystore) VerifyRequest(r *http.Request) error {
	x := r.Header.Get("Authorization")
	if x == "" {
		return ErrInvalidClient.as("client did not provide authentication")
	}
	parts := strings.Split(x, " ")
	if len(parts) != 2 {
		return ErrInvalidRequest.as("malformed authentication")
	}
	if parts[0] != "Bearer" {
		return ErrInvalidClient.as("only bearer scheme supported")
	}
	return ks.Verify(parts[1])
}

// parsePublicKey parses PEM encoded PKCS1 public key.
func parsePublicKey(bin []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(bin)
	if block == nil {
		return nil, fmt.Errorf("key must be PEM encoded")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	if pub, ok := key.(*rsa.PublicKey); ok {
		return pub, nil
	}

	return nil, fmt.Errorf("key not correct type")
}

// base64Decode decodes the Base64url encoded string.
func base64Decode(s string) ([]byte, error) {
	// add back missing padding
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}
