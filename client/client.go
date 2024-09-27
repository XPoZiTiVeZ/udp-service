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
	}

	address := os.Args[1]
	fileName = os.Args[2]

	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}

	fmt.Println("Подключён к", address)

	_, err = conn.Write([]byte("Greetings"))
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}

	var id int
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

		buf := bytes.NewBuffer(buffer[:4])
		fmt.Println("От сервера", address, "получено", bytesRead, "байт")

		var bytesLength uint32
		if n, _ := fmt.Sscan(string(buffer[:bytesRead]), &id); n != 1 {
			binary.Read(buf, binary.LittleEndian, &bytesLength)

			if bytesLength == 0 {
				os.WriteFile(fileName, fileData.Bytes(), 0644)
				return
			}

			fileData.Write(buffer[4:bytesRead])
			part += 1
		}

		_, err = conn.Write([]byte(fmt.Sprint(id, part)))
		if err != nil {
			fmt.Println("Ошибка", err)
			return
		}
	}
}
