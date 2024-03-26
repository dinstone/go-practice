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

	// 发送消息给服务器
	fmt.Fprintln(conn, "Hello, server!")

	go func() {
		// 发送心跳包以保持连接
		ticker := time.NewTicker(5 * time.Second)
		for t := range ticker.C {
			fmt.Printf("Sending heartbeat to server at %v\n", t)
			conn.Write([]byte("HEARTBEAT\n"))
		}
	}()

	for {
		// 读取服务器的响应
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		fmt.Print("Message from server:", message)
	}

	// 阻塞主线程，以便保持连接打开
	select {}
}
