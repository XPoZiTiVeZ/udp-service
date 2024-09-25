package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

func slice(data []byte, s, e int) []byte {
    if s < 0 {
        s = 0
    }

    if e > len(data) {
        e = len(data)
    }

    return data[s:e]
}

func main() {
	if len(os.Args) != 3 {
        fmt.Println("Предоставьте ip:порт и имя файла в формате <программа ip:port file_name>")
    }

    network := os.Args[1]
    file_name := os.Args[2]

    file_data, err := os.ReadFile(file_name)
    if err != nil {
        fmt.Println("Ошибка", err)
        return
    }

    id := 0
    buflen := 16384
    addr, err := net.ResolveUDPAddr("udp", network)
    if err != nil {
        fmt.Println("Ошибка", err)
        return
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Println("Ошибка", err)
        return
    }

    for {
        buffer := make([]byte, 65536)
        bytesread, addr, _ := conn.ReadFromUDP(buffer)
        request := strings.TrimSpace(string(buffer))
        fmt.Println(request)

        var cid, cn int
        if n, _ := fmt.Sscanf(request, "%d %d", &cid, &cn); n != 2 {
            fmt.Println(string(buffer[:bytesread]))
            _, err := conn.WriteToUDP([]byte(fmt.Sprint(id)), addr)
            if err != nil {
                fmt.Println("Ошибка", err)
            }
            id += 1
        } else if buflen * cn >= len(file_data) {
            _, err := conn.WriteToUDP([]byte{0, 0, 0, 0}, addr)
            if err != nil {
                fmt.Println("Ошибка", err)
            }
        } else {
            fmt.Println("Отправлено клиенту", cid)
            data := slice(file_data, buflen * cn, buflen * (cn + 1))
            fmt.Println(len(data), buflen * cn, buflen * (cn + 1))
            buf := bytes.NewBuffer([]byte{})
            binary.Write(buf, binary.LittleEndian, []uint32{uint32(len(data))})
            buf.Write(data)
            fmt.Println(buf.Bytes()[:16], len(data))
            _, err := conn.WriteToUDP(buf.Bytes(), addr)
            if err != nil {
                fmt.Println("Ошибка", err)
            }
        }
    }
}