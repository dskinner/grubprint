// Package keystore provides rigid storage, retrieval, and verification
// of public keys for oauth2 bearer authorization schemes for servers.
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

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jws"
)

var (
	// Default provides a default in-memory keystore implementation.
	Default = New(nil)

	// The request parameters are malformed.
	ErrInvalidRequest = &errorCode{err: "invalid_request", code: 400}

	// Client authentication failed
	// (e.g., unknown client, no client authentication included, or unsupported authentication method)
	ErrInvalidClient = &errorCode{err: "invalid_client", code: 401}

	// The provided authorization grant is invalid, expired, revoked.
	ErrInvalidGrant = &errorCode{err: "invalid_grant", code: 400}

	// The authenticated client is not authorized to use this authorization grant type.
	ErrUnauthorizedClient = &errorCode{err: "unauthorized_client", code: 400}

	// The authorization grant type is not supported by the authorization server.
	ErrUnsupportedGrantType = &errorCode{err: "unsupported_grant_type", code: 400}

	// The requested scope is invalid, unknown, malformed, or exceeds the scope granted by the resource owner.
	ErrInvalidScope = &errorCode{err: "invalid_scope", code: 400}
)

type Error interface {
	error
	Code() int
}

type errorCode struct {
	err  string `json:"error"`
	desc string `json:"error_description"`
	code int    `json:omit`
}

func (e *errorCode) Error() string { return e.desc }
func (e *errorCode) Code() int     { return e.code }

// as returns a copy of error with description.
func (e errorCode) as(desc string) *errorCode {
	e.desc = desc
	return &e
}

// HandleError is a helper for writing appropriate headers and body to response.
// Writes Content-Type and Cache-Control headers, and writes serialized error as json.
// If error results in 401, WWW-Authenticate header will be set.
func HandleError(w http.ResponseWriter, e Error) {
	if e.Code() == 401 {
		w.Header().Set("WWW-Authenticate", "Bearer")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "private, no-store")
	w.WriteHeader(e.Code())
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
			HandleError(w, err)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "private, no-store")
	if err := json.NewEncoder(w).Encode(oauth2.Token{AccessToken: x}); err != nil {
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

// Get is a helper that operates on Default.
func Get(id string) ([]byte, error) { return Default.Get(id) }

// Set is a helper that operates on Default.
func Set(id string, key []byte) error { return Default.Set(id, key) }

// Verify is a helper that operates on Default.
func Verify(token string) error { return Default.Verify(token) }

// VerifyRequest is a helper that operates on Default.
func VerifyRequest(r *http.Request) error { return Default.VerifyRequest(r) }

func (st *memstore) Set(id string, key []byte) error {
	st.Lock()
	defer st.Unlock()

	if _, err := parsePublicKey(key); err != nil {
		return err
	}
	st.keymap[id] = key
	return nil
}

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
	payloadJson, err := base64Decode(parts[1])
	if err != nil {
		return ErrInvalidRequest.as("base64 decode payload failed")
	}
	var cs jws.ClaimSet
	if err := json.NewDecoder(bytes.NewBuffer(payloadJson)).Decode(&cs); err != nil {
		return ErrInvalidRequest.as("json decode payload failed")
	}
	if time.Unix(cs.Iat, 0).After(time.Now()) {
		return ErrInvalidGrant.as("invalid timestamp")
	}
	if time.Unix(cs.Exp, 0).Before(time.Now()) {
		return ErrInvalidGrant.as("token expired")
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
