package main

import (
	"application-design/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		panic(err)
	}
}
