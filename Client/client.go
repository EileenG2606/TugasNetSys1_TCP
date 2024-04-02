package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

func menu() {
	scanner := bufio.NewScanner(os.Stdin) //Buat declare scanenr
	for {
		fmt.Println("Welcome user!")
		fmt.Println("1.Send Message")
		fmt.Println("2.Exit")
		fmt.Println("Please enter your choice: ")
		scanner.Scan() //Tempat menerima input
		cho := scanner.Text()
		if cho == "1" {
			sendMessageMenu()
		} else if cho == "2" {
			fmt.Println("Thank you for using this program")
			break
		} else {
			fmt.Println("Please choose one of the menu")
		}
	}
}

func sendMessageMenu() {
	scanner := bufio.NewScanner(os.Stdin)
	var message string
	for {
		fmt.Print("Please enter your message: ")
		scanner.Scan()
		message = scanner.Text()
		if len(message) < 4 {
			fmt.Println("Message cannot be less than 4 characters")
		} else {
			break
		}
	}
	sendMessagetoServer(message)
}

func sendMessagetoServer(message string) {
	serverConn, err := net.DialTimeout("tcp", "127.0.0.1:1456", 3*time.Second)
	if err != nil {
		panic(err)
	}
	defer serverConn.Close()

	err = binary.Write(serverConn, binary.LittleEndian, uint32(len(message)))
	if err != nil {
		panic(err)
	}

	//Set Write Deadline untuk menetapkan batas waktu untuk operasi penulisan
	err = serverConn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		panic(err)
	}

	_, err = serverConn.Write([]byte(message))
	if err != nil {
		// Handle write timeout
		if netErr, yes := err.(net.Error); yes && netErr.Timeout() {
			fmt.Println("Write timeout")
			return
		}
		panic(err)
	}

	var size uint32

	// Set Read Deadline untuk menetapkan batas waktu untuk operasi pembacaan 
	err = serverConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		panic(err)
	}

	err = binary.Read(serverConn, binary.LittleEndian, &size)
	if err != nil {
		// Handle read timeout
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("Read timeout")
			return
		}
		panic(err)
	}

	bytReply := make([]byte, size)
	_, err = serverConn.Read(bytReply)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Replied: %s\n", string(bytReply))
}

func main() {
	menu()
}
