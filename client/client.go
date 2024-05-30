package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"
)

var fileString = flag.String("files", "", "files split by comma")

func init() {
	flag.Parse()
}

func main() {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 5150,
	})
	if err != nil {
		log.Fatalf("error with %v", err)
	}

	// be sure to close connections
	defer conn.Close()

	files := strings.Split(*fileString, ",")
	for _, targetFile := range files {
		data, err := os.ReadFile(targetFile)
		if err != nil {
			log.Printf("could not read '%s': %v", targetFile, err)
			continue
		}

		n, err := conn.Write(data)
		if err != nil {
			log.Printf("error while writing data to connection, %d\n", err)
			return
		}

		log.Printf("written '%s' with %v bytes", targetFile, n)
	}

}
