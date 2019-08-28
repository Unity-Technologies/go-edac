package main

import (
	"fmt"
	"github.com/multiplay/go-edac/lib/edac"
	"log"
)

func main() {
	mcs, err := edac.MemoryControllers()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range mcs {
		i, err := c.Info()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%#v\n", i)
		ranks, err := c.DimmRanks()
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range ranks {
			fmt.Printf("%#v\n", r)
		}
	}
}
