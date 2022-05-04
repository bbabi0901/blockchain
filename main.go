package main

import (
	"github.com/bbabi0901/blockchain/explorer"
	"github.com/bbabi0901/blockchain/rest"
)

func main() {
	go explorer.Start(3000) // goroutine을 사용해서 explorer and rest를 동시에 실행?
	// causes err bc same url(?) -> change handler nil, default serveMux, to my own handler in ListenAndServe()

	rest.Start(4000)
}
