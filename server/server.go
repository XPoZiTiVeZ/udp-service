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
	if len(os.Args) != 3 {
        fmt.Println("Предоставьте ip:порт и имя файла в формате <программа ip:port file_name>")
    }

    network := os.Args[1]
    file_name := os.Args[2]
    client_table := make(map[int]int)

    file_data, err := os.ReadFile(file_name)
    if err != nil {
        fmt.Println("Ошибка", err)
        return
    }

    id := 0
    buflen := 8192
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

    fmt.Println("Сервер включён и ожидает подключение по порту", strings.Split(network, ":")[1])

    for {
        buffer := make([]byte, 65536)
        _, addr, _ := conn.ReadFromUDP(buffer)
        request := strings.TrimSpace(string(buffer))

        var cid, cn int
        if n, _ := fmt.Sscanf(request, "%d %d", &cid, &cn); n != 2 {
            _, err := conn.WriteToUDP([]byte(fmt.Sprint(id)), addr)
            if err != nil {
                fmt.Println("Ошибка", err)
            }
            client_table[id] = 0
            id += 1
        } else if buflen * cn >= len(file_data) {
            _, err := conn.WriteToUDP([]byte{0, 0, 0, 0}, addr)
            if err != nil {
                fmt.Println("Ошибка", err)
            }
        } else {
            if cn == 1 {
                client_table[cid] += 1
            }

            data, s, e := slice(file_data, buflen * client_table[cid], buflen * (client_table[cid] + 1))
            buf := bytes.NewBuffer([]byte{})
            binary.Write(buf, binary.LittleEndian, []uint32{uint32(len(data))})
            buf.Write(data)
            
            fmt.Println("Отправлено клиенту", cid, s, "-", e)

            _, err := conn.WriteToUDP(buf.Bytes(), addr)
            if err != nil {
                fmt.Println("Ошибка", err)
            }
        }
    }
}