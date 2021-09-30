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
type PlayerLoadout struct {
	ObjectID    primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	User_id     string             `json:"user_id" default:"" bson:"user_id, omitempty"`
	Head        Item               `json:"head" bson:"Head, omitempty"`
	Body        Item               `json:"body" bson:"Body, omitempty"`
	Feet        Item               `json:"feet" bson:"Feet, omitempty"`
	Weapon      Item               `json:"weapon" bson:"Weapon, omitempty"`
	Accessory_1 Item               `json:"accessory_1" bson:"Accessory_1, omitempty"`
	Accessory_2 Item               `json:"accessory_2" bson:"Accessory_2, omitempty"`
	Accessory_3 Item               `json:"accessory_3" bson:"Accessory_3, omitempty"`
}

type Purse struct {
	Bits float32 `default:"0" json:"bits" bson:"bits, omitempty"`
}
type Item struct {
	ObjectID    primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	Item_id     string             `json:"item_id" default:"" bson:"item_id, omitempty"`
	Item_type   string             `json:"item_type" default:"" bson:"item_type, omitempty"`
	Entity      string             `json:"entity" bson:"entity, omitempty"`
	Name        string             `json:"name" default:"" bson:"name, omitempty"`
	Num         int32              `json:"num" default:"" bson:"num, omitempty"`
	Description string             `json:"description" default:"" bson:"description, omitempty"`
	Stats       Stats              `json:"stats" bson:"stats, omitempty"`
}
type Stats struct {
	Health       float64 `json:"health" default:"0" bson:"health"`
	Mana         float64 `json:"mana" default:"0" bson:"mana"`
	Attack       float64 `json:"attack" default:"0" bson:"attack"`
	MagicAttack  float64 `json:"magicAttack" default:"0" bson:"magicAttack"`
	Defense      float64 `json:"defense" default:"0" bson:"defense"`
	MagicDefense float64 `json:"magicDefense" default:"0" bson:"magicDefense"`
	Armor        float64 `json:"armor" default:"0" bson:"armor"`
	Evasion      float64 `json:"evasion" default:"0" bson:"evasion"`
	Accuracy     float64 `json:"accuracy" default:"0" bson:"accuracy"`
	Agility      float64 `json:"agility" default:"0" bson:"agility"`
	Willpower    float64 `json:"willpower" default:"0" bson:"willpower"`
	FireRes      float64 `json:"fireRes" default:"0" bson:"fireRes"`
	WaterRes     float64 `json:"waterRes" default:"0" bson:"waterRes"`
	EarthRes     float64 `json:"earthRes" default:"0" bson:"earthRes"`
	WindRes      float64 `json:"windRes" default:"0" bson:"windRes"`
	IceRes       float64 `json:"iceRes" default:"0" bson:"iceRes"`
	EnergyRes    float64 `json:"energyRes" default:"0" bson:"energyRes"`
	NatureRes    float64 `json:"natureRes" default:"0" bson:"natureRes"`
	PoisonRes    float64 `json:"poisonRes" default:"0" bson:"poisonRes"`
	MetalRes     float64 `json:"metalRes" default:"0" bson:"metalRes"`
	LightRes     float64 `json:"lightRes" default:"0" bson:"lightRes"`
	DarkRes      float64 `json:"darkRes" default:"0" bson:"darkRes"`
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
type PlayerProfile struct {
	Name        string  `json:"name" default:"" bson:"name, omitempty"`
	Level       int     `json:"level" default:"" bson:"level, omitempty"`
	Age         int     `json:"age" default:"" bson:"age, omitempty"`
	Title       string  `json:"title" default:"" bson:"title,omitempty"`
	Current_EXP float64 `json:"current_exp" default:"" bson:"current_exp, omitempty"`
	Total_EXP   float64 `json:"total_exp" default:"" bson:"total_exp, omitempty"`
	Max_EXP     float64 `json:"max_exp" default:"" bson:"max_exp, omitempty"`
	Race_id     string  `json:"race_id" default:"" bson:"race_id, omitempty"`
	Race_name   string  `json:"race_name" default:"" bson:"race_name, omitempty"`
	Class_id    string  `json:"class_id" default:"" bson:"class_id, omitempty"`
	Class_name  string  `json:"class_name" default:"" bson:"class_name, omitempty"`
}
type PlayerVital struct {
	ObjectID      primitive.ObjectID `json:"objectID" bson:"_id,omitempty"`
	User_id       string             `json:"user_id" default:"" bson:"user_id,omitempty"`
	PlayerProfile PlayerProfile      `json:"profile" default:"" bson:"profile, omitempty"`
	Stats         Stats              `json:"stats" default:"" bson:"stats, omitempty"`
	BaseStats     Stats              `json:"base_stats" default:"" bson:"base_stats, omitempty"`
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
	byteLimiter := 1024
	for {
		netData, err := bufio.NewReader(clientConnection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		packetCode, packetMessage := packetDissect(netData)
		fmt.Println(netData)
		//Login
		if packetCode == "L0#" {
			clientResponse = packetCode
			fmt.Println("Login packet received!")
			username, password := processTier2Packet(packetMessage)
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
			writeResponse(clientResponse, clientConnection)
		}
		//Register
		if packetCode == "R0#" {
			clientResponse = packetCode
			fmt.Println("Register packet received!")
			username, password := processTier2Packet(packetMessage)
			registerResponse, valid := handleRegistration(username, password, mongoClient)
			if valid {
				//Register success
				clientResponse = "RS#"
			} else if !valid {
				//Register fail
				clientResponse = "RF#"
			}
			clientResponse += "?" + registerResponse
			writeResponse(clientResponse, clientConnection)
		}
		//Inventory add
		if packetCode == "IA#" {
			clientResponse = packetCode
			fmt.Println("Add Inventory packet received!")
			fmt.Println("Packet message : ", packetMessage)
			userID, itemID := processTier2Packet(packetMessage)
			clientResponse += addInventoryItem(userID, itemID, mongoClient)
			writeResponse(clientResponse, clientConnection)
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
			fmt.Println("Packet message : ", packetMessage)
			username, itemID := processTier2Packet(packetMessage)

			clientResponse += "?" + addInventoryItem(username, itemID, mongoClient)
			writeResponse(clientResponse, clientConnection)
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
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "CE#" {
			clientResponse = packetCode
			fmt.Println("Create Everything packet received!")
			userID := strings.Split(packetMessage, "?")[0]
			createInventory(userID, mongoClient)
			createLoadout(userID, mongoClient)
			createSpellIndex(userID, mongoClient)
			createVital(userID, mongoClient)
			clientResponse += "?Everything created successfully"
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "IR#" {
			fmt.Println("Read Inventory packet received!")
			userID := strings.Split(packetMessage, "?")[0]
			inventory, _ := getInventory(userID, mongoClient)
			inventoryJSON, _ := json.Marshal(inventory)
			inventoryData := "?" + string(inventoryJSON)
			chainWriteResponse(packetCode, inventoryData, byteLimiter, clientConnection, "INVENTORY")
		}
		if packetCode == "LC#" {
			clientResponse = packetCode
			fmt.Println("Create loadout packet received!")
			userID := processTier1Packet(packetMessage)
			success := createLoadout(userID, mongoClient)
			clientResponse += "?" + strconv.FormatBool(success)
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "LR#" {
			clientResponse = packetCode
			fmt.Println("Read loadout packet received!")
			userID := processTier1Packet(packetMessage)
			loadout, _ := getLoadout(userID, mongoClient)
			loadoutJSON, _ := json.Marshal(loadout)
			loadoutData := "?" + string(loadoutJSON)
			chainWriteResponse(packetCode, loadoutData, byteLimiter, clientConnection, "LOADOUT")
		}
		if packetCode == "LE#" {
			clientResponse = packetCode
			fmt.Println("Equip to loadout")
			userID, itemID := processTier2Packet(packetMessage)
			equipFeedback := equipItem(userID, itemID, mongoClient)
			clientResponse += equipFeedback
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "LUE#" {
			clientResponse = packetCode
			fmt.Println("Unequip from loadout")
			userID, itemID := processTier2Packet(packetMessage)
			unequipFeedback := unequipItem(userID, itemID, mongoClient)
			clientResponse += unequipFeedback
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "SC#" {
			clientResponse = packetCode
			fmt.Println("Spell index creation packet received!")
			userID := processTier1Packet(packetMessage)
			success := createSpellIndex(userID, mongoClient)
			clientResponse += "?" + strconv.FormatBool(success)
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "SR#" {
			fmt.Println("Read Spell Index packet received!")
			userID := strings.Split(packetMessage, "?")[0]
			spellIndex, _ := getSpellIndex(userID, mongoClient)
			spellIndexJSON, _ := json.Marshal(spellIndex)
			spellIndexData := "?" + string(spellIndexJSON)
			chainWriteResponse(packetCode, spellIndexData, byteLimiter, clientConnection, "SPELLINDEX")
		}
		if packetCode == "SU#" {
			clientResponse = packetCode
			fmt.Println("Update Spell Index packet received!")
			userID, spellID := processTier2Packet(packetMessage)
			clientResponse += "?" + addSpell(userID, spellID, mongoClient)
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "VR#" {
			clientResponse = packetCode
			fmt.Println("Load vital packet received!")
			userID := processTier1Packet(packetMessage)
			vital, _ := getVital(userID, mongoClient)
			vitalJSON, _ := json.Marshal(vital)
			vitalData := "?" + string(vitalJSON)
			chainWriteResponse(packetCode, vitalData, byteLimiter, clientConnection, "PROFILE")
		}
		if packetCode == "LU#" {
			clientResponse = packetCode
			fmt.Println("Update EXP packet received!")
			userID, streamedEXPString := processTier2Packet(packetMessage)
			streamedEXP, _ := strconv.ParseFloat(streamedEXPString, 64)
			clientResponse += "?" + updateProfile_EXP(userID, streamedEXP, mongoClient)
			writeResponse(clientResponse, clientConnection)
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
	sep := strings.Index(data, "#")
	packetCode := data[0 : sep+1]
	packetMessage := strings.Replace(data, packetCode, "", -1)
	return packetCode, packetMessage
}
func processTier1Packet(packetMessage string) string {
	item := strings.Split(packetMessage, "?")[0]
	return item
}
func processTier2Packet(packetMessage string) (string, string) {
	item1 := strings.Split(packetMessage, "?")[0]
	item2 := strings.Split(packetMessage, "?")[1]
	return item1, item2
}
func writeResponse(clientResponse string, clientConnection net.Conn) {
	clientConnection.Write([]byte(strings.Trim(strconv.QuoteToASCII(clientResponse), "\"")))
}
func chainWriteResponse(packetCode string, totalData string, byteLimiter int, clientConnection net.Conn, serviceType string) {
	base := strings.Replace(packetCode, "#", "", -1)
	totalByteData := []byte(strings.Trim(strconv.QuoteToASCII(totalData), "\""))
	dataPartitions := len(totalByteData) / byteLimiter
	fullPartitions := 0
	remainingBytes := len(totalByteData) % byteLimiter
	if remainingBytes >= 1 {
		fullPartitions = dataPartitions
		dataPartitions++
	}
	for i := 1; i <= dataPartitions; i++ {
		start := 0
		end := 0
		currentPartition := strconv.Itoa(i)
		constructedPacketCode := base + currentPartition + strconv.Itoa(dataPartitions) + "#"
		//[0 : 1024] -> [1024 : 2048]
		start = (i - 1) * byteLimiter
		if i > fullPartitions {
			end = len(totalByteData)
		} else {
			end = (i) * byteLimiter
		}
		partitionedInventory := totalByteData[start:end]
		clientResponse := ""
		if i > 1 {
			clientResponse = constructedPacketCode + "?" + string(partitionedInventory)
		} else {
			clientResponse = constructedPacketCode + string(partitionedInventory)
		}
		fmt.Println("("+serviceType+") "+"Sent message back to client : ", clientResponse)
		clientConnection.Write([]byte(clientResponse))
	}
	fmt.Println("Size of ", serviceType, " data in bytes : ", len(totalByteData))
	fmt.Println("Size of remaining ", serviceType, " in bytes : ", len(totalByteData)%byteLimiter)
	fmt.Println("Size of ", serviceType, " partitions : ", dataPartitions)
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
func createVital(userID string, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profile := database.Collection("vital")

	var freshVital PlayerVital
	freshVital.ObjectID = primitive.NewObjectID()
	freshVital.User_id = userID
	freshVital.PlayerProfile.Name = userID
	freshVital.PlayerProfile.Level = 1
	freshVital.PlayerProfile.Age = 0
	freshVital.PlayerProfile.Title = "Just a newbie"
	freshVital.PlayerProfile.Current_EXP = 0
	freshVital.PlayerProfile.Total_EXP = 0
	freshVital.PlayerProfile.Max_EXP = 100
	freshVital.PlayerProfile.Race_id = "human"
	freshVital.PlayerProfile.Race_name = "Human"
	freshVital.PlayerProfile.Class_id = "stranger"
	freshVital.PlayerProfile.Class_name = "Stranger"
	freshVital.Stats.Health = 100
	freshVital.Stats.Mana = 100
	freshVital.Stats.Attack = 1
	freshVital.Stats.MagicAttack = 1
	freshVital.Stats.Defense = 1
	freshVital.Stats.MagicDefense = 1
	freshVital.Stats.Accuracy = 1
	freshVital.Stats.Agility = 1
	freshVital.Stats.Willpower = 1
	freshVital.BaseStats = freshVital.Stats

	insertResult, err := profile.InsertOne(cxt, freshVital)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Fresh Vital created for user: ", userID, " insertID: ", insertResult.InsertedID)
	return true
}
func getVital(userID string, mongoClient *mongo.Client) (*PlayerVital, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profile := database.Collection("vital")
	filterCursor, err := profile.Find(cxt, bson.M{"user_id": userID})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var filterResult []PlayerVital
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	if len(filterResult) == 1 {
		return &filterResult[0], true
	}
	return nil, true
}
func updateProfile_EXP(userID string, streamed_exp float64, mongoClient *mongo.Client) string {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	playerVital, profileFound := getVital(userID, mongoClient)
	newTotalExp := playerVital.PlayerProfile.Total_EXP + streamed_exp
	if profileFound {
		database := mongoClient.Database("player")
		profile := database.Collection("vital")
		match := bson.M{"user_id": userID}
		totalEXP := playerVital.PlayerProfile.Current_EXP + streamed_exp
		if totalEXP >= playerVital.PlayerProfile.Max_EXP {
			bufferEXP := 0.0
			levelUpperLimit := 0
			levelUpperLimitEXP := playerVital.PlayerProfile.Max_EXP
			for totalEXP > bufferEXP {
				levelUpperLimitEXP += float64(levelUpperLimit * 50.0)
				bufferEXP += playerVital.PlayerProfile.Max_EXP + float64(levelUpperLimit*50.0)
				levelUpperLimit++
			}
			newCurrentEXP := levelUpperLimitEXP - (bufferEXP - totalEXP)
			newLevel := playerVital.PlayerProfile.Level + levelUpperLimit
			newMaxEXP := levelUpperLimitEXP
			change := bson.D{
				{"$set", bson.D{{"PlayerProfile.level", newLevel}, {"PlayerProfile.current_exp", newCurrentEXP}, {"PlayerProfile.max_exp", newMaxEXP}, {"PlayerProfile.total_exp", newTotalExp}}},
			}
			_, err := profile.UpdateOne(cxt, match, change)
			if err != nil {
				fmt.Println(err)
				return "Profile level and EXP update failed!"
			}
			return "Profile level and EXP updated successfully!"
		} else if totalEXP < playerVital.PlayerProfile.Max_EXP {
			newCurrentEXP := totalEXP
			change := bson.D{
				{"$set", bson.D{{"PlayerProfile.current_exp", newCurrentEXP}}},
			}
			_, err := profile.UpdateOne(cxt, match, change)
			if err != nil {
				fmt.Println(err)
				return "Profile EXP update failed!"
			}
			return "Profile EXP updated successfully!"
		}
	}
	return "vital entry not found"
}
func createSpellIndex(userID string, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	spellIndex := database.Collection("spellIndex")

	var freshSpellIndex PlayerSpellIndex
	freshSpellIndex.ObjectID = primitive.NewObjectID()
	freshSpellIndex.User_id = userID
	freshSpellIndex.Spell_index = make([]Spell, 0)

	insertResult, err := spellIndex.InsertOne(cxt, freshSpellIndex)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Fresh Spell index created for user: ", userID, " insertID: ", insertResult.InsertedID)
	return true
}
func getSpellIndex(userID string, mongoClient *mongo.Client) (*PlayerSpellIndex, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	spellIndex := database.Collection("spellIndex")
	filterCursor, err := spellIndex.Find(cxt, bson.M{"user_id": userID})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var filterResult []PlayerSpellIndex
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	if len(filterResult) == 1 {
		return &filterResult[0], true
	}
	return nil, true
}
func addSpell(userID string, spellID string, mongoClient *mongo.Client) string {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	spellIndex := database.Collection("spellIndex")
	retrievedSpell, spellFound := getSpell(spellID, mongoClient)
	if spellFound {
		entryArea := "spell_index"
		match := bson.M{"user_id": userID}
		change := bson.M{"$push": bson.M{entryArea: retrievedSpell}}
		updateResponse, err := spellIndex.UpdateOne(cxt, match, change)
		fmt.Println(updateResponse)
		if err != nil {
			fmt.Println(err)
			return "Spell addition failed!"
		}
		return "Spell added successfully!"
	}
	return "Spell does not exist!"
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
	freshInventory.ObjectID = primitive.NewObjectID()
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
func createLoadout(userID string, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	inventory := database.Collection("loadout")

	var freshLoadout PlayerLoadout
	freshLoadout.ObjectID = primitive.NewObjectID()
	freshLoadout.User_id = userID

	freshLoadout.Head = Item{}
	freshLoadout.Body = Item{}
	freshLoadout.Feet = Item{}
	freshLoadout.Weapon = Item{}
	freshLoadout.Accessory_1 = Item{}
	freshLoadout.Accessory_2 = Item{}
	freshLoadout.Accessory_3 = Item{}

	insertResult, err := inventory.InsertOne(cxt, freshLoadout)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Fresh Loadout created for user: ", userID, " insertID: ", insertResult.InsertedID)
	return true
}
func getLoadout(userID string, mongoClient *mongo.Client) (*PlayerLoadout, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profile := database.Collection("loadout")
	filterCursor, err := profile.Find(cxt, bson.M{"user_id": userID})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var filterResult []PlayerLoadout
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	if len(filterResult) == 1 {
		return &filterResult[0], true
	}
	return nil, true
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
	fmt.Println("Item retrieved : ", itemResult[0].Item_id)
	retrievedItem := itemResult[0]
	return &retrievedItem, true
}
func equipItem(userID string, itemID string, mongoClient *mongo.Client) string {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	loadout := database.Collection("loadout")
	retrievedItem, itemFound := getItem(itemID, mongoClient)
	if itemFound {
		entryArea := retrievedItem.Item_type
		match := bson.M{"user_id": userID}
		change := bson.M{"$set": bson.M{entryArea: retrievedItem}}
		updateResponse, err := loadout.UpdateOne(cxt, match, change)
		fmt.Println(updateResponse)
		if err != nil {
			fmt.Println(err)
			return "EQUIP$0"
		}
		updateVital(userID, retrievedItem, "add", mongoClient)
		return "EQUIP$1"
	}
	return "EQUIP$0"
}
func unequipItem(userID string, itemID string, mongoClient *mongo.Client) string {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	loadout := database.Collection("loadout")
	retrievedItem, itemFound := getItem(itemID, mongoClient)
	if itemFound {
		entryArea := retrievedItem.Item_type
		match := bson.M{"user_id": userID}
		change := bson.M{"$set": bson.M{entryArea: Item{}}}
		updateResponse, err := loadout.UpdateOne(cxt, match, change)
		fmt.Println(updateResponse)
		if err != nil {
			fmt.Println(err)
			return "UNEQUIP$0"
		}
		updateVital(userID, retrievedItem, "remove", mongoClient)
		return "UNEQUIP$1"
	}
	return "UNEQUIP$0"
}
func updateVital(userID string, item *Item, operation string, mongoClient *mongo.Client) *PlayerVital {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	vital := database.Collection("vital")
	retrievedVital, _ := getVital(userID, mongoClient)
	originalStats := retrievedVital.Stats
	updatedStats := updateStatsByItem(&originalStats, item, operation)
	retrievedVital.Stats = *updatedStats
	match := bson.M{"user_id": userID}
	change := bson.M{"$set": bson.D{{"profile", retrievedVital.PlayerProfile}, {"stats", retrievedVital.Stats}}}
	updateResponse, err := vital.UpdateOne(cxt, match, change)
	fmt.Println(updateResponse)
	if err != nil {
		fmt.Println(err)
	}
	return retrievedVital
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
func updateStatsByItem(originalStats *Stats, item *Item, operation string) *Stats {
	op := 0.0
	if operation == "add" {
		op = 1.0
	} else if operation == "remove" {
		op = -1.0
	}
	originalStats.Health += (op * item.Stats.Health)
	originalStats.Mana += (op * item.Stats.Mana)
	originalStats.Attack += (op * item.Stats.Attack)
	originalStats.MagicAttack += (op * item.Stats.MagicAttack)
	originalStats.Defense += (op * item.Stats.Defense)
	originalStats.MagicDefense += (op * item.Stats.MagicDefense)
	originalStats.Armor += (op * item.Stats.Armor)
	originalStats.Evasion += (op * item.Stats.Evasion)
	originalStats.Accuracy += (op * item.Stats.Accuracy)
	originalStats.Agility += (op * item.Stats.Agility)
	originalStats.Willpower += (op * item.Stats.Willpower)
	originalStats.FireRes += (op * item.Stats.FireRes)
	originalStats.WaterRes += (op * item.Stats.WaterRes)
	originalStats.EarthRes += (op * item.Stats.EarthRes)
	originalStats.WindRes += (op * item.Stats.WindRes)
	originalStats.IceRes += (op * item.Stats.IceRes)
	originalStats.EnergyRes += (op * item.Stats.EnergyRes)
	originalStats.NatureRes += (op * item.Stats.NatureRes)
	originalStats.PoisonRes += (op * item.Stats.PoisonRes)
	originalStats.MetalRes += (op * item.Stats.MetalRes)
	originalStats.LightRes += (op * item.Stats.LightRes)
	originalStats.DarkRes += (op * item.Stats.DarkRes)

	return originalStats
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
