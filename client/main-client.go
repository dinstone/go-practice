package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	server := "localhost:8080"
	conn, err := net.Dial("tcp", server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to server", server)

	go func() {
		// 发送心跳包以保持连接
		ticker := time.NewTicker(5 * time.Second)
		for t := range ticker.C {
			fmt.Printf("Sending heartbeat to server at %v\n", t)
			conn.Write([]byte("CLIENT HEARTBEAT\n"))
		}
	}()

	for {
		// 读取服务器的响应
		reader := bufio.NewReader(conn)
		var readBytes [1024]byte
		n, err := reader.Read(readBytes[:])
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		data := string(readBytes[:n])
		fmt.Print("Message from server: ", data)
	}

	// 阻塞主线程，以便保持连接打开
	select {}
}
