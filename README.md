# go-cat

https://github.com/kdarkhan/go-cat

Implementing `cat` from `coreutils` using Golang.

This implementation tries to be close to `cat` supporting all of its
command line arguments.

`cat --help` and `go-cat --help` is almost the same on my Linux machine.

The tests in `cat_test.go` execute `cat` and `go-cat` binaries and
verify that the output of this Go implementation is the same as of `cat`.

It supports the following command line parameters:

```bash
$ ./go-cat --help
Usage: go-cat [OPTION]... [FILE]...
Concatenate FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.

  -A, --show-all           equivalent to -vET
  -b, --number-nonblank    number nonempty output lines, overrides -n
  -e, --e                  equivalent to -vE
  -E, --show-ends          display $ at end of each line
  -n, --number             number all output lines
  -s, --squeeze-blank      suppress repeated empty output lines
  -t, --t                  equivalent to -vT
  -T, --show-tabs          display TAB characters as ^I
  -u, --u                  (ignored)
  -v, --show-nonprinting   use ^ and M- notation, except for LFD and TAB
      --help               display this help and exit
      --version            output version information and exit

Examples:
  go-cat f - g  Output f's contents, then standard input, then g's contents.
  go-cat        Copy standard input to standard output.
```
