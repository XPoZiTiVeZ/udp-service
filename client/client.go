package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
    if len(os.Args) < 4 {
        panic("Предоставьте адрес, порт и комманду <программа ip-адрес порт get | inc>")
    }

    address := os.Args[1]
    port := os.Args[2]
    command := os.Args[3]

    addr, err := net.ResolveUDPAddr("udp", fmt.Sprint(address, ":", port))
    if err != nil { panic(err) }

    conn, err := net.DialUDP("udp", nil, addr)
    if err != nil { panic(err) }

    if len([]byte(command)) > 1000 { command = command[:1000] }

    _, err = conn.Write([]byte(command))
    if err != nil { panic(err) }

    buffer := make([]byte, 65536)
    n, err := conn.Read(buffer)
    if err != nil { panic(err) }

    fmt.Println(string(buffer[:n]))
}
