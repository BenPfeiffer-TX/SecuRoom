package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	input = strings.TrimSpace(input)
	fmt.Printf("sending: %s\n", input)

	data := []byte(input)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()
}
