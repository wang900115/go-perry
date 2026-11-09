// Error wrapping in Go: %w, errors.Is / errors.As / errors.AsType (Go 1.26),
// custom Is/Unwrap methods, multi-error trees (errors.Join, multiple %w),
// and the typed-nil pitfall.
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// ---------- sentinel errors ----------

var (
	ErrNotFound   = errors.New("not found")
	ErrPermission = errors.New("permission denied")
)

// ---------- a custom error type ----------

// QueryError carries structured context and wraps a cause.
type QueryError struct {
	Table string
	ID    int
	Err   error
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("query %s id=%d: %v", e.Table, e.ID, e.Err)
}

// Unwrap exposes the cause so errors.Is / errors.As can look through it.
func (e *QueryError) Unwrap() error { return e.Err }

// ---------- a type with a custom Is method ----------

// HTTPError matches ErrNotFound when its status code is 404,
// even though it does not wrap ErrNotFound.
type HTTPError struct{ Code int }

func (e HTTPError) Error() string { return fmt.Sprintf("http status %d", e.Code) }

func (e HTTPError) Is(target error) bool {
	return target == ErrNotFound && e.Code == 404
}

// ---------- layered "application" code ----------

func repoFind(id int) error {
	if id == 42 {
		return &QueryError{Table: "users", ID: id, Err: ErrNotFound}
	}
	return nil
}

func serviceGetUser(id int) error {
	if err := repoFind(id); err != nil {
		return fmt.Errorf("service: get user: %w", err) // %w keeps the chain
	}
	return nil
}

func serviceGetUserLossy(id int) error {
	if err := repoFind(id); err != nil {
		return fmt.Errorf("service: get user: %v", err) // %v flattens to text
	}
	return nil
}

// ---------- tree printer ----------

// printTree walks the error tree the same way errors.Is does:
// the node itself, then depth-first through Unwrap() error / Unwrap() []error.
func printTree(err error, depth int) {
	fmt.Printf("%s- %T: %q\n", strings.Repeat("  ", depth), err, err.Error())
	switch u := err.(type) {
	case interface{ Unwrap() []error }:
		for _, child := range u.Unwrap() {
			printTree(child, depth+1)
		}
	case interface{ Unwrap() error }:
		if child := u.Unwrap(); child != nil {
			printTree(child, depth+1)
		}
	}
}

// ---------- typed-nil pitfall ----------

func validateBad(ok bool) error {
	var qe *QueryError // nil pointer of a concrete type
	if !ok {
		qe = &QueryError{Table: "orders", ID: 1, Err: ErrPermission}
	}
	return qe // BUG: a nil *QueryError inside a non-nil error interface
}

func validateGood(ok bool) error {
	if !ok {
		return &QueryError{Table: "orders", ID: 1, Err: ErrPermission}
	}
	return nil // return an untyped nil
}

func main() {
	fmt.Printf("== 1. %%w keeps the chain, %%v breaks it ==\n")
	err := serviceGetUser(42)
	fmt.Println("err:", err)
	fmt.Println("err == ErrNotFound         :", err == ErrNotFound)
	fmt.Println("errors.Is(err, ErrNotFound):", errors.Is(err, ErrNotFound))
	lossy := serviceGetUserLossy(42)
	fmt.Println("lossy errors.Is            :", errors.Is(lossy, ErrNotFound))
	fmt.Println("errors.Unwrap(err)         :", errors.Unwrap(err))

	fmt.Println("\n== 2. errors.As vs errors.AsType (Go 1.26) ==")
	var qe *QueryError
	if errors.As(err, &qe) {
		fmt.Printf("As     -> table=%s id=%d\n", qe.Table, qe.ID)
	}
	if qe2, ok := errors.AsType[*QueryError](err); ok {
		fmt.Printf("AsType -> table=%s id=%d\n", qe2.Table, qe2.ID)
	}
	_, ok := errors.AsType[*fs.PathError](err)
	fmt.Println("AsType[*fs.PathError] found:", ok)

	fmt.Println("\n== 3. custom Is method ==")
	h404 := fmt.Errorf("fetch profile: %w", HTTPError{Code: 404})
	h500 := fmt.Errorf("fetch profile: %w", HTTPError{Code: 500})
	fmt.Println("404 is ErrNotFound:", errors.Is(h404, ErrNotFound))
	fmt.Println("500 is ErrNotFound:", errors.Is(h500, ErrNotFound))

	fmt.Println("\n== 4. multi-error trees ==")
	joined := errors.Join(
		serviceGetUser(42),
		nil, // nil entries are dropped
		fmt.Errorf("cache: %w", ErrPermission),
	)
	fmt.Printf("Join message (newline-separated):\n%v\n", joined)
	fmt.Println("Is ErrNotFound  :", errors.Is(joined, ErrNotFound))
	fmt.Println("Is ErrPermission:", errors.Is(joined, ErrPermission))
	fmt.Println("errors.Unwrap(joined) == nil:", errors.Unwrap(joined) == nil)
	fmt.Println("errors.Join(nil, nil) == nil:", errors.Join(nil, nil) == nil)

	_, statErr := os.Stat("definitely-missing-file.txt")
	multi := fmt.Errorf("load config: %w; fallback: %w", statErr, h404)
	fmt.Printf("multi-wrap message: %v\n", multi)
	fmt.Println("Is fs.ErrNotExist:", errors.Is(multi, fs.ErrNotExist))
	fmt.Println("Is ErrNotFound   :", errors.Is(multi, ErrNotFound))
	fmt.Println("tree:")
	printTree(multi, 1)

	fmt.Println("\n== 5. typed-nil pitfall ==")
	bad := validateBad(true)
	good := validateGood(true)
	fmt.Printf("validateBad(true)  != nil: %v (dynamic type %T)\n", bad != nil, bad)
	fmt.Printf("validateGood(true) != nil: %v\n", good != nil)
}
