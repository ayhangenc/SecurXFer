/*

Bu args ile - basic

package main



import (
	"fmt"
	"os"
)

func main() {

	// baby step 1: parse the arguments from the command line...
	fmt.Println("len", len(os.Args))

	for _, arg := range os.Args[1:] {
		fmt.Println(arg)
	}
}


*/

// bu da parse paketi

package main

import (
	"flag"
	"fmt"
)

var (
	file   = flag.String("file", "", "file-name")
	count  = flag.Int("count", 2, "count params")
	repeat = flag.Bool("repeat", false, "Repeat execution")
)

func main() {
	flag.Parse()

	fmt.Println("file name: ", *file)
	fmt.Println("count: ", *count)
	fmt.Println("repeat: ", *repeat)
}
