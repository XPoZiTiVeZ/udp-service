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
	fileName := ""
	if len(os.Args) < 2 {
		fmt.Println("Предоставьте ip:порт и имя файла в формате <программа ip-адрес:порт имя_файла>")
		return
	} else if len(os.Args) == 2 {
		fileName = "file.txt"
	} else {
		fileName = os.Args[2]
	}

	addr_port := os.Args[1]

	addr, err := net.ResolveUDPAddr("udp", addr_port)
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}

	fmt.Println("Подключён к", addr_port)

	_, err = conn.Write([]byte("Greetings"))
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}

	var id string
	var part int
	var fileData bytes.Buffer
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	for {
		buffer := make([]byte, 65536)
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Ошибка", err)
			return
		}

		var bytesLength uint32
		buf := bytes.NewBuffer(buffer[:4])
		err = binary.Read(buf, binary.LittleEndian, &bytesLength)
		fmt.Println("От сервера", addr_port, "получено", bytesRead, "байт")

		if err == nil && bytesRead-4 == int(bytesLength) {
			if bytesLength == 0 {
				os.WriteFile(fileName, fileData.Bytes(), 0644)
				return
			}

			fileData.Write(buffer[4:bytesRead])
			part += 1
		} else {
			id = string(buffer[:bytesRead])
		}

		_, err = conn.Write([]byte(fmt.Sprint(id, " ", part)))
		if err != nil {
			fmt.Println("Ошибка", err)
			return
		}
	}
}
