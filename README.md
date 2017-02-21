# xkcd
A command line utility to list xkcd comics.

## Background

The idea for this program originates from *The Go Programing Language* book. I
decided to implement the program to teach myself some Go concepts and to make
a program that I would enjoy using!

## Usage

The `xkcd` program can look up individual xkcd comics.

    xkcd -n 323 // Lists metadata for the "Ballmer Peak" comic

It will list the most recent comic by default, if no number is specified.

To list the number, date, and title of the twenty most recent xkcd comics,
execute

    xkcd -l
