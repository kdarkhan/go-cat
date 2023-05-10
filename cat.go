package main

import flag "github.com/spf13/pflag"
import "fmt"
import "os"
import "strings"

func main() {

	showAll := flag.BoolP("show-all", "A", false, "equivalent to -vET")
	numberNonEmpty := flag.BoolP("number-nonblank", "b", false, "number nonempty output lines, overrides -n")
	showNonPrintingAndEnds := flag.BoolP("e", "e", false, "equivalent to -vE")
	showEnds := flag.BoolP("show-ends", "E", false, "display $ at end of each line")
	numberAll := flag.BoolP("number", "n", false, "number all output lines")
	squezeBlank := flag.BoolP("squeeze-blank", "s", false, "suppress repeated empty output lines")
	showNonPrintingAndTabs := flag.BoolP("t", "t", false, "equivalent to -vT")
	showTabs := flag.BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	_ = flag.BoolP("u", "u", false, "(ignored)")
	showNonPrinting := flag.BoolP("show-nonprinting", "v", false, "use ^ and M- notation, except for LFD and TAB")
	helpFlag := flag.Bool("help", false, "display this help and exit")
	verFlag := flag.Bool("version", false, "output version information and exit")

	flag.CommandLine.SortFlags = false

	flag.Parse()

	if *verFlag {
		printVersionInformation()
	}

	if *helpFlag {
		fmt.Printf(strings.TrimLeft(usagePre, "\n") + "\n")
		flag.CommandLine.PrintDefaults()
		fmt.Printf(usagePost)
	}

	if *numberNonEmpty {
		*numberAll = false
	}

	if *showNonPrintingAndTabs {
		*showNonPrinting = true
		*showTabs = true
	}

	if *showNonPrintingAndEnds {
		*showNonPrinting = true
		*showEnds = true
	}

	if *showAll {
		*showNonPrinting = true
		*showTabs = true
		*showEnds = true
	}

	cat(os.Args[len(os.Args)-flag.NArg():], *numberNonEmpty, *numberAll, *showEnds, *showNonPrinting, *showTabs, *squezeBlank)
}

func cat(inputs []string, numberNonEmpty, numberAll, showEnds, showNonPrinting, showTabs, squezeBlank bool) {
}

const usagePre string = `
Usage: go-cat [OPTION]... [FILE]...
Concatenate FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.
`

const usagePost string = `
Examples:
  go-cat f - g  Output f's contents, then standard input, then g's contents.
  go-cat        Copy standard input to standard output.
`

const version string = `
go-cat 0.1
Hello world
`

func printVersionInformation() {
	fmt.Printf(strings.TrimLeft(version, "\n"))
}
