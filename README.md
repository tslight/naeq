 ![CI Result](https://github.com/tslight/naeq/actions/workflows/build.yml/badge.svg?event=push) [![Go Report Card](https://goreportcard.com/badge/github.com/tslight/naeq)](https://goreportcard.com/report/github.com/tslight/naeq) [![Go Reference](https://pkg.go.dev/badge/github.com/tslight/naeq.svg)](https://pkg.go.dev/github.com/tslight/naeq)
# English Qabalah CLI & REST API

*The Secret Cipher of the UFOnauts as an API & CLI, because* `¯\_(ツ)_/¯`

Inspired from [Allen
Greenfield's](https://en.wikipedia.org/wiki/Allen_H._Greenfield) bizarre &
fascinating
[book](https://www.amazon.co.uk/Complete-SECRET-CIPHER-UfOnauts/dp/171864535X)
& it's usage in the equally bizarre & fascinating
[Hellier](https://www.hellier.tv/) TV show, taking cues from
[Wren Collier](https://www.naeq.io/) & [Chad Milburn's](http://www.naequery.com/) web based
versions.

## Installation

``` shell
go install github.com/tslight/naeq@latest
```

Alternatively, download a suitable pre-compiled binary for your architecture
and operating system from the
[releases](https://github.com/tslight/naeq/releases) page and move it to
somewhere in your `$PATH`.

**N.B.**

The binaries are fairly large because of the embedding of all the files in
`assets/books/*.json` into them.

## CLI Usage

Running `alw-cli` with no arguments will prompt for words and use Liber Al
Vegis (The Book of the Law) by default, and will display all the results.

If you want to choose another book and limit the number of results use the `-b`
and `-n` flags respectively. For example:

``` shell
alw-cli -b liber-i -n 31 hellier
```

To list the baked in books run `alw-cli -l` and to give the tool your own book
in json format use the `-p` flag. If you just want to calculate the NAEQ/ALW
cipher sum of a given string (without finding matches in a mystical text) use
the `-s` flag.

``` text
Usage: alw-cli [options...] <words>:
  -b string
        embedded book (default "liber-al")
  -l    list embedded books
  -n int
        number of matches to show
  -p string
        path to alternative book
  -s    only return naeq sum
```

## API Usage

This is very much W.I.P. and my first time writing an API from scratch so bare
with me...

``` text
Usage: alw-api [options...]:
  -p int
        Port to listen on (default 8080)
  -v    print version info
```

I have instances deployed [here](https://naeq.onrender.com) and
[here](https://alw.up.railway.app/), if you'd like to try it out:

``` shell
curl -X GET https://naeq.onrender.com?words=hellier
curl -X GET https://naeq.onrender.com?words=hellier&book=liber-i.json
curl -X POST https://naeq.onrender.com -d '{"words": "hellier"}'
curl -X POST https://naeq.onrender.com -d '{"book": "liber-x.json", "words": "hellier"}'
```

``` shell
curl -X GET https://alw.up.railway.app?words=hellier
curl -X GET https://alw.up.railway.app?words=hellier&book=liber-i.json
curl -X POST https://alw.up.railway.app -d '{"words": "hellier"}'
curl -X POST https://alw.up.railway.app -d '{"book": "liber-x.json", "words": "hellier"}'
```

## Context/Background

The sections below were ripped directly from [here](https://www.naeq.io/about),
so full credit/attribution to [Wren Collier](https://liminalroom.com/) & Alynne
Keith.

### Gematria and Qabalah

Hermetic Qabalah is a derivative of a school of Jewish mysticism called
Kabbalah that was developed by modern Western esotericists as a way to explore
the divine and the nature of the universe. One aspect of Qabalah is Gematria, a
mystical interpretation of a holy text using specific mathematical laws. A
Gematria is a system used to assign a numeric value for each letter of a word,
which is then summed. This sum can be referred to as the “key” of that word, or
phrase. Words and phrases that have the same numeric value are thought to have
similar properties and can be used to meditate on hidden meanings or
relationships contained within those similarities.

### The Book of The Law

The Book of the Law — or Liber AL vel Legis — is the central holy text of
Thelema, a spiritual and social philosophy derived from Western esotericism and
founded by magician Aleister Crowley. Liber AL vel Legis was dictated to
Crowley over the course of three days in 1904 by a discarnate entity called
Aiwass.

### The New Aeon English Qabalah

Since Hebrew Gematria uses the Hebrew script to derive values, it does not
neccessarily apply as well to English or Roman scripts and thus efforts have
been made over the years to develop an “English Qabalah” that could be utilized
with texts in those languages.

In Chapter II verse 55 of Liber AL vel Legis, Crowley writes: “Thou shalt
obtain the order & value of the English Alphabet, thou shalt find new symbols
to attribute them unto”. When analyzing this information later, Crowley
realized that this implied there was a cipher contained within Liber AL that
had yet to be discovered or developed. Later on, in Chapter III sheet 16 of the
original documents, there is a unique page that contains a grid of numbers, a
slash through it, and a circle within a cross.

This page in particular is where the ALW cipher, or New Aeon English Qabalah,
is derived. James Lees discovered in 1976 that this key was based on the
magical number 11. By taking every eleventh letter of the alphabet as the
order, and then assigning them sequential values you are able to arrive at this
cipher:

``` text
  A=1 L=2 W=3 H=4 S=5 D=6 O=7 Z=8 K=9 V=10 G=11 R=12 C=13 N=14 Y=15 J=16 U=17
  F=18 Q=19 B=20 M=21 X=22 I=23 T=24 E=25 P=26
```

``` text
  A=1 B=20 C=13 D=6 E=25 F=18 G=11 H=4 I=23 J=16 K=9 L=2 M=21 N=14 O=7 P=26 Q=19
  R=12 S=5 T=24 U=17 V=10 W=3 X=22 Y=15 Z=8
```

### The Secret Cipher of the UFOnauts

Secret Cipher of the UFOnauts, a book published in 1994 and written by
occultist and ufologist Allen H. Greenfield seeks to help elucidate the nature
of the mysterious Aiwass that dictated Liber AL to Crowley, and the myriad
entities encountered by UFO contactees in the 20th and 21st century. His theory
was that ultraterrestrials - a term used by researcher and author John Keel
used to describe the beings - transmitted information to humans in a secret,
enciphered format, in the guise of enlightened messages. These entities often
used authors to produce voluminous works, such as Jane Roberts' “Seth” and the
“Ephraim” that communicated with poet James Merrill during his writing of The
Changing Light at Sandover. In particular, Greenfield conjectures that Liber AL
was the ultimate key to this secret cipher utilized by the ultraterrestrials.

Greenfield believed that by utilizing the ALW cipher/NAEQ you could take the
names of these entities, or information they provided and compare their cipher
value using the values of phrases or words from Liber AL as a key. By analyzing
words or phrases in Liber AL that had the same values, one could arrive at
further insight, or obtain major revelations about the nature of these entities
or their motivations.

### How to Use this Tool

Type a word or phrase into the search bar. The ALW cipher value will
automatically be calculated and a list of letters, words and phrases (of up to
16 words) from Liber AL of equal value will be shown.

How you interpret this data is up to you. We have found insight stringing
together multiple phrases or words with the same value that fit grammatically
into something resembling cut-up method poetry, as well as picking out phrases
that resonated with the name or phrase we were analyzing. Results are open to
interpretation and a bit of psychic intuition plays part. Use your imagination.

Do not read deeply into the order of the results displayed. Individual matches
are not neccessarily meant to be interpreted as one string of text; indeed,
results will often make little grammatical sense, though occasional happy
accidents occur.
