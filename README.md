[![Status badge for tests](https://github.com/fd0/randall/workflows/Build%20and%20tests/badge.svg)](https://github.com/fd0/randall/actions?query=workflow%3A%22Build+and+tests%22)

# randall

The program `randall` implements a password generator based on a sequence of words, similar to [Diceware](https://en.wikipedia.org/wiki/Diceware). It includes the EFF's [large diceware worldist](https://www.eff.org/document/passphrase-wordlists) as well as a German wordlist.

## Usage

Get help:

    $ randall --help
    Usage of randall:
      -n, --passphrases n       generate n passphrases (default 1)
      -r, --reconstruct         interactively reconstruct a password based on a wordlist
      -l, --wordlist wordlist   use wordlist as the source for words (valid: en, de) (default "en")
      -w, --words n             generate passphrase with n words (default 5)

Generate passphrase with four words in German:

    $ randall --wordlist de --words 4
    SchlangeWaldQuasiToll

Reconstruct passphrase interactively:

    $ randall --wordlist de --words 4 --reconstruct
    reconstruct password consisting of 4 words using wordlist de

    type first word, complete with <tab>, press <enter> to add word
    > Schlange
    add word "Schlange" to password
    Schlange > Wald
    add word "Wald" to password
    SchlangeWald > Quasi
    add word "Quasi" to password
    SchlangeWaldQuasi > Toll
    add word "Toll" to password
    password is: SchlangeWaldQuasiToll
