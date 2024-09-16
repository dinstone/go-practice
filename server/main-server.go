package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("New connection established")

	go func() {
		// 假设我们每隔一段时间发送一个心跳包以保持连接
		ticker := time.NewTicker(5 * time.Second)
		for t := range ticker.C {
			fmt.Printf("Sending heartbeat at %v\n", t)
			conn.Write([]byte("HEARTBEAT\n"))
		}
	}()

	// 读取客户端发送的数据
	reader := bufio.NewReader(conn)
	for {
		var readBytes [1024]byte
		n, err := reader.Read(readBytes[:])
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from connection:", err)
			}
			break
		}
		data := string(readBytes[:n])
		fmt.Print("Message received: ", data)

		// 发送回应给客户端
		conn.Write([]byte(data))
	}
}
