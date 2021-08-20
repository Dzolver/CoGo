package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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

type User struct {
	ObjectID primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	User_id  string             `json:"user_id" default:"" bson:"user_id, omitempty"`
	Password string             `json:"password" default:"" bson:"password, omitempty"`
	Logins   int                `json:"logins" default:"" bson:"logins, omitempty"`
}
type PlayerInventory struct {
	ObjectID  primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	User_id   string             `json:"user_id" default:"" bson:"user_id, omitempty"`
	Purse     Purse              `json:"purse" bson:"purse, omitempty"`
	Equipment Equipment          `json:"equipment" bson:"equipment, omitempty"`
	Items     Items              `json:"items" bson:"items, omitempty"`
}

type Purse struct {
	Bits float32 `default:"0" json:"bits" bson:"bits, omitempty"`
}
type Item struct {
	ObjectID    primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	Item_id     string             `json:"item_id" bson:"item_id, omitempty"`
	Item_type   string             `json:"item_type" bson:"item_type, omitempty"`
	Entity      string             `json:"entity" bson:"entity, omitempty"`
	Name        string             `json:"name" bson:"name, omitempty"`
	Num         int32              `json:"num" bson:"num, omitempty"`
	Description string             `json:"description" bson:"description, omitempty"`
	Stats       Stats              `json:"stats" bson:"stats, omitempty"`
}
type Stats struct {
	Health       float32 `json:"health" bson:"health"`
	Mana         float32 `json:"mana" bson:"mana"`
	Attack       float32 `json:"attack" bson:"attack"`
	MagicAttack  float32 `json:"magicAttack" bson:"magicAttack"`
	Defense      float32 `json:"defense" bson:"defense"`
	MagicDefense float32 `json:"magicDefense" bson:"magicDefense"`
	Armor        float32 `json:"armor" bson:"armor"`
	Evasion      float32 `json:"evasion" bson:"evasion"`
	Accuracy     float32 `json:"accuracy" bson:"accuracy"`
	Agility      float32 `json:"agility" bson:"agility"`
	Willpower    float32 `json:"willpower" bson:"willpower"`
	FireRes      float32 `json:"fireRes" bson:"fireRes"`
	WaterRes     float32 `json:"waterRes" bson:"waterRes"`
	EarthRes     float32 `json:"earthRes" bson:"earthRes"`
	WindRes      float32 `json:"windRes" bson:"windRes"`
	IceRes       float32 `json:"iceRes" bson:"iceRes"`
	EnergyRes    float32 `json:"energyRes" bson:"energyRes"`
	NatureRes    float32 `json:"natureRes" bson:"natureRes"`
	PoisonRes    float32 `json:"poisonRes" bson:"poisonRes"`
	MetalRes     float32 `json:"metalRes" bson:"metalRes"`
	LightRes     float32 `json:"lightRes" bson:"lightRes"`
	DarkRes      float32 `json:"darkRes" bson:"darkRes"`
}
type ItemRange struct {
	Collection []Item `json:"collection" bson:"collection"`
}
type Equipment struct {
	Head      ItemRange `json:"head" bson:"Head, omitempty"`
	Body      ItemRange `json:"body" bson:"Body, omitempty"`
	Feet      ItemRange `json:"feet" bson:"Feet, omitempty"`
	Weapon    ItemRange `json:"weapon" bson:"Weapon, omitempty"`
	Accessory ItemRange `json:"accessory" bson:"Accessory, omitempty"`
}
type Items struct {
	Consumable ItemRange `json:"consumable" bson:"Consumable, omitempty"`
	Crafting   ItemRange `json:"crafting" bson:"Crafting, omitempty"`
	Quest      ItemRange `json:"quest" bson:"Quests, omitempty"`
}
type PlayerSpellIndex struct {
	ObjectID    primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	User_id     string             `json:"user_id" bson:"user_id, omitempty"`
	Spell_index []Spell            `json:"spell_index" bson:"spell_index, omitempty"`
}
type Spell struct {
	ObjectID    primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	Spell_id    string             `json:"spell_id" bson:"spell_id, omitempty"`
	Name        string             `json:"name" bson:"name, omitempty"`
	Mana_cost   int32              `json:"mana_cost" bson:"mana_cost, omitempty"`
	Spell_type  string             `json:"spell_type" bson:"spell_type, omitempty"`
	Targetable  string             `json:"targetable" bson:"targetable, omitempty"`
	Spell       string             `json:"spell" bson:"spell, omitempty"`
	Damage      int32              `json:"damage" bson:"damage, omitempty"`
	Element     string             `json:"element" bson:"element, omitempty"`
	Level       int32              `json:"level" bson:"level, omitempty"`
	Init_block  int32              `json:"init_block" bson:"init_block, omitempty"`
	Block_count int32              `json:"block_count" bson:"block_count, omitempty"`
	Effect      Effect             `json:"effect" bson:"effect, omitempty"`
}
type Effect struct {
	Name             string `json:"name" bson:"name, omitempty"`
	Effect_id        string `json:"effect_id" bson:"effect_id, omitempty"`
	Element          string `json:"element" bson:"element, omitempty"`
	Effect_type      string `json:"effect_type" bson:"effect_type, omitempty"`
	Buff_element     string `json:"buff_element" bson:"buff_element, omitempty"`
	Debuff_element   string `json:"debuff_element" bson:"debuff_element, omitempty"`
	Damage_per_cycle int32  `json:"damage_per_cycle" bson:"damage_per_cycle, omitempty"`
	Lifetime         int32  `json:"lifetime" bson:"lifetime, omitempty"`
	Ticks_left       int32  `json:"ticks_left" bson:"ticks_left, omitempty"`
	Scalar           int32  `json:"scalar" bson:"scalar, omitempty"`
	Description      string `json:"description" bson:"description, omitempty"`
	Effector         string `json:"effector" bson:"effector, omitempty"`
}

var count = 0
var processLimitChan = make(chan int, 1000)
var wg sync.WaitGroup
var connectedUsers = make(map[string]bool)
var portNumbers = make(map[string]int)
var portNumbersReversed = make(map[string]string)
var portIndex = 1

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
		//Login
		if packetCode == "L0#" {
			clientResponse = packetCode
			fmt.Println("Login packet received!")
			username, password := processLoginPacket(packetMessage)
			loginResponse, valid := handleLogin(username, password, mongoClient)
			if valid {
				//Login success
				clientResponse = "LS#"
			} else if !valid {
				//Login fail
				clientResponse = "LF#"
			}
			//get a good port number for the udplistener and for the client to connect to
			wg.Add(1)
			designatedPortNumber := getPortFromIndex(portIndex)
			go udpListener(designatedPortNumber)
			clientResponse += "?" + loginResponse + "?" + designatedPortNumber
			portIndex--
			clientConnection.Write([]byte(strings.Trim(strconv.QuoteToASCII(clientResponse), "\"")))
		}
		//Register
		if packetCode == "R0#" {
			clientResponse = packetCode
			fmt.Println("Register packet received!")
			username, password := processRegisterPacket(packetMessage)
			registerResponse, valid := handleRegistration(username, password, mongoClient)
			if valid {
				//Register success
				clientResponse = "RS#"
			} else if !valid {
				//Register fail
				clientResponse = "RF#"
			}
			clientResponse += "?" + registerResponse
			clientConnection.Write([]byte(strings.Trim(strconv.QuoteToASCII(clientResponse), "\"")))
		}
		//Inventory add
		if packetCode == "IA#" {
			clientResponse = packetCode
			fmt.Println("Add Inventory packet received!")
			userID, itemID := processInventoryPacket(packetMessage)
			addInventoryItem(userID, itemID, mongoClient)
			clientConnection.Write([]byte(strings.Trim(strconv.QuoteToASCII(clientResponse), "\"")))
		}
		//Inventory delete
		if packetCode == "ID#" {
			clientResponse = packetCode
			fmt.Println("Delete Inventory packet received!")
		}
		//Inventory update
		if packetCode == "IU#" {
			clientResponse = packetCode
			fmt.Println("Update Inventory packet received!")
		}
		//Inventory create
		if packetCode == "IC#" {
			clientResponse = packetCode
			fmt.Println("Create Inventory packet received!")
			userID := strings.Split(packetMessage, "?")[0]
			success := createInventory(userID, mongoClient)
			if success {
				//Inventory success
				clientResponse = "IS#"
				clientResponse += "?Message from server : Inventory created succcesfully"
			}
			clientConnection.Write([]byte(strings.Trim(strconv.QuoteToASCII(clientResponse), "\"")))
		}
		if packetCode == "IR#" {
			clientResponse = packetCode
			fmt.Println("Read Inventory packet received!")
			userID := strings.Split(packetMessage, "?")[0]
			inventory, _ := getInventory(userID, mongoClient)
			inventoryJSON, _ := json.Marshal(inventory)
			clientResponse += "?" + string(inventoryJSON)
			clientConnection.Write([]byte(strings.Trim(strconv.QuoteToASCII(clientResponse), "\"")))
		}
		fmt.Println("Sent message back to client : ", clientResponse)
		if packetMessage == "STOP" {
			fmt.Println("Client connection has exited")
			count--
			break
		}

	}
	clientConnection.Close()
}
func packetDissect(netData string) (string, string) {
	data := strings.TrimSpace(string(netData))
	packetCode := data[0:3]
	packetMessage := strings.Replace(data, packetCode, "", -1)
	return packetCode, packetMessage
}
func processLoginPacket(packetMessage string) (string, string) {
	username := strings.Split(packetMessage, "?")[0]
	password := strings.Split(packetMessage, "?")[1]
	fmt.Println(username + " : " + password)
	return username, password
}
func processRegisterPacket(packetMessage string) (string, string) {
	username := strings.Split(packetMessage, "?")[0]
	password := strings.Split(packetMessage, "?")[1]
	return username, password
}
func processInventoryPacket(packetMessage string) (string, string) {
	username := strings.Split(packetMessage, "?")[0]
	itemID := strings.Split(packetMessage, "?")[1]
	return username, itemID
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

func handleLogin(username string, password string, mongoClient *mongo.Client) (string, bool) {
	if !lookForUser(username, mongoClient) {
		if validateUser(username, password, mongoClient) {
			fmt.Println("Login successful!")
			return "Login successful;" + username + ";" + strconv.Itoa(getLoginInfo(username, mongoClient)), true
		} else {
			fmt.Println("Login failed!")
			return "Login Failed;" + username + ";0", false
		}
	} else {
		if validateUser(username, password, mongoClient) {
			fmt.Println("Login successful!")
			return "Login Successful;" + username + ";" + strconv.Itoa(getLoginInfo(username, mongoClient)), true
		} else {
			fmt.Println("Login failed!")
			return "Login failed;" + username + ";0", false
		}
	}
}
func handleRegistration(username string, password string, mongoClient *mongo.Client) (string, bool) {
	if !lookForUser(username, mongoClient) {
		createUser(username, password, mongoClient)
		return "Account created", true
	} else {
		fmt.Println(username, " is not available")
		return "Username is not available", false
	}
}
func createUser(username string, password string, mongoClient *mongo.Client) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database := mongoClient.Database("player")
	users := database.Collection("users")
	createResult, err := users.InsertOne(cxt, bson.D{
		{"user_id", username},
		{"password", password},
		{"logins", 0},
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
	filterCursor, err := users.Find(cxt, bson.M{"user_id": username})
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
	filterCursor, err := users.Find(cxt, bson.M{"user_id": username, "password": password})
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
		_, err := users.UpdateOne(
			cxt,
			bson.M{"user_id": username},
			bson.D{
				{"$inc", bson.D{{"logins", 1}}},
			}, options.Update().SetUpsert(true))
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error with incrementing user login amount!")
			return false
		}
		return true

	}
	return true
}
func addInventoryItem(userID string, itemID string, mongoClient *mongo.Client) string {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	inventory := database.Collection("inventory")
	retrievedItem, itemFound := getItem(itemID, mongoClient)
	if itemFound {
		entryArea := "equipment." + retrievedItem.Item_type + ".collection"
		match := bson.M{"user_id": userID}
		change := bson.M{"$push": bson.M{entryArea: retrievedItem}}
		updateResponse, err := inventory.UpdateOne(cxt, match, change)
		fmt.Println(updateResponse)
		if err != nil {
			fmt.Println(err)
			return "Item addition failed!"
		}
		return "Item added successfully!"
	}
	return "Item does not exist!"
}
func addSpell(userID string, spellID string, mongoClient *mongo.Client) string {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	inventory := database.Collection("spellIndex")
	retrievedSpell, spellFound := getSpell(spellID, mongoClient)
	if spellFound {
		entryArea := "spell_index"
		match := bson.M{"user_id": userID}
		change := bson.M{"$push": bson.M{entryArea: retrievedSpell}}
		updateResponse, err := inventory.UpdateOne(cxt, match, change)
		fmt.Println(updateResponse)
		if err != nil {
			fmt.Println(err)
			return "Spell addition failed!"
		}
		return "Spell added successfully!"
	}
	return "Spell does not exist!"
}
func getLoginInfo(userID string, mongoClient *mongo.Client) int {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	users := database.Collection("users")
	filterCursor, err := users.Find(cxt, bson.M{"user_id": userID})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var filterResult []User
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	if len(filterResult) == 1 {
		return filterResult[0].Logins
	}
	return 0
}
func getInventory(userID string, mongoClient *mongo.Client) (*PlayerInventory, bool) {
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
	if len(filterResult) == 1 {
		return &filterResult[0], true
	}
	return nil, true
}
func createInventory(userID string, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	inventory := database.Collection("inventory")

	var freshInventory PlayerInventory
	freshInventory.User_id = userID

	freshInventory.Equipment.Head.Collection = make([]Item, 0)
	freshInventory.Equipment.Body.Collection = make([]Item, 0)
	freshInventory.Equipment.Feet.Collection = make([]Item, 0)
	freshInventory.Equipment.Weapon.Collection = make([]Item, 0)
	freshInventory.Equipment.Accessory.Collection = make([]Item, 0)

	freshInventory.Items.Consumable.Collection = make([]Item, 0)
	freshInventory.Items.Crafting.Collection = make([]Item, 0)
	freshInventory.Items.Quest.Collection = make([]Item, 0)

	insertResult, err := inventory.InsertOne(cxt, freshInventory)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Fresh Inventory created for user: ", userID, " insertID: ", insertResult.InsertedID)
	return true
}
func getItem(itemID string, mongoClient *mongo.Client) (*Item, bool) {
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
	if len(itemResult) == 0 {
		fmt.Println("Item not found!")
		emptyItem := new(Item)
		return emptyItem, false
	}
	fmt.Println("Item retrieved : ", itemResult[0].Stats.Attack)
	retrievedItem := itemResult[0]
	return &retrievedItem, true
}
func getSpell(spellID string, mongoClient *mongo.Client) (*Spell, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	spells := database.Collection("spells")
	filterCursor, err := spells.Find(cxt, bson.M{"spell_id": spellID})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var spellResult []Spell
	if err = filterCursor.All(cxt, &spellResult); err != nil {
		log.Fatal(err)
	}
	if len(spellResult) == 0 {
		fmt.Println("Item not found!")
		emptyItem := new(Spell)
		return emptyItem, false
	}
	retrievedItem := spellResult[0]
	return &retrievedItem, true
}
func IsNewUser(userID string) bool {
	var isNewUser bool
	if connectedUsers[userID] {
		isNewUser = false
	} else {
		isNewUser = true
	}
	return isNewUser
}
func createPorts() {
	min := 10000
	max := 20000
	portCandidate := strconv.Itoa(rand.Intn(max-min) + min)
	if portNumbers[portCandidate] != 0 {
		createPorts()
	} else {
		portNumbers[portCandidate] = portIndex
		portNumbersReversed[strconv.Itoa(portIndex)] = portCandidate
		if len(portNumbers) == 2000 {
			return
		}
		portIndex++
	}
}
func getPortFromIndex(index int) string {
	return portNumbersReversed[strconv.Itoa(index)]
}
func getIndexFromPort(port string) string {
	return strconv.Itoa(portNumbers[port])
}
func main() {
	createPorts()
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
	// addInventoryItem("asd", "WizardHat", mongoClient)
	// addInventoryItem("asd", "WizardHat", mongoClient)
	//disconnect mongoDB client on return
	defer func() {
		if err = mongoClient.Disconnect(cxt); err != nil {
			panic(err)
		}
	}()
	PORT := ":" + arguments[1]

	// add to sync.WaitGroup
	wg.Add(1)
	go tcpListener(PORT, cxt, mongoClient)
	wg.Wait()

}
