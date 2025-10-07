package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/ed25519"
)

type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
}

// 發行 v2.local
func issueV2Local(symKey []byte, c Claims, ttl time.Duration) (string, error) {
	v2 := paseto.NewV2()
	now := time.Now().UTC()

	tok := paseto.JSONToken{
		Issuer:     "perry",
		Subject:    c.UserID,
		Audience:   "example-app",
		Jti:        randomJTI(),
		IssuedAt:   now,
		NotBefore:  now.Add(-5 * time.Second),
		Expiration: now.Add(ttl),
	}
	tok.Set("uid", c.UserID)
	tok.Set("role", c.Role)
	footer := "k=v2.local.demo"
	return v2.Encrypt(symKey, tok, footer)
}

// 驗證 v2.local
func verifyV2Local(symKey []byte, token string) (*paseto.JSONToken, string, error) {
	v2 := paseto.NewV2()
	var out paseto.JSONToken
	var footer string
	if err := v2.Decrypt(token, symKey, &out, &footer); err != nil {
		return nil, "", err
	}
	// 時效驗證
	if err := out.Validate(paseto.ValidAt(time.Now().UTC())); err != nil {
		return nil, "", err
	}
	return &out, footer, nil
}

// 發行 v2.Public
func issueV2Public(priv ed25519.PrivateKey, c Claims, ttl time.Duration) (string, error) {
	v2 := paseto.NewV2()
	now := time.Now().UTC()

	tok := paseto.JSONToken{
		Issuer:     "perry",
		Subject:    c.UserID,
		Audience:   "example-app",
		Jti:        randomJTI(),
		IssuedAt:   now,
		NotBefore:  now.Add(-5 * time.Second),
		Expiration: now.Add(ttl),
	}
	tok.Set("uid", c.UserID)
	tok.Set("role", c.Role)
	footer := "k=v2.public.demo"
	return v2.Sign(priv, tok, footer)
}

// 驗證 v2.public
func verifyV2Public(pub ed25519.PublicKey, token string) (*paseto.JSONToken, string, error) {
	v2 := paseto.NewV2()
	var out paseto.JSONToken
	var footer string
	if err := v2.Verify(token, pub, &out, &footer); err != nil {
		return nil, "", err
	}
	if err := out.Validate(paseto.ValidAt(time.Now().UTC())); err != nil {
		return nil, "", err
	}
	return &out, footer, nil
}

func randomSymKey() []byte {
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	return key
}

func randomJTI() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func main() {
	// local
	symKey := randomSymKey()
	localTok, err := issueV2Local(symKey, Claims{UserID: "1", Role: "admin"}, 30*time.Minute)
	must(err)
	fmt.Println("v2.local token:", localTok)

	localOut, localFooter, err := verifyV2Local(symKey, localTok)
	must(err)
	fmt.Printf("llocal Footer: %s\n", localFooter)
	b, _ := json.MarshalIndent(localOut, "", "  ")
	fmt.Println("local Payload JSON:", string(b))
	// public

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	must(err)

	publicTok, err := issueV2Public(priv, Claims{UserID: "2", Role: "manager"}, 30*time.Minute)
	must(err)

	publicOut, publicFooter, err := verifyV2Public(pub, publicTok)
	must(err)
	fmt.Printf("public Footer: %s\n", publicFooter)

	b, _ = json.MarshalIndent(publicOut, "", "  ")
	fmt.Println("public Payload JSON:", string(b))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
