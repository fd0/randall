package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/c-bata/go-prompt"
	"github.com/fd0/randall/wordlists"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/spf13/pflag"
)

const (
	JoinModeSpace = "space"
	JoinModeCamel = "camel"
	JoinModeMinus = "minus"
)

// Options collects options set via command-line flags.
type Options struct {
	Words       uint
	Wordlist    string
	Passphrases uint
	Reconstruct bool
	JoinMode    string
}

func die(msg string, args ...interface{}) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func capitalizeWord(w string) string {
	r, size := utf8.DecodeRuneInString(w)
	return string(strings.ToUpper(string(r))) + w[size:]
}

func capitalize(words []string) {
	for i, word := range words {
		words[i] = capitalizeWord(word)
	}
}

func generate(n uint, list []string, joinMode string) string {
	words := make([]string, 0, n)
	for i := uint(0); i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(list))))
		if err != nil {
			die("fatal: unable to select random word: %v", err)
		}

		if !num.IsInt64() {
			die("fatal: number is too large")
		}

		idx := num.Int64()
		if idx > int64(len(list)) {
			die("fatal: selected number is larger than wordlist")
		}

		words = append(words, list[idx])
	}

	switch strings.ToLower(joinMode) {
	case JoinModeSpace, "":
		return strings.Join(words, " ")
	case JoinModeMinus:
		return strings.Join(words, "-")
	case JoinModeCamel:
		capitalize(words)
		return strings.Join(words, "")
	default:
		die("unknown join mode %q: choose between %q, %q and %q", joinMode,
			JoinModeSpace, JoinModeCamel, JoinModeMinus)
		return ""
	}
}

func completer(wordlist []string) prompt.Completer {
	capitalize(wordlist)
	return func(d prompt.Document) []prompt.Suggest {
		matches := fuzzy.RankFindFold(d.GetWordBeforeCursor(), wordlist)
		sort.Sort(matches)

		var suggestions []prompt.Suggest
		for _, match := range matches {
			suggestions = append(suggestions, prompt.Suggest{
				Text: match.Target,
			})
		}

		return suggestions
	}
}

func main() {
	var opts Options

	flags := pflag.NewFlagSet("randall", pflag.ContinueOnError)
	flags.UintVarP(&opts.Words, "words", "w", 5, "generate passphrase with `n` words")
	flags.UintVarP(&opts.Passphrases, "passphrases", "n", 1, "generate `n` passphrases")
	flags.StringVarP(&opts.Wordlist, "wordlist", "l", "en", fmt.Sprintf("use `wordlist` as the source for words (valid: %v)", strings.Join(wordlists.Names(), ", ")))
	flags.BoolVarP(&opts.Reconstruct, "reconstruct", "r", false, "interactively reconstruct a password based on a wordlist")

	mode := JoinModeSpace
	if v, ok := os.LookupEnv("RANDALL_JOIN_MODE"); ok {
		mode = v
	}

	flags.StringVarP(&opts.JoinMode, "join-mode", "j", mode, "choose how words are joined (\"space\" (like this), \"camel\" (LikeThis) or \"minus\" (like-this))")

	err := flags.Parse(os.Args)
	if err == pflag.ErrHelp {
		os.Exit(0)
	}
	if err != nil {
		die("error: %v", err)
	}

	list, err := wordlists.Get(opts.Wordlist)
	if err != nil {
		die("error: %v", err)
	}

	if opts.Reconstruct {
		fmt.Printf("reconstruct password consisting of %d words using wordlist %s\n", opts.Words, opts.Wordlist)
		fmt.Printf("\ntype first word, complete with <tab>, press <enter> to add word\n")

		validWords := make(map[string]struct{})
		for _, word := range list {
			validWords[capitalizeWord(word)] = struct{}{}
		}

		prefix := "> "
		var words []string

		p := prompt.New(
			func(word string) {
				if word == "" {
					fmt.Printf("aborting\n")
					os.Exit(0)
				}

				if _, ok := validWords[word]; !ok {
					fmt.Printf("word %q is not in wordlist, try again\n", word)
					return
				}

				words = append(words, word)
				fmt.Printf("add word %q to password\n", word)

				if uint(len(words)) == opts.Words {
					fmt.Printf("password is: %s\n", strings.Join(words, ""))
					os.Exit(0)
				}
			},
			completer(list),
			prompt.OptionPrefix(prefix),
			prompt.OptionLivePrefix(func() (string, bool) {
				if len(words) == 0 {
					return "", false
				}

				return strings.Join(words, "") + " > ", true
			}),
			prompt.OptionTitle("foo"),
		)

		p.Run()
		return
	}

	for i := uint(0); i < opts.Passphrases; i++ {
		fmt.Printf("%v\n", generate(opts.Words, list, opts.JoinMode))
	}
}
