// script.go
package main

import (
	"fmt"

	_ "micscript/calc"
	_ "micscript/engineauth"
	_ "micscript/extable"
	_ "micscript/filter"
	_ "micscript/numgo"
	_ "micscript/regression"
	_ "micscript/statistic"
)

func main() {
	fmt.Println("Hello World!")
}
