package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

func warn(msg string, args ...interface{}) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, msg, args...)
}

func verbose(msg string, args ...interface{}) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, msg, args...)
}

const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZöäüÖÄÜß"

func check(word string) bool {
	if !utf8.ValidString(word) {
		// verbose("word %q is not valid UTF-8", word)
		return false
	}

	// note: range over string iterates over runes (decoding UTF-8 characters)
	runes := 0
	for _, r := range word {
		if r == utf8.RuneError {
			// verbose("word %q contains the UTF-8 replacement character", word)
			return false
		}

		if !strings.ContainsAny(allowedChars, string(r)) {
			// verbose("word %q contains invalid character %U (%c)", word, r, r)
			return false
		}

		runes++
	}

	if runes < 3 {
		// verbose("word %q is too short", word)
		return false
	}

	if runes > 11 {
		// verbose("word %q is too long", word)
		return false
	}

	// filter special characters
	if strings.ContainsAny(word, `"'-./ 0123456789`) {
		// verbose("word %q contains special character", word)
		return false
	}

	return true
}

func tolower(s string) (result string) {
	for _, r := range s {
		result += string(unicode.ToLower(r))
	}
	return result
}

const wordcount = 7776

func main() {
	if len(os.Args) != 2 {
		warn("usage: %v FILE", os.Args[0])
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	words := make(map[string]struct{})

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if sc.Err() != nil {
			if err != nil {
				panic(err)
			}
		}

		fields := strings.Fields(sc.Text())
		word := tolower(fields[len(fields)-1])
		if !check(word) {
			continue
		}

		if _, ok := words[word]; ok {
			// verbose("word %q already found", word)
			continue
		}

		words[word] = struct{}{}
		if len(words) == wordcount {
			fmt.Printf("%d words found\n", len(words))
			break
		}
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	sorted := make([]string, 0, len(words))

	for word := range words {
		sorted = append(sorted, word)
	}

	sort.Strings(sorted)
	for _, word := range sorted {
		fmt.Printf("%v\n", word)
	}
}
