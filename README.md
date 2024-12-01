# Advent of Code 2024

Solutions for AoC 2024, in Go. You can build with just `go build` but the
makefile assumes you have `docker` and GNU `date` available.

## Building

You can build the project with `make` or `make build`. For a clean build you
do, as you've probably guessed, `make clean build`. Both create an `aoc2024`
binary in the `bin` folder.

## Running

The `aoc2024` binary has two commands: `input` and `run`. Both require a
specific day (so `1` for day 1, `2` for day 2 and so on until day 25). The
`input` command also requires a session token that you can set in the
`AOC_SESSION` environment variable or as a value to the `-s/--session`
parameter. See `aoc2024 --help` for more info.

For your convenience, `make run-all` runs the solutions for all available
puzzles and `make run DAY=XX` runs the solutions for day XX.

## Acknowledgements

The solver framework was largely inspired by [obalenenko's AoC package](https://github.com/obalunenko/advent-of-code).
