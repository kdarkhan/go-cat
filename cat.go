package main

import flag "github.com/spf13/pflag"
import "fmt"
import "strings"


func main() {

		verFlag := flag.BoolP("version", "v", false, "output version information and exit")

		flag.CommandLine.SortFlags = false
		flag.Parse()

		if *verFlag {
			printVersionInformation()
		}

}


const version string =
`
go-cat 0.1
Hello world
`

func printVersionInformation() {
	fmt.Printf(strings.TrimLeft(version, "\n"))
}
