package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "<ip:port>")
		return
	}

	network := os.Args[1]
	addr, err := net.ResolveUDPAddr("udp", network)
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}

	part := 0
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}

	fmt.Println("Connected to", network)

	_, err = conn.Write([]byte("Greetings"))
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}

	var id int
	var file_data bytes.Buffer
	for {
		buffer := make([]byte, 65536)
		bytesread, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Ошибка", err)
			return
		}

		buf := bytes.NewBuffer(buffer[:4])
		fmt.Println(buf.Bytes())

		var data uint32
		if n, _ := fmt.Sscan(string(buffer[:bytesread]), &id); n != 1 {
			binary.Read(buf, binary.LittleEndian, &data)

			if data == 0 {
				os.WriteFile("file.txt", file_data.Bytes(), 0644)
				return
			}
			file_data.Write(buffer[4:bytesread])
			part += 1
		}

		_, err = conn.Write([]byte(fmt.Sprint(id, part)))
		if err != nil {
			fmt.Println("Ошибка", err)
			return
		}
		time.Sleep(time.Second)
	}
}