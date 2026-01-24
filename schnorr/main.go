package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// This is a simple implementation of Schnorr signatures in Go.
// It demonstrates key generation, signing, and verification.

// x = private key
// g = generator
// p = prime modulus
// y = g^x mod p (public key)

var (
	p = big.NewInt(23)
	g = big.NewInt(5)
)

func randMod(n *big.Int) *big.Int {
	r, _ := rand.Int(rand.Reader, n)
	return r
}

type Prover struct {
	x *big.Int // private key
	y *big.Int // public key
}

func NewProver() *Prover {
	x := randMod(new(big.Int).Sub(p, big.NewInt(2)))
	x.Add(x, big.NewInt(1))
	y := new(big.Int).Exp(g, x, p)
	return &Prover{x: x, y: y}
}

// step 1: Prover generates a commitment
func (pr *Prover) Commit() (*big.Int, *big.Int) {
	r := randMod(new(big.Int).Sub(p, big.NewInt(1)))
	t := new(big.Int).Exp(g, r, p)
	return r, t
}

// step 3: Prover generates a response
func (pr *Prover) Respond(r, c *big.Int) *big.Int {
	mod := new(big.Int).Sub(p, big.NewInt(1))
	s := new(big.Int).Mul(c, pr.x)
	s.Add(s, r)
	s.Mod(s, mod)
	return s
}

func Verify(t, s, c, y *big.Int) bool {
	// left = g^s mod p
	left := new(big.Int).Exp(g, s, p)

	// right = t * y^c mod p
	yc := new(big.Int).Exp(y, c, p)
	right := new(big.Int).Mul(t, yc)
	right.Mod(right, p)

	return left.Cmp(right) == 0
}

func main() {
	prover := NewProver()

	fmt.Println("Public key y =", prover.y)

	// Step 1: Commit
	r, t := prover.Commit()
	fmt.Println("Commit t =", t)

	// Step 2: Challenge (Verifier)
	c := randMod(new(big.Int).Sub(p, big.NewInt(1)))
	fmt.Println("Challenge c =", c)

	// Step 3: Response
	s := prover.Respond(r, c)
	fmt.Println("Response s =", s)

	// Step 4: Verify
	ok := Verify(t, s, c, prover.y)
	fmt.Println("Verify result =", ok)
}
