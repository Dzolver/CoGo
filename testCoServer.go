package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var count = 0
var processLimitChan = make(chan int, 1000)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port")
		return
	}
	PORT := ":" + arguments[1]
	var wg sync.WaitGroup
	wg.Add(1)
	go tcpListener(PORT)
	wg.Add(1)
	go udpListener(PORT)
	wg.Wait()
}
func handleTCPConnection(clientConnection net.Conn) {
	fmt.Print(".")
	for {
		netData, err := bufio.NewReader(clientConnection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		rcvMessage := strings.TrimSpace(string(netData))
		if rcvMessage == "STOP" {
			fmt.Println("Client connection has exited")
			break
		}
		counter := strconv.Itoa(count) + "\n"
		clientConnection.Write([]byte(string(counter)))
	}
	clientConnection.Close()
}
func handleUDPConnection(listenerConnection *net.UDPConn, processLimitChan chan int) {
	fmt.Print("*" + string(count))
	buffer := make([]byte, 1024)
	for {
		n, clientAddress, err := listenerConnection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Received message from client: " + string(buffer[0:n-1]))
		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("UDP client has exited!")
			break
		}
		data := []byte("UDP server acknowledges!")
		v, err := listenerConnection.WriteToUDP(data, clientAddress)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(v)
		count++
		processLimitChan <- count
	}
}
func tcpListener(PORT string) {
	listenerConnection, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fmt.Println("done!")
	defer listenerConnection.Close()
	for {
		clientConnection, err := listenerConnection.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleTCPConnection(clientConnection)
		count++
	}
}

func udpListener(PORT string) {
	buffer := make([]byte, 1024)
	serverAddress, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	listenerConnection, err := net.ListenUDP("udp4", serverAddress)
	defer listenerConnection.Close()
	for {
		n, clientAddress, err := listenerConnection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Received message from client: " + string(buffer[0:n-1]))
		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("UDP client has exited!")
			break
		}
		data := []byte("UDP server acknowledges!")
		v, err := listenerConnection.WriteToUDP(data, clientAddress)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(v)
		count++
		processLimitChan <- count
	}
}
