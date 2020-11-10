//服务端解包过程
package main

import (
	"fmt"
	"github.com/LinFight/proto-serial-port/protocol"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	"os"
)

func main() {
	done := make(chan bool)
	options := serial.OpenOptions{
		PortName:        "/dev/cu.usbserial-A1085TXF",
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

	go handleConnection(port)

	<- done
}
func handleConnection(conn io.ReadWriteCloser) {
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)
	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}
func reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			Log(string(data))
		}
	}
}
func Log(v ...interface{}) {
	fmt.Println(v...)
}
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
