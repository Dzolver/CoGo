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

func handleTCPConnection(clientConnection net.Conn, cxt context.Context, mongoClient *mongo.Client) {
	fmt.Print(".")
	clientResponse := "DEFAULT"
	for {
		netData, err := bufio.NewReader(clientConnection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		data := strings.TrimSpace(string(netData))
		packetCode := data[0:2]
		rcvMessage := strings.Replace(data, packetCode, "", -1)
		if packetCode == "L#" {
			fmt.Println("Login packet received!")
			username := strings.Split(rcvMessage, "?")[0]
			password := strings.Split(rcvMessage, "?")[1]
			clientResponse = handleLogin(username, password, mongoClient)
		}
		if packetCode == "R#" {
			fmt.Println("Register packet received!")
			username := strings.Split(rcvMessage, "?")[0]
			password := strings.Split(rcvMessage, "?")[1]
			//THIS NEEDS TO BE REDIRECTED TO A REGISTER FUNCTION
			clientResponse = handleRegistration(username, password, mongoClient)
		}
		fmt.Println("Sent message back to client : ", clientResponse)
		fmt.Println("Received message from client : ", rcvMessage)
		if rcvMessage == "STOP" {
			fmt.Println("Client connection has exited")
			count--
			break
		}
		clientConnection.Write([]byte(strings.Trim(strconv.QuoteToASCII(clientResponse), "\"")))
	}
	clientConnection.Close()
}

func tcpListener(PORT string, cxt context.Context, mongoClient *mongo.Client) {
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
		go handleTCPConnection(clientConnection, cxt, mongoClient)
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

func findMongoDataExample(cxt context.Context, mongoClient *mongo.Client) {
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Mongo DB client has been pinged successfully!")
	database := mongoClient.Database("test")
	testCollection := database.Collection("test")
	filterCursor, err := testCollection.Find(cxt, bson.M{"score": 500})
	if err != nil {
		fmt.Println(err)
		return
	}
	var filterResult []bson.M
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Data found from DB: ", filterResult[0]["name"], " scored ", filterResult[0]["score"], " points!")
}
func handleLogin(username string, password string, mongoClient *mongo.Client) string {
	if !lookForUser(username, mongoClient) {
		if validateUser(username, password, mongoClient) {
			fmt.Println("Login successful!")
			return "Login successful"
		} else {
			fmt.Println("Login failed!")
			return "Login Failed"
		}
	} else {
		if validateUser(username, password, mongoClient) {
			fmt.Println("Login successful!")
			return "Login Successful"
		} else {
			fmt.Println("Login failed!")
			return "Login failed"
		}
	}
}
func handleRegistration(username string, password string, mongoClient *mongo.Client) string {
	if !lookForUser(username, mongoClient) {
		createUser(username, password, mongoClient)
		return "Account created"
	} else {
		fmt.Println(username, " is not available")
		return "Username is not available"
	}
}
func createUser(username string, password string, mongoClient *mongo.Client) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database := mongoClient.Database("player")
	users := database.Collection("users")
	createResult, err := users.InsertOne(cxt, bson.D{
		{"username", username},
		{"password", password},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New user added to db : ", createResult.InsertedID)
}
func lookForUser(username string, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database := mongoClient.Database("player")
	users := database.Collection("users")
	filterCursor, err := users.Find(cxt, bson.M{"username": username})
	if err != nil {
		fmt.Println(err)
		return false
	}
	var filterResult []bson.M
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	if len(filterResult) == 0 {
		//no users found
		return false
	} else if len(filterResult) == 1 {
		return true
	}
	return true
}
func validateUser(username string, password string, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	users := database.Collection("users")
	filterCursor, err := users.Find(cxt, bson.M{"username": username, "password": password})
	if err != nil {
		fmt.Println(err)
		return false
	}
	var filterResult []bson.M
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	if len(filterResult) == 0 {
		//no users found
		return false
	} else if len(filterResult) == 1 {
		return true
	}
	return true
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
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(cxt, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	//disconnect mongoDB client on return
	defer func() {
		if err = mongoClient.Disconnect(cxt); err != nil {
			panic(err)
		}
	}()
	//sample mongoDB query
	findMongoDataExample(cxt, mongoClient)

	//PORT = ":20001"
	PORT := ":" + arguments[1]

	// add to sync.WaitGroup
	wg.Add(1)
	go tcpListener(PORT, cxt, mongoClient)
	wg.Add(1)
	go udpListener(PORT)
	wg.Wait()

}
