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
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2/jws"
)

var (
	Default = New()

	// TODO better errors throughout
	ErrInvalidToken     = &errorCode{"Invalid token", 400}
	ErrInvalidTimestamp = &errorCode{"Invalid timestamp", 400}
	ErrInvalidIssuer    = &errorCode{"Invalid issuer", 401}
	ErrExpiredToken     = &errorCode{"Expired token", 401}
)

type Error interface {
	error
	Code() int
}

type errorCode struct {
	err  string
	code int
}

func (e *errorCode) Error() string { return e.err }
func (e *errorCode) Code() int     { return e.code }

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

type Keystore struct {
	sync.RWMutex
	keymap map[string][]byte
}

func New() *Keystore {
	return &Keystore{keymap: make(map[string][]byte)}
}

func (ks *Keystore) Get(id string) ([]byte, error) {
	ks.RLock()
	defer ks.RUnlock()

	key, ok := ks.keymap[id]
	if !ok {
		return nil, ErrInvalidIssuer
	}
	return key, nil
}

func Get(id string) ([]byte, error) { return Default.Get(id) }

func (ks *Keystore) Set(id string, key []byte) error {
	ks.Lock()
	defer ks.Unlock()

	if _, err := parsePublicKey(key); err != nil {
		return err
	}
	ks.keymap[id] = key
	return nil
}

func Set(id string, key []byte) error { return Default.Set(id, key) }

func (ks *Keystore) Verify(token string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ErrInvalidToken
	}

	// verify claim set
	payloadJson, err := base64Decode(parts[1])
	if err != nil {
		return err
	}
	var cs jws.ClaimSet
	if err := json.NewDecoder(bytes.NewBuffer(payloadJson)).Decode(&cs); err != nil {
		return err
	}
	if time.Unix(cs.Iat, 0).After(time.Now()) {
		return ErrInvalidTimestamp
	}
	if time.Unix(cs.Exp, 0).Before(time.Now()) {
		return ErrExpiredToken
	}

	// verify signature
	bin, err := ks.Get(cs.Iss)
	if err != nil {
		return err
	}
	publicKey, err := parsePublicKey(bin)
	if err != nil {
		return err
	}
	sig, err := base64Decode(parts[2])
	if err != nil {
		return err
	}
	if !crypto.SHA256.Available() {
		return fmt.Errorf("Hash unavailable")
	}
	hash := crypto.SHA256.New()
	hash.Write([]byte(strings.Join(parts[:2], ".")))

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash.Sum(nil), sig)
}

func Verify(token string) error { return Default.Verify(token) }

func (ks *Keystore) VerifyRequest(r *http.Request) error {
	x := r.Header.Get("Authorization")
	parts := strings.Split(x, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fmt.Errorf("Not verified")
	}
	return ks.Verify(parts[1])
}

func VerifyRequest(r *http.Request) error { return Default.VerifyRequest(r) }

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
