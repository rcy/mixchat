package main

import (
	"fmt"
	"gap/internal/server"
)

func main() {
	server := server.NewServer()

	fmt.Println("listening on", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
