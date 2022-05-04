package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/bbabi0901/blockchain/explorer"
	"github.com/bbabi0901/blockchain/rest"
)

func usage() {
	fmt.Printf("Welcome to Quad-O Coin\n\n")
	fmt.Printf("Please use the following flags\n\n")
	fmt.Printf("-port=4000: 	Set the PORT of the server\n")
	fmt.Printf("-mode=rest: 	Choose between 'html' and 'rest'\n")
	os.Exit(0)
}

func Start() {
	// making command set for "rest"
	// rest := flag.NewFlagSet("rest", flag.ExitOnError) // command set은 flag가 여러개일 경우 더 유용하다.
	// portFlag := rest.Int("port", 4000, "Sets the port of the server")

	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Sets the port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}
	fmt.Println(*port, *mode)
}
