package main

import (
	basehttp "net/http"

	"gitlab.com/contextualcode/go-object-store/http"
	"gitlab.com/contextualcode/go-object-store/store"
)

func main() {
	config, err := store.LoadConfig("./config.yaml")
	if err != nil {
		panic(err)
	}
	basehttp.Handle("/", basehttp.FileServer(basehttp.Dir("./")))
	if err := http.Listen(config); err != nil {
		panic(err)
	}

}
