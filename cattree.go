package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"unicode/utf8"

	ignore "github.com/sabhiram/go-gitignore"
)

const sniffLen = 8000 // Number of bytes to read for detection

// isTextFile determines if a file is text by reading its first few KB.
func isTextFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, sniffLen)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}
	return isText(buf[:n])
}

// isText returns true if the data is likely text.
func isText(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	// If NUL byte is present, it's binary
	if bytes.IndexByte(data, 0) != -1 {
		return false
	}
	// Count non-printable bytes
	nonPrintable := 0
	total := 0
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		if r == utf8.RuneError && size == 1 {
			nonPrintable++
		} else if r < 32 && r != 9 && r != 10 && r != 13 {
			nonPrintable++
		}
		total++
		data = data[size:]
	}
	// If more than 10% are non-printable, treat as binary
	return float64(nonPrintable)/float64(total) < 0.1
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [directory]\n\nPrints all text files in a directory tree, respecting .gitignore.\n",
			os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	// Load .gitignore if present in the root directory
	var ign *ignore.GitIgnore
	gitignorePath := filepath.Join(root, ".gitignore")
	if f, err := os.Open(gitignorePath); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		ign = ignore.CompileIgnoreLines(lines...)
	}

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", path, err)
			return nil
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		relPath, _ := filepath.Rel(root, path)
		if ign != nil && ign.MatchesPath(relPath) {
			return nil
		}
		// Only print text files
		if !isTextFile(path) {
			return nil
		}
		fmt.Printf("==%s==\n", relPath)
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("  [error opening file: %v]\n", err)
			return nil
		}
		defer f.Close()
		io.Copy(os.Stdout, f)
		fmt.Println()
		return nil
	})
}
