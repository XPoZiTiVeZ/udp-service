package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
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
	var part int
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
				return
			}
			f, _ := os.OpenFile("file.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			f.Write(buffer[4:bytesread])
			f.Close()
		}
		if n, _ := fmt.Scan(&part); n == 0 {
			return
		}
		for part != 0 && part != 1 {
			if n, _ := fmt.Scan(&part); n == 0 {
				return
			}
		}

		_, err = conn.Write([]byte(fmt.Sprint(id, part)))
		if err != nil {
			fmt.Println("Ошибка", err)
			return
		}
	}
}