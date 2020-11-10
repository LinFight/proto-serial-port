package main

import (
	"flag"
	"fmt"
	"github.com/LinFight/proto-serial-port/pb"
	"github.com/LinFight/proto-serial-port/protocol"
	"github.com/golang/protobuf/proto"
	"github.com/jacobsa/go-serial/serial"
	serialList "go.bug.st/serial"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var portListen = flag.String("p", "", "serial port to listen")
var baudRate = flag.Uint("rate", 9600, "the rate for baud")

func showSerialPortList() {
	ports, err := serialList.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		if strings.Contains(port, "Bluetooth") {
			continue
		}
		fmt.Printf("Found port: %v\n", port)
	}
	if *portListen == "" {
		fmt.Println("Please use -p to choose listen port")
		os.Exit(404)
	}
}

func serialPortInit() (io.ReadWriteCloser, error) {
	flag.Parse()
	// 输出命令行参数
	fmt.Println("portListen=", *portListen)
	showSerialPortList()
	options := serial.OpenOptions{
		PortName:        *portListen,
		BaudRate:        *baudRate,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	return port, nil
}
func main() {
	done := make(chan bool)
	port, _ := serialPortInit()
	readerChannel := make(chan []byte, 16)
	writeChannel := make(chan []byte, 16)
	go handleServer(port, readerChannel)
	go handleClient(port, writeChannel)
	go reader(readerChannel)
	// 往chan里面写入
	for i := 0; i < 1000; i++ {
		user1 := pb.User{
			Id:   *proto.Int32(1),
			Name: *proto.String("Mike"),
		}

		user2 := pb.User{
			Id:   2,
			Name: "John",
		}

		users := pb.MultiUser{
			Users: []*pb.User{&user1, &user2},
		}
		data, err := proto.Marshal(&users)
		if err != nil {
			log.Fatalln("Marshal data error: ", err)
		}
		writeChannel <- data
		time.Sleep(time.Second)
	}
	defer port.Close()
	<-done
}

func handleClient(conn io.ReadWriteCloser, writeChannel chan []byte) {
	for {
		data := <-writeChannel
		//time.Sleep(1 * time.Second)
		conn.Write(protocol.Packet(data))
	}
}

func handleServer(conn io.ReadWriteCloser, readerChannel chan []byte) {
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)
	//声明一个管道用于接收解包的数据
	//readerChannel := make(chan []byte, 16)
	//go reader(readerChannel)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		fmt.Println(buffer)
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
			fmt.Println(string(data))
		}
	}
}
