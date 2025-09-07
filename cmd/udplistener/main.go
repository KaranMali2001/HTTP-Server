package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.ResolveUDPAddr("udp", ":8080")
	listenr, err := net.ListenUDP("udp", conn)

	if err != nil {
		fmt.Println("error while reading from tcp stream", err)
	}
	fmt.Println("server is listening at ", listenr.LocalAddr())
	buf := make([]byte, 1024)

	defer listenr.Close()
	for {
		// conn, err := listenr.Accept()
		n, _, err := listenr.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error while reading from UDP")
		}
		if n > 0 {
			fmt.Println("MESSAGE", string(buf[:n]))
		}

	}
}
