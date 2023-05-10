package main

import flag "github.com/spf13/pflag"
import "fmt"
import "os"
import "strings"
import "bufio"

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
		return
	}

	if *helpFlag {
		fmt.Printf(strings.TrimLeft(usagePre, "\n") + "\n")
		flag.CommandLine.PrintDefaults()
		fmt.Printf(usagePost)
		return
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

	failed := cat(os.Args[len(os.Args)-flag.NArg():], *numberNonEmpty, *numberAll, *showEnds, *showNonPrinting, *showTabs, *squezeBlank)
	if failed {
		// cat exists with error code 1 when file is not found, do the same
		os.Exit(1)
	}
}

func cat(inputs []string, numberNonEmpty, numberAll, showEnds, showNonPrinting, showTabs, squezeBlank bool) bool {
	if len(inputs) == 0 {
		inputs = []string{"-"}
	}
	failed := false

	lineNum := 1

	for _, input := range inputs {
		if input == "-" {
			catReader(bufio.NewReader(os.Stdin), &lineNum, numberNonEmpty, numberAll, showEnds, showNonPrinting, showTabs, squezeBlank)
		} else {
			file, err := os.Open(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cat: %s: No such file or directory\n", input)
				failed = true
				continue
			}
			catReader(bufio.NewReader(file), &lineNum, numberNonEmpty, numberAll, showEnds, showNonPrinting, showTabs, squezeBlank)
		}
	}
	return failed
}

func catReader(reader *bufio.Reader, lineNum *int, numberNonEmpty, numberAll, showEnds, showNonPrinting, showTabs, squezeBlank bool) {
	fileScanner := bufio.NewScanner(reader)

	fileScanner.Split(bufio.ScanLines)

	previousWasEmpty := false
	endLineSymbol := ""
	if showEnds {
		endLineSymbol = "$"
	}

	for fileScanner.Scan() {
		line := fileScanner.Text()
		isEmpty := line == ""

		if !isEmpty || !previousWasEmpty || !squezeBlank {
			if showTabs {
				line = strings.ReplaceAll(line, "\t", "^I")
			}
			if showNonPrinting {
				line = escapeString(line)
			}

			numberPrefix := ""
			if (numberNonEmpty && !isEmpty) || numberAll {
				numberPrefix = fmt.Sprintf("%6d	", *lineNum)
			}

			fmt.Printf("%s%s%s\n", numberPrefix, line, endLineSymbol)
			if !isEmpty || numberAll {
				*lineNum += 1
			}
		}

		previousWasEmpty = isEmpty
	}
}

// Replace non-printable characters using Caret notation
// https://en.wikipedia.org/wiki/C0_and_C1_control_codes
func escapeString(line string) string {
	var sb strings.Builder
	sb.Grow(len(line))

	for _, r := range line {
		int_r := int(r)
		if int_r >= 0 && int_r < 31 && int_r != 9 {
			sb.WriteString("^" + string(64+int_r))
		} else if int_r == 127 {
			sb.WriteString("^?")
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

func splitNewLine(line string) (string, string) {
	for i := len(line) - 1; i >= 0; i-- {
		if line[i] != '\r' && line[i] != '\n' {
			return line[:i+1], line[i+1:]
		}
	}

	return "", line
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
`

func printVersionInformation() {
	fmt.Printf(strings.TrimLeft(version, "\n"))
}
