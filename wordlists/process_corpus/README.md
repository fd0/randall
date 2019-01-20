This program processes a file listing the most common words in a language as
downloaded from the [Leipzig Corpora Collection](https://wortschatz.uni-leipzig.de/en/download/)

To recreate the German word list, download the "Mixed Typical" corpus of 1
million words collected in 2011:
http://pcai056.informatik.uni-leipzig.de/downloads/corpora/deu_mixed-typical_2011_1M.tar.gz

Extract the file `deu_mixed-typical_2011_1M-words.txt`, then run the program
and supply the file as the single argument:

    go run main.go deu_mixed-typical_2011_1M-words.txt > ../corpus_leipzig_german_words.txt
