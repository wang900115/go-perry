// os_root demonstrates os.Root (Go 1.24+, extra methods in Go 1.25):
// traversal-resistant file access confined to one directory tree.
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	// Layout (inside a temp dir):
	//   secret.txt            <- must never be reachable from "public"
	//   public/hello.txt
	//   public/docs/a.txt
	//   public/link_in  -> docs/a.txt     (symlink that stays inside)
	//   public/link_out -> ../secret.txt  (symlink that escapes)
	base, err := os.MkdirTemp("", "osroot-*")
	must(err)
	defer os.RemoveAll(base)

	public := filepath.Join(base, "public")
	must(os.MkdirAll(filepath.Join(public, "docs"), 0o755))
	must(os.WriteFile(filepath.Join(base, "secret.txt"), []byte("TOP SECRET"), 0o644))
	must(os.WriteFile(filepath.Join(public, "hello.txt"), []byte("hello"), 0o644))
	must(os.WriteFile(filepath.Join(public, "docs", "a.txt"), []byte("doc A"), 0o644))

	// Symlinks may need Developer Mode / admin rights on Windows.
	symlinks := os.Symlink(filepath.Join("docs", "a.txt"), filepath.Join(public, "link_in")) == nil &&
		os.Symlink(filepath.Join("..", "secret.txt"), filepath.Join(public, "link_out")) == nil
	if !symlinks {
		fmt.Println("(symlinks not available on this system; symlink cases skipped)")
	}

	inputs := []string{"hello.txt", "docs/a.txt", "../secret.txt"}
	if symlinks {
		inputs = append(inputs, "link_in", "link_out")
	}

	fmt.Println("== 1. naive filepath.Join ==")
	for _, name := range inputs {
		report(name, func() ([]byte, error) {
			return os.ReadFile(filepath.Join(public, name))
		})
	}

	fmt.Println("== 2. filepath.IsLocal check (string-only) ==")
	for _, name := range inputs {
		report(name, func() ([]byte, error) {
			if !filepath.IsLocal(name) {
				return nil, errors.New("rejected: not a local path")
			}
			return os.ReadFile(filepath.Join(public, name))
		})
	}

	fmt.Println("== 3. os.Root ==")
	root, err := os.OpenRoot(public)
	must(err)
	defer root.Close()
	for _, name := range inputs {
		report(name, func() ([]byte, error) {
			return root.ReadFile(name) // Root.ReadFile: Go 1.25+
		})
	}

	fmt.Println("== 4. writes are confined too ==")
	err = root.WriteFile("../pwned.txt", []byte("x"), 0o644) // Go 1.25+
	fmt.Println("WriteFile ../pwned.txt ->", errLine(err))
	must(root.MkdirAll("uploads/2026", 0o755)) // Go 1.25+
	must(root.WriteFile("uploads/2026/ok.txt", []byte("ok"), 0o644))
	fmt.Println("WriteFile uploads/2026/ok.txt -> ok")

	fmt.Println("== 5. Root.FS() as an fs.FS ==")
	fsys := root.FS()
	must(fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&fs.ModeSymlink != 0 {
			fmt.Println("  ", p, "(symlink)")
			return nil
		}
		if !d.IsDir() {
			fmt.Println("  ", p)
		}
		return nil
	}))
	_, err = fs.ReadFile(fsys, "../secret.txt")
	fmt.Println("fs.ReadFile ../secret.txt ->", errLine(err))
}

func report(name string, read func() ([]byte, error)) {
	data, err := read()
	if err != nil {
		fmt.Printf("  %-14s -> %s\n", name, errLine(err))
		return
	}
	fmt.Printf("  %-14s -> %q\n", name, data)
}

// errLine strips the operation and path from an error so the output is
// portable (full messages contain temp paths and OS-specific op names).
func errLine(err error) string {
	switch {
	case err == nil:
		return "ok"
	case errors.Is(err, fs.ErrInvalid):
		return "ERROR (invalid path)"
	default:
		var pe *fs.PathError
		if errors.As(err, &pe) {
			return "ERROR (" + pe.Err.Error() + ")"
		}
		return "ERROR (" + err.Error() + ")"
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
