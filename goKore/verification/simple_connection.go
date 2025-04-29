package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// Connect to the server
	fmt.Println("Connecting to server: 192.168.5.219:6900")
	conn, err := net.DialTimeout("tcp", "192.168.5.219:6900", 5*time.Second)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()
	
	fmt.Println("Connected to server successfully")
	
	// Send a simple packet (just to test the connection)
	_, err = conn.Write([]byte{0x64, 0x00}) // 0x0064 is the master_login packet ID
	if err != nil {
		fmt.Printf("Failed to send data: %v\n", err)
		return
	}
	
	fmt.Println("Data sent successfully")
	
	// Try to receive a response
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Failed to receive data: %v\n", err)
		return
	}
	
	fmt.Printf("Received %d bytes of data\n", n)
	fmt.Println("Connection test successful")
}
