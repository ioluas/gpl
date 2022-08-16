# The Go Programming Language book

I have worked briefly with Go in the past, but never followed through since I mainly work with PHP & Javascript. This 
is where I pick it up again with the book "The Go Programming Language" taking it easy and going through the
examples and exercises in the book.

I've set this up to use [Cobra CLI](https://github.com/spf13/cobra) framework and attempt to solve the exercises 
in the book in single file command/sub-commands as I go.

---

### Command dedup

Simple command to count the number of occurrences of text line from input. Input can be `os.Stdin` or file(s)

```shell
# From os.Stdout
cat cmd/text.txt | ./gpl dedup
50      Some text
50      Followed by different text
50      Then more
50      And some more

./gpl dedup < data/dedup/text.txt
50      And some more
50      Some text
50      Followed by different text
50      Then more

# From file
./gpl dedup file data/dedup/text.txt
File:   data/dedup/text.txt
50      Some text
50      Followed by different text
50      Then more
50      And some more
```

---

### Command lissajous

Simple command to create GIF images with Lissajous curve graphs. sci-fi effect from the 1960s

```shell
cat ./gpl lissajous
```
This outputs a file `lissajous.gif` in current directory.

---

### Command fetch

Simple command to fetch given urls concurrently, reports status code and content length. Follows redirects keeping cookies