package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"sync"

	_ "modernc.org/sqlite"
)

var mx sync.Mutex

type Request struct {
    method string
    addr   *net.UDPAddr
}

func main() {
    if len(os.Args) < 3 {
        panic("Предоставьте порт и имя файла в формате <программа порт имя_файла>")
    }

    port := os.Args[1]
    fileName := os.Args[2]

    db, err := sql.Open("sqlite", fileName)
    if err != nil { panic(err) }

    defer db.Close()

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS counter(id INT PRIMARY KEY, value INT)`)
    if err != nil { panic(err) }

    rows, err := db.Query(`SELECT value FROM counter WHERE id = 1`)
    if err != nil { panic(err) }

    var value int = -1
    rows.Next()
    rows.Scan(&value)
    rows.Close()

    if value == -1 {
        _, err = db.Exec(`INSERT INTO counter VALUES (1, 0)`)
        if err != nil { panic(err) }
    }

    addr, err := net.ResolveUDPAddr("udp", fmt.Sprint("0.0.0.0:", port))
    if err != nil { panic(err) }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil { panic(err) }

    var query []Request
    go func(){
        for {
        buffer := make([]byte, 65536)
        n, addr, _ := conn.ReadFromUDP(buffer)

        value := string(buffer[:n])
    
        if value == "get" || value == "inc" {
            mx.Lock()
            query = append(query, Request{value, addr})
            mx.Unlock()
        } else {
            conn.WriteTo([]byte("error"), addr)
        }
        }
    }()

    for {
        if len(query) == 0 {
            continue
        }

        req := query[0]
        
        mx.Lock()
        query = query[1:]
        mx.Unlock()

        if req.method == "get" {
            rows, err := db.Query(`SElECT value FROM counter WHERE id = 1`)
            if err != nil { panic(err) }

            
            var value int
            if rows.Next() {
                rows.Scan(&value)
                rows.Close()
            }

            conn.WriteToUDP([]byte(fmt.Sprint(value)), req.addr)
            continue
        }
        
        if req.method == "inc" {
            rows, err := db.Query(`SElECT value FROM counter WHERE id = 1`)
            if err != nil { conn.WriteToUDP([]byte("fail"), req.addr) }

            var value int
            if rows.Next() {
                rows.Scan(&value)
                rows.Close()
            }
            if rows.Err() == sql.ErrNoRows { conn.WriteToUDP([]byte("fail"), req.addr)}

            _, err = db.Exec(`UPDATE counter SET value = ? WHERE id = 1`, value + 1)
            if err != nil { conn.WriteToUDP([]byte("fail"), req.addr) }

            _, err = conn.WriteToUDP([]byte([]byte("ok")), req.addr)
            if err != nil { panic(err) }
        }
    }
}
