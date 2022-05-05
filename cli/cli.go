package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/bbabi0901/blockchain/explorer"
	"github.com/bbabi0901/blockchain/rest"
)

func usage() {
	fmt.Printf("Welcome to Quad-O Coin\n\n")
	fmt.Printf("Please use the following flags\n\n")
	fmt.Printf("-port=4000: 	Set the PORT of the server\n")
	fmt.Printf("-mode=rest: 	Choose between 'html' and 'rest'\n")
	// os.Exit(0) 이제는 defer를 실행시 exit하도록 다른 방법을 쓴다.
	runtime.Goexit() // 설명 읽어보면 Goexit()는 실행되기 전에 defer를 먼저 실행시켜준다 -> usage()를 통한 종료시 main.go의 db.Close()를 먼저 실행
}

func Start() {
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
