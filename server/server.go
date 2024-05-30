package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 5150,
	})

	if err != nil {
		log.Fatalf("error with %v", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed by the client")
				return
			}
			log.Printf("error reading from connection: %v", err)
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, conn)

	if err != nil {
		if err == io.EOF {
			log.Println("Connection closed by the client")
			return
		}
		log.Printf("error reading from connection: %v", err)
		return
	}

	os.WriteFile("out.mkv", buf.Bytes(), 0664)
	log.Printf("incomming message: %v", buf)
}
