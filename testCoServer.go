package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var count = 0
var processLimitChan = make(chan int, 1000)
var wg sync.WaitGroup

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
			count--
			break
		}
		counter := strconv.Itoa(count) + "\n"
		clientConnection.Write([]byte(string(counter)))
	}
	clientConnection.Close()
}

func tcpListener(PORT string) {
	listenerConnection, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
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
	if err != nil {
		fmt.Println(err)
		return
	}
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
			count--
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

func findMongoDataExample(context context.Context, mongoClient *mongo.Client) {
	if err := mongoClient.Ping(context, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Mongo DB client has been pinged successfully!")
	database := mongoClient.Database("test")
	testCollection := database.Collection("test")
	filterCursor, err := testCollection.Find(context, bson.M{"score": 500})
	if err != nil {
		fmt.Println(err)
		return
	}
	var filterResult []bson.M
	if err = filterCursor.All(context, &filterResult); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Data found from DB: ", filterResult[0]["name"], " scored ", filterResult[0]["score"], " points!")
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port")
		return
	}
	//mongoDB specs
	uri := "mongodb+srv://admin:london1234@cluster0.8acnf.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

	//initialize mongoDB client
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(context, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	//disconnect mongoDB client on return
	defer func() {
		if err = mongoClient.Disconnect(context); err != nil {
			panic(err)
		}
	}()
	//sample mongoDB query
	findMongoDataExample(context, mongoClient)

	//PORT = ":20001"
	PORT := ":" + arguments[1]

	// add to sync.WaitGroup
	wg.Add(1)
	go tcpListener(PORT)
	wg.Add(1)
	go udpListener(PORT)
	wg.Wait()

}
