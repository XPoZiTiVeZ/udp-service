package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

func slice(data []byte, s, e int) ([]byte, int, int) {
    if s < 0 {
        s = 0
    }

    if s > len(data) {
        s = len(data)
    }

    if e > len(data) {
        e = len(data)
    }

    if e < 0 {
        e = 0
    }

    return data[s:e], s, e
}

func main() {
	if len(os.Args) < 3 {
        fmt.Println("Предоставьте порт и имя файла в формате <программа порт имя_файла>")
    }

    port := os.Args[1]
    fileName := os.Args[2]

    fileData, err := os.ReadFile(fileName)
    if err != nil {
        fmt.Println("Ошибка", err)
        return
    }

    addr, err := net.ResolveUDPAddr("udp", fmt.Sprint("0.0.0.0:", port))
    if err != nil {
        fmt.Println("Ошибка", err)
        return
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Println("Ошибка", err)
        return
    }

    id := 0
    bytesLength := 32768
    for {
        buffer := make([]byte, 65536)
        bytesRead, addr, _ := conn.ReadFromUDP(buffer)
        request := strings.TrimSpace(string(buffer[:bytesRead]))

        var cid, cn int
        if n, _ := fmt.Sscanf(request, "%d %d", &cid, &cn); n != 2 {
            _, err := conn.WriteToUDP([]byte(fmt.Sprint(id)), addr)
            if err != nil {
                fmt.Println("Ошибка", err)
            }
            id += 1
        } else if bytesLength * cn >= len(fileData) || 0 > cid || cid > id {
            _, err := conn.WriteToUDP([]byte{0, 0, 0, 0}, addr)
            if err != nil {
                fmt.Println("Ошибка", err)
            }
        } else {
            data, s, e := slice(fileData, bytesLength * cn, bytesLength * (cn + 1))

            buf := bytes.NewBuffer([]byte{})
            binary.Write(buf, binary.LittleEndian, []uint32{uint32(len(data))})

            buf.Write(data)
            fmt.Println(s, "-", e, "байт отправлено клиенту", cid)

            _, err := conn.WriteToUDP(buf.Bytes(), addr)
            if err != nil {
                fmt.Println("Ошибка", err)
            }
        }
    }
}
