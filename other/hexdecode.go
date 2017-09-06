package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "give hex code")
	}

	out, err := base64.StdEncoding.DecodeString(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("decode data in hex is ", hex.EncodeToString(out))
}
