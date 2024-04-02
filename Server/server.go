package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:1456")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleServerConnection(clientConn)
	}
}

func handleServerConnection(client net.Conn) {
	defer client.Close()

	var size uint32
	err := binary.Read(client, binary.LittleEndian, &size)
	if err != nil {
		panic(err)
	}

	// Set read deadline
	err = client.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		panic(err)
	}

	bytMsg := make([]byte, size)
	_, err = client.Read(bytMsg)
	if err != nil{
		if netErr, yes := err.(net.Error); yes && netErr.Timeout() {
			fmt.Println("Read timeout")
			return
		}
		panic(err)
	}

	strMsg := string(bytMsg)
	fmt.Printf("Received: %s\n", strMsg)

	var reply string
	if strings.HasSuffix(strMsg, ".zip") {
		reply = "File has been received"
	} else {
		reply = "Message has been received"
	}

	err = binary.Write(client, binary.LittleEndian, uint32(len(reply)))
	if err != nil {
		panic(err)
	}
	
	_, err = client.Write([]byte(reply))
	if err != nil {
		panic(err)
	}

}