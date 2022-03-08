package main

import (
	"log"

	"github.com/nooderg/autotest_templating/pkg/parsing"
)

func main() {
	_, err := parsing.ParseOpenApi()
	if err != nil {
		log.Println(err)
	}

	// log.Println(data)
}
