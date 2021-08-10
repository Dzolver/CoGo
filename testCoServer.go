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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type PlayerInventory struct {
	ObjectID  primitive.ObjectID `bson:"_id,omitempty"`
	User_id   string             `bson:"user_id,omitempty"`
	Purse     Purse              `bson:"purse,omitempty"`
	Equipment Equipment          `bson:"equipment,omitempty"`
	Items     Items              `bson:"items,omitempty"`
}

type Purse struct {
	Bits float32 `bson:"bits,omitempty"`
}
type Item struct {
	ObjectID    primitive.ObjectID `bson:"_id,omitempty"`
	Item_id     string             `bson:"item_id,omitempty"`
	Item_type   string             `bson:"item_type,omitempty"`
	Entity      string             `bson:"entity,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Num         int32              `bson:"num,omitempty"`
	Description string             `bson:"description,omitempty"`
	Stats       Stats              `bson:"stats,omitempty"`
}
type Stats struct {
	Health       float32 `bson:"health,omitempty"`
	Mana         float32 `bson:"mana,omitempty"`
	Attack       float32 `bson:"attack,omitempty"`
	MagicAttack  float32 `bson:"magicAttack,omitempty"`
	Defense      float32 `bson:"defense,omitempty"`
	MagicDefense float32 `bson:"magicDefense,omitempty"`
	Armor        float32 `bson:"armor,omitempty"`
	Evasion      float32 `bson:"evasion,omitempty"`
	Accuracy     float32 `bson:"accuracy,omitempty"`
	Agility      float32 `bson:"agility,omitempty"`
	Willpower    float32 `bson:"willpower,omitempty"`
	FireRes      float32 `bson:"fireRes,omitempty"`
	WaterRes     float32 `bson:"waterRes,omitempty"`
	EarthRes     float32 `bson:"earthRes,omitempty"`
	WindRes      float32 `bson:"windRes,omitempty"`
	IceRes       float32 `bson:"iceRes,omitempty"`
	EnergyRes    float32 `bson:"energyRes,omitempty"`
	NatureRes    float32 `bson:"natureRes,omitempty"`
	PoisonRes    float32 `bson:"poisonRes,omitempty"`
	MetalRes     float32 `bson:"metalRes,omitempty"`
	LightRes     float32 `bson:"lightRes,omitempty"`
	DarkRes      float32 `bson:"darkRes,omitempty"`
}
type ItemRange struct {
	Collection []Item `bson:"collection,omitempty"`
}
type Equipment struct {
	Head      ItemRange `bson:"Head,omitempty"`
	Body      ItemRange `bson:"Body,omitempty"`
	Feet      ItemRange `bson:"Feet,omitempty"`
	Weapon    ItemRange `bson:"Weapon,omitempty"`
	Accessory ItemRange `bson:"Accessory,omitempty"`
}
type Items struct {
	Consumable ItemRange `bson:"Consumable,omitempty"`
	Crafting   ItemRange `bson:"Crafting,omitempty"`
	Quest      ItemRange `bson:"Quests,omitempty"`
}
type TestItem struct {
	ObjectID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	User_id  string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Items    []string           `json:"items,omitempty" bson:"items,omitempty"`
}

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
		packetCode, packetMessage := packetDissect(netData)
		if packetCode == "L#" {
			fmt.Println("Login packet received!")
			username, password := processLoginPacket(packetMessage)
			clientResponse = handleLogin(username, password, mongoClient)
		}
		if packetCode == "R#" {
			fmt.Println("Register packet received!")
			username, password := processRegisterPacket(packetMessage)
			clientResponse = handleRegistration(username, password, mongoClient)
		}
		if packetCode == "I#" {
			fmt.Println("Inventory packet received!")
			// userID, itemType, itemID := processInventoryPacket(packetMessage)
			// addInventoryItem(userID, itemType, itemID, mongoClient)
		}
		fmt.Println("Sent message back to client : ", clientResponse)
		fmt.Println("Received message from client : ", packetMessage)
		if packetMessage == "STOP" {
			fmt.Println("Client connection has exited")
			count--
			break
		}
		clientConnection.Write([]byte(strings.Trim(strconv.QuoteToASCII(clientResponse), "\"")))
	}
	clientConnection.Close()
}
func packetDissect(netData string) (string, string) {
	data := strings.TrimSpace(string(netData))
	packetCode := data[0:2]
	packetMessage := strings.Replace(data, packetCode, "", -1)
	return packetCode, packetMessage
}
func processLoginPacket(packetMessage string) (string, string) {
	username := strings.Split(packetMessage, "?")[0]
	password := strings.Split(packetMessage, "?")[1]
	return username, password
}
func processRegisterPacket(packetMessage string) (string, string) {
	username := strings.Split(packetMessage, "?")[0]
	password := strings.Split(packetMessage, "?")[1]
	return username, password
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

// func addInventoryItem(userID string, itemID string, mongoClient *mongo.Client) string {
// 	var equipmentMap = map[string]int{"Head":0,"Body":1,"Feet":2,"Weapon":3,"Accessory":4}
// 	var itemMap = map[string]int{"Consumable":0,"Crafting":1,"Quest":2}
// 	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
// 		panic(err)
// 	}
// 	database := mongoClient.Database("player")
// 	inventory := database.Collection("users")
// 	item,itemFound := getItem(itemID, mongoClient)
// 	query := bson.M{"user_id":userID}

// 	itemType := bson.M{"type":item["type"]}
// 	if index, existing := equipmentMap[item["type"]]; existing{
// 		itemQuery := bson.M{"inventory":bson.M{"equipment"}}
// 	} else if index, existing:= itemMap[item["type"]]; existing{
// 		itemQuery := bson.M{"items"}
// 	}

// 	change := bson.M{"$push":bson.M{"inventory":itemArea}}
// 	filterCursor, err := inventory.UpdateOne(cxt, bson.M{"user_id": userID},bson.M{$push})
// 	return "Item added successfully!"
// }

func inventoryTest(userID string, mongoClient *mongo.Client) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	inventory := database.Collection("inventory")
	filterCursor, err := inventory.Find(cxt, bson.M{"user_id": userID})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var filterResult []PlayerInventory
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}

	fmt.Println(filterResult[0].Purse.Bits)
}
func getItem(itemID string, mongoClient *mongo.Client) (Item, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	items := database.Collection("items")
	filterCursor, err := items.Find(cxt, bson.M{"item_id": itemID})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var itemResult []Item
	if err = filterCursor.All(cxt, &itemResult); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Item retrieved : ", itemResult[0].Stats.Attack)
	return itemResult[0], false
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
	//getItem("WizardHat", mongoClient)
	inventoryTest("testUser", mongoClient)
	PORT := ":" + arguments[1]

	// add to sync.WaitGroup
	wg.Add(1)
	go tcpListener(PORT, cxt, mongoClient)
	wg.Add(1)
	go udpListener(PORT)
	wg.Wait()

}
