package main

import (
	"fmt"
	"github.com/LinFight/proto-serial-port/protocol"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	"time"
)

func sender(conn io.ReadWriteCloser) {
	for i := 0; i < 1000; i++ {
		words := "{\"Id\":1,\"Name\":\"golang\",\"Message\":\"message\"}"
		conn.Write(protocol.Packet([]byte(words)))
	}
	fmt.Println("send over")
}
func main() {
	options := serial.OpenOptions{
		PortName:        "/dev/cu.usbserial-AI06APD9",
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer port.Close()
	go sender(port)
	for {
		time.Sleep(1 * 1e9)
	}
}
