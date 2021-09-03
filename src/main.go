package main

import (
	"fmt"

	d7024e "github.com/maxlengdell/D7024E/d7024e"
)

func main() {
	fmt.Println("hello")
	d7024e.Listen("0.0.0.0", 8080)
}
