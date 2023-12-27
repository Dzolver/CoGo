package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	ObjectID   primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	Account_id uuid.UUID          `json:"uuid" bson:"uuid,omitempty"`
	User_id    string             `json:"user_id" default:"" bson:"user_id, omitempty"`
	Password   string             `json:"password" default:"" bson:"password, omitempty"`
	Active     int                `json:"active" default:"" bson:"active,omitempty"`
	Logins     int                `json:"logins" default:"" bson:"logins, omitempty"`
}
type Profile struct {
	ObjectID     primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	Account_id   uuid.UUID          `json:"uuid" bson:"uuid,omitempty"`
	Name         string             `json:"name" default:"" bson:"name, omitempty"`
	Level        int                `json:"level" default:"" bson:"level, omitempty"`
	Age          int                `json:"age" default:"" bson:"age, omitempty"`
	Title        string             `json:"title" default:"" bson:"title,omitempty"`
	Current_EXP  float64            `json:"current_exp" default:"" bson:"current_exp, omitempty"`
	Total_EXP    float64            `json:"total_exp" default:"" bson:"total_exp, omitempty"`
	Max_EXP      float64            `json:"max_exp" default:"" bson:"max_exp, omitempty"`
	Race_id      string             `json:"race_id" default:"" bson:"race_id, omitempty"`
	Race_name    string             `json:"race_name" default:"" bson:"race_name, omitempty"`
	Class_id     string             `json:"class_id" default:"" bson:"class_id, omitempty"`
	Class_name   string             `json:"class_name" default:"" bson:"class_name, omitempty"`
	LastPosition Position           `json:"last_position" bson:"last_position,omitempty"`
	LastRegion   string             `json:"last_region" default:"" bson:"last_region,omitempty"`
	LastLevel    string             `json:"last_level" default:"" bson:"last_level,omitempty"`
	Items        ItemRange          `json:"items" default:"" bson:"items,omitempty"`
	Purse        Purse              `json:"purse" default:"" bson:"purse,omitempty"`
	Loadout      Loadout            `json:"loadout" default:"" bson:"loadout,omitempty"`
	Stats        Stats              `json:"stats" default:"" bson:"stats,omitempty"`
	BaseStats    Stats              `json:"base_stats" default:"" bson:"base_stats,omitempty"`
	SpellIndex   []string           `json:"spell_index" default:"" bson:"spell_index, omitempty"`
	Description  string             `json:"description" default:"" bson:"description, omitempty"`
}
type Loadout struct {
	Head        string `json:"head" bson:"head, omitempty"`
	Body        string `json:"body" bson:"body, omitempty"`
	Feet        string `json:"feet" bson:"feet, omitempty"`
	Weapon      string `json:"weapon" bson:"weapon, omitempty"`
	Accessory_1 string `json:"accessory_1" bson:"accessory_1, omitempty"`
	Accessory_2 string `json:"accessory_2" bson:"accessory_2, omitempty"`
	Accessory_3 string `json:"accessory_3" bson:"accessory_3, omitempty"`
}
type Purse struct {
	Bits float64 `json:"bits" default:"0" bson:"bits, omitempty"`
}
type Item struct {
	ObjectID     primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	Item_id      string             `json:"item_id" default:"" bson:"item_id, omitempty"`
	Item_type    string             `json:"item_type" default:"" bson:"item_type, omitempty"`
	Item_subtype string             `json:"item_subtype" default:"" bson:"item_subtype, omitempty"`
	Entity       string             `json:"entity" bson:"entity, omitempty"`
	Name         string             `json:"name" default:"" bson:"name, omitempty"`
	Num          int32              `json:"num" default:"" bson:"num, omitempty"`
	Description  string             `json:"description" default:"" bson:"description, omitempty"`
	Stats        Stats              `json:"stats" bson:"stats, omitempty"`
	BaseValue    float64            `json:"base_value" bson:"base_value,omitempty"`
}
type Stats struct {
	Strength     float64 `json:"strength" default:"0" bson:"strength, omitempty"`
	Intelligence float64 `json:"intelligence" default:"0" bson:"intelligence, omitempty"`
	Dexterity    float64 `json:"dexterity" default:"0" bson:"dexterity, omitempty"`
	Charisma     float64 `json:"charisma" default:"0" bson:"charisma, omitempty"`
	Luck         float64 `json:"luck" default:"0" bson:"luck, omitempty"`
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
	Collection []string `json:"collection" bson:"collection"`
}
type ShopItem struct {
	Item_uuid uuid.UUID `json:"uuid" bson:"uuid,omitempty"`
	Item      Item      `json:"shop_item" bson:"shop_item"`
	Price     float64   `json:"price" bson:"price"`
}
type Spell struct {
	ObjectID       primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	Spell_id       string             `json:"spell_id" bson:"spell_id, omitempty"`
	Name           string             `json:"name" bson:"name, omitempty"`
	Mana_cost      int32              `json:"mana_cost" bson:"mana_cost, omitempty"`
	Spell_type     string             `json:"spell_type" bson:"spell_type, omitempty"`
	Targetable     string             `json:"targetable" bson:"targetable, omitempty"`
	Spell          string             `json:"spell" bson:"spell, omitempty"`
	Damage         int32              `json:"damage" bson:"damage, omitempty"`
	Element        string             `json:"element" bson:"element, omitempty"`
	Level          int32              `json:"level" bson:"level, omitempty"`
	Spell_duration int32              `json:"spell_duration" bson:"spell_duration, omitempty"`
	Init_block     int32              `json:"init_block" bson:"init_block, omitempty"`
	Block_count    int32              `json:"block_count" bson:"block_count, omitempty"`
	Effect         Effect             `json:"effect" bson:"effect, omitempty"`
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
type Client struct {
	Account_id       uuid.UUID    `json:"uuid" bson:"uuid, omitempty"`
	ConnectTime      time.Time    `json:"connect_time" bson:"connect_time, omitempty"`
	UDPAddress       *net.UDPAddr `json:"udp_addr" bson:"udp_addr, omitempty"`
	BroadcastAddress *net.UDPAddr `json:"broadcast_addr" bson:"broadcast_addr, omitempty"`
	Position         Position     `json:"position" bson:"position, omitempty"`
}
type Position struct {
	Position_x float64 `json:"pos_x" default:"0" bson:"pos_x, omitempty"`
	Position_y float64 `json:"pos_y" default:"1" bson:"pos_y, omitempty"`
	Position_z float64 `json:"pos_z" default:"0" bson:"pos_z, omitempty"`
}
type BattlePacket struct {
	BattleID        uuid.UUID  `json:"battle_id" default:"" bson:"battle_id"`
	PlayerProfile   *Profile   `json:"player_profile" default:"" bson:"player_profile, omitempty"`
	Monsters        *[]Monster `json:"monsters" default:"" bson:"monsters, omitempty"`
	MonsterQuantity int        `json:"monster_quantity" default:"0" bson:"monster_quantity, omitempty"`
}
type LoginSecretPacket struct {
	User          *User       `json:"user_data" default:"" bson:"user_data,omitempty"`
	PlayerProfile *Profile    `json:"player_profile" default:"" bson:"player_profile, omitempty"`
	Region        *RegionData `json:"region_data" default:"" bson:"region_data,omitempty"`
}
type Region struct {
	ObjectID   primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	RegionID   string             `json:"region_id" default:"" bson:"regionID,omitempty"`
	RegionName string             `json:"region_name" default:"" bson:"regionName,omitempty"`
	Levels     []string           `json:"levels" bson:"levels,omitempty"`
}
type Level struct {
	ObjectID  primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	LevelID   string             `json:"level_id" default:"" bson:"levelID,omitempty"`
	LevelName string             `json:"level_name" default:"" bson:"levelName,omitempty"`
	ZIP       string             `json:"zip" default:"" bson:"ZIP,omitempty"`
	Monsters  []string           `json:"monsters" bson:"monsters,omitempty"`
	Residents []string           `json:"residents" bson:"residents,omitempty"`
}
type Resident struct {
	ObjectID primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	NpcID    string             `json:"npc_id" default:"" bson:"npcID, omitempty"`
	NpcName  string             `json:"npc_name" default:"" bson:"npcName, omitempty"`
	Dialogue []string           `json:"dialogue" default:"" bson:"dialogue, omitempty"`
}
type ShopKeeper struct {
	ObjectID  primitive.ObjectID `json:"objectID" bson:"_id,omitempty"`
	NpcID     string             `json:"npc_id" default:"" bson:"npcID,omitempty"`
	Catalogue []ShopItem         `json:"catalogue" default:"" bson:"catalogue,omitempty"`
	Purse     Purse              `json:"purse" default:"" bson:"purse,omitempty"`
}
type Monster struct {
	ObjectID       primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	MobID          string             `json:"mob_id" default:"" bson:"mobID,omitempty"`
	MonsterType    string             `json:"monster_type" default:"" bson:"monsterType,omitempty"`
	GoldGain       int                `json:"gold_gain" default:"" bson:"goldGain,omitempty"`
	ExperienceGain int                `json:"experience_gain" default:"" bson:"experienceGain,omitempty"`
	Profile        *Profile           `json:"profile" default:"" bson:"mobVitals,omitempty"`
	Stats          *Stats             `json:"stats" default:"" bson:"stats,omitempty"`
	Actions        *[]Spell           `json:"actions" default:"" bson:"attackActions,omitempty"`
	Element        string             `json:"element" default:"" bson:"element, omitempty"`
	Regions        []string           `json:"regions" bson:"regions, omitempty"`
}
type RegionData struct {
	Region    *Region    `json:"region" default:"" bson:"region"`
	LevelData *LevelData `json:"level_data" default:"" bson:"level"`
}
type LevelData struct {
	Level     *Level      `json:"level" default:"" bson:"level"`
	Residents *[]Resident `json:"residents" default:"" bson:"residents"`
}
type Sessions struct {
	Battles map[uuid.UUID]BattleSession `json:"battle_sessions" default:"" bson:"battle_sessions"`
}
type BattleSession struct {
	BattleID     uuid.UUID  `json:"battle_id" default:"" bson:"battle_id"`
	Status       int        `json:"status" default:"0" bson:"status"`
	Monsters     *[]Monster `json:"monsters" default:"" bson:"monsters"`
	RewardMatrix []int      `json:"reward_matrix" default:"" bson:"reward_matrix"`
	Reward       Reward     `json:"reward" default:"" bson:"reward"`
}
type Reward struct {
	Gold     float64 `json:"gold" default:"0" bson:"gold"`
	Exp      float64 `json:"exp" default:"0" bson:"exp"`
	TotalExp float64 `json:"total_exp" default:"0" bson:"total_exp"`
}
type Packet struct {
	PacketID    uuid.UUID `json:"packet_id" default:""`
	PacketCode  string    `json:"packet_code" default:""`
	Chain       bool      `json:"chain" default:""`
	ServiceType string    `json:"service_type" default:""`
	Content     string    `json:"content" default:""`
}
type PlayerPacketCache struct {
	PacketCache map[uuid.UUID]Packet `json:"packet_cache" default:""`
}

var PACKET_SIZE = 10000
var wg sync.WaitGroup
var ALLshopkeepers = make(map[string]ShopKeeper)
var playerPacketCache = make(map[uuid.UUID]PlayerPacketCache)
var MASTER_ITEM_TABLE = make(map[string]Item)
var MASTER_SPELL_TABLE = make(map[string]Spell)
var MASTER_MONSTER_TABLE = make(map[string]Monster)
var MASTER_LEVEL_TABLE = make(map[string]Level)

var sessions Sessions
var (
	Info           = Teal
	IncomingPacket = Magenta
	Warn           = Yellow
	Fata           = Red
	Success        = Green
	Failure        = Yellow
	Internal       = Mint
)
var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
	Mint    = Color("\033[1;158m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func handleTCPConnection(clientConnection net.Conn, cxt context.Context, mongoClient *mongo.Client) {
	fmt.Print(".")
	clientResponse := "DEFAULT"
	byteLimiter := PACKET_SIZE
	for {
		netData, err := bufio.NewReader(clientConnection).ReadString('\n')
		if err != nil {
			fmt.Println(Failure(err))
			return
		}
		packetCode, packetMessage := packetDissect(netData)
		fmt.Println(packetMessage)
		if packetCode == "BR#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Read for Battle packet received!"))
			requestIDSTR, accountIDSTR, levelID := processTier3Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			var freshBattlePacket BattlePacket
			playerProfile, _ := getProfile(accountID, mongoClient)
			level := getLevel(levelID, mongoClient)
			monsters := getMonsters(level.Monsters, mongoClient)
			//create BattleSession out of this information and add to the list of sessions
			freshBattlePacket.MonsterQuantity = 1
			battle := createBattle(monsters, freshBattlePacket.MonsterQuantity)

			freshBattlePacket.BattleID = battle.BattleID
			freshBattlePacket.PlayerProfile = playerProfile
			freshBattlePacket.Monsters = monsters

			contentJSON, _ := json.Marshal(freshBattlePacket)
			packet := createMultiDeliveryPacket(requestIDSTR, packetCode, "BATTLE", contentJSON)
			chainWriteResponse(accountID, requestIDSTR, packet, byteLimiter, clientConnection, false)
		}
		if packetCode == "B0#" {
			fmt.Println(IncomingPacket("Battle State Confirmation packet received!"))
			// accountIDSTR := processTier1Packet(packetMessage)
			// accountID, _ = uuid.Parse(accountIDSTR)
		}
		if packetCode == "BF#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Battle Finish packet received!"))
			fmt.Println(Info(packetMessage))
			//e.g: battleID?BF#1,1,1 -> battleID, [1,1,1]
			requestIDSTR, accountIDSTR, battleIDSTR, rewardMatrixSTR := processTier4Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			battleID, _ := uuid.Parse(battleIDSTR)

			//find out how much gold and exp is earned from the reward matrix
			fmt.Println(Info(rewardMatrixSTR))
			rewardMatrix := getArrayFromString(rewardMatrixSTR)
			exp := 0.0
			gold := 0.0
			updateStatus := "False"
			if entry, ok := sessions.Battles[battleID]; ok {
				entry.RewardMatrix = rewardMatrix
				entry.Status = 1
				for index, reward := range entry.RewardMatrix {
					monster := (*entry.Monsters)[index]
					if reward == 1 {
						entry.Reward.Exp += float64(monster.ExperienceGain)
						entry.Reward.Gold += float64(monster.GoldGain)
					}
				}
				exp = entry.Reward.Exp
				gold = entry.Reward.Gold
				//calculate exp and return total exp to entry.Reward.Totalexp
				addBits(accountID, gold, true, mongoClient)
				entry.Reward.TotalExp = updateProfile_EXP(accountID, entry.Reward.Exp, mongoClient)
				sessions.Battles[battleID] = entry
				updateStatus = "True"
			}
			profile, _ := getProfile(accountID, mongoClient)
			profileJSON, _ := json.Marshal(profile)
			contentJSON := strconv.FormatFloat(exp, 'f', -1, 64) + "|" + strconv.FormatFloat(gold, 'f', -1, 64) + "|" + updateStatus + "|" + string(profileJSON)
			packet := createMultiDeliveryPacket(requestIDSTR, packetCode, "BATTLEFINISH", []byte(contentJSON))
			chainWriteResponse(accountID, requestIDSTR, packet, byteLimiter, clientConnection, false)
		}
		if packetCode == "HB#" {
			// fmt.Println("Heartbeat packet received!")
			accountID, x, y, z := processTier4Packet(packetMessage)
			target_uuid, _ := uuid.Parse(accountID)
			var lastPosition Position
			lastPosition.Position_x, _ = strconv.ParseFloat(x, 64)
			lastPosition.Position_y, _ = strconv.ParseFloat(y, 64)
			lastPosition.Position_z, _ = strconv.ParseFloat(z, 64)
			updateUserLastPosition(target_uuid, &lastPosition, mongoClient)
		}
		//Inventory add
		if packetCode == "IA#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Add Inventory packet received!"))
			requestIDSTR, accountIDSTR, itemID := processTier3Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			content := addInventoryItem(accountID, itemID, mongoClient)
			packet := createSimpleDeliveryPacket(requestIDSTR, packetCode, "INVENTORY", content)
			writeResponse(accountID, requestIDSTR, packet, clientConnection, false)
		}
		if packetCode == "ID#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Delete Inventory packet received!"))
		}
		//Inventory update
		if packetCode == "IU#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Update Inventory packet received!"))
			requestIDSTR, accountIDSTR, itemID := processTier3Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			content := addInventoryItem(accountID, itemID, mongoClient)
			packet := createSimpleDeliveryPacket(requestIDSTR, packetCode, "INVENTORY", content)
			writeResponse(accountID, requestIDSTR, packet, clientConnection, false)
		}
		//Login
		if packetCode == "L0#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Login packet received!"))
			fmt.Println(Info(packetMessage))
			requestIDSTR, username, password := processTier3Packet(packetMessage)
			loginResponse, valid := handleLogin(username, password, mongoClient)
			if valid {
				//Login success
				packetCode = "LS#"
				/*var LSP LoginSecretPacket
				User, _ := getUser(username, mongoClient)
				Profile, _ := getProfile(User.Account_id, mongoClient)
				LSP.User = User
				LSP.PlayerProfile = Profile

				var FRD *RegionData = new(RegionData)
				Region := getRegion(Profile.LastRegion, mongoClient)
				FRD.Region = Region

				var FLD *LevelData = new(LevelData)
				Level := getLevel(Profile.LastLevel, mongoClient)
				Residents := getNPCs(Level.Residents, mongoClient)
				FLD.Level = Level
				FLD.Residents = Residents
				FRD.LevelData = FLD
				LSP.Region = FRD

				contentJSON, _ := json.Marshal(LSP)
				packet := createMultiDeliveryPacket(requestIDSTR, packetCode, "LSP", contentJSON)
				chainWriteResponse(User.Account_id, requestIDSTR, packet, byteLimiter, clientConnection, false)*/
				contentJSON, _ := json.Marshal(getLevelFromCache("00001"))
				packet := createMultiDeliveryPacket(requestIDSTR, packetCode, "LSP", contentJSON)
				chainWriteResponse(uuid.New(), requestIDSTR, packet, byteLimiter, clientConnection, false)
			} else if !valid {
				//Login fail
				packetCode = "LF#"
				content := loginResponse
				packet := createSimpleDeliveryPacket(requestIDSTR, packetCode, "LOGIN", content)
				fmt.Println(Info(clientResponse))
				fakeID := uuid.New()
				writeResponse(fakeID, requestIDSTR, packet, clientConnection, true)
			}
			//get a good port number for the udplistener and for the client to connect to
			//wg.Add(1)
			//designatedPortNumber := getPortFromIndex(portIndex)
			//go udpListener(designatedPortNumber, cxt, mongoClient)
			//portIndex--
		}
		if packetCode == "LE#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Equip to loadout"))
			requestIDSTR, accountIDSTR, itemID := processTier3Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			equipFeedback := equipItem(accountID, itemID, mongoClient)
			packet := createSimpleDeliveryPacket(requestIDSTR, packetCode, "LOADOUT", equipFeedback)
			writeResponse(accountID, requestIDSTR, packet, clientConnection, false)
		}
		if packetCode == "LL#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Level packet received!"))
			requestIDSTR, accountIDSTR, levelID := processTier3Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			var freshLevel LevelData
			level := getLevel(levelID, mongoClient)
			NPC := getNPCs(level.Residents, mongoClient)
			freshLevel.Level = level
			freshLevel.Residents = NPC
			contentJSON, _ := json.Marshal(freshLevel)
			packet := createMultiDeliveryPacket(requestIDSTR, packetCode, ":EVEL", contentJSON)
			chainWriteResponse(accountID, requestIDSTR, packet, byteLimiter, clientConnection, false)
		}
		if packetCode == "PR#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Read loadout packet received!"))
			requestIDSTR, accountIDSTR := processTier2Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			profile, _ := getProfile(accountID, mongoClient)
			profileJSON, _ := json.Marshal(profile)
			packet := createMultiDeliveryPacket(requestIDSTR, packetCode, "PROFILE", profileJSON)
			chainWriteResponse(accountID, requestIDSTR, packet, byteLimiter, clientConnection, false)
		}
		if packetCode == "LU#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Update EXP packet received!"))
			requestIDSTR, accountIDSTR, streamedEXPString := processTier3Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			streamedEXP, _ := strconv.ParseFloat(streamedEXPString, 64)
			newTotalEXP := updateProfile_EXP(accountID, streamedEXP, mongoClient)
			content := strconv.FormatFloat(newTotalEXP, 'E', -1, 64)
			packet := createSimpleDeliveryPacket(requestIDSTR, packetCode, "EXP", content)
			writeResponse(accountID, requestIDSTR, packet, clientConnection, false)
		}
		if packetCode == "LUE#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Unequip from loadout"))
			requestIDSTR, accountIDSTR, itemID := processTier3Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			content := unequipItem(accountID, itemID, mongoClient)
			packet := createSimpleDeliveryPacket(requestIDSTR, packetCode, "LOADOUT", content)
			writeResponse(accountID, requestIDSTR, packet, clientConnection, false)
		}
		if packetCode == "OK#" {
			fmt.Println(IncomingPacket("OK Packet received!"))
			requestIDSTR, accountIDSTR := processTier2Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			requestID, _ := uuid.Parse(requestIDSTR)
			_, found := playerPacketCache[accountID].PacketCache[requestID]
			if found {
				delete(playerPacketCache[accountID].PacketCache, requestID)
				if len(playerPacketCache[accountID].PacketCache) == 0 {
					delete(playerPacketCache, accountID)
				}
			}
		}
		//Register
		if packetCode == "R0#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Register packet received!"))
			fmt.Println(Info(packetMessage))
			requestIDSTR, username, password := processTier3Packet(packetMessage)
			registerResponse, valid, accountID := handleRegistration(username, password, mongoClient)
			if valid {
				//Register success
				createProfile(username, accountID, mongoClient)
				addSpell(accountID, "Fireball", mongoClient)
				addSpell(accountID, "Scorch", mongoClient)
				addInventoryItem(accountID, "WizardRobe", mongoClient)
				addInventoryItem(accountID, "WizardHat", mongoClient)
				clientResponse = "RS#"
			} else if !valid {
				//Register fail
				clientResponse = "RF#"
			}
			packet := createSimpleDeliveryPacket(requestIDSTR, packetCode, "REGISTER", clientResponse+registerResponse)
			writeResponse(accountID, requestIDSTR, packet, clientConnection, true)
		}
		if packetCode == "RLL#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Region and Level packet received!"))
			fmt.Println(Info(packetMessage))
			requestIDSTR, accountIDSTR, regionID, levelID := processTier4Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			var FRD RegionData
			var FLD LevelData
			region := getRegion(regionID, mongoClient)
			level := getLevel(levelID, mongoClient)
			NPC := getNPCs(level.Residents, mongoClient)
			FRD.Region = region
			FLD.Level = level
			FLD.Residents = NPC
			FRD.LevelData = &FLD
			contentJSON, _ := json.Marshal(FRD)
			packet := createMultiDeliveryPacket(requestIDSTR, packetCode, "REGION", contentJSON)
			chainWriteResponse(accountID, requestIDSTR, packet, byteLimiter, clientConnection, false)
		}
		if packetCode == "SH#" {
			fmt.Println(IncomingPacket("Shopkeeper Request Packet received"))
			requestIDSTR, accountIDSTR, npcID := processTier3Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			if entry, found := ALLshopkeepers[npcID]; found {
				shopkeeper := entry
				contentJSON, _ := json.Marshal(shopkeeper)
				packet := createMultiDeliveryPacket(requestIDSTR, packetCode, "SHOPKEEPER", contentJSON)
				chainWriteResponse(accountID, requestIDSTR, packet, byteLimiter, clientConnection, false)
			}
		}
		if packetCode == "SOS#" {
			fmt.Println(IncomingPacket("SOS Packet received!"))
			requestIDSTR, accountIDSTR := processTier2Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			requestID, _ := uuid.Parse(requestIDSTR)
			SOSPacket, found := playerPacketCache[accountID].PacketCache[requestID]
			if found {
				if SOSPacket.Chain {
					chainWriteResponse(accountID, requestIDSTR, SOSPacket, byteLimiter, clientConnection, true)
				} else if !SOSPacket.Chain {
					writeResponse(accountID, requestIDSTR, SOSPacket, clientConnection, true)
				}
			} else {
				fmt.Println(Warn("SOS Packet ID : " + requestIDSTR + " is NOT Found!"))
			}
		}
		if packetCode == "SU#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Update Spell Index packet received!"))
			requestIDSTR, accountIDSTR, spellID := processTier3Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			content := addSpell(accountID, spellID, mongoClient)
			packet := createSimpleDeliveryPacket(requestIDSTR, packetCode, "SPELL", content)
			writeResponse(accountID, requestIDSTR, packet, clientConnection, false)
		}
		if packetCode == "TT#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("TEST MESSAGE received!"))
			testMessageSTR := processTier1Packet(packetMessage)
			responseMessage := "Server response to : " + testMessageSTR
			clientConnection.Write([]byte(strings.Trim(strconv.QuoteToASCII(responseMessage), "\"")))
		}
		if packetCode == "XX#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("CLIENT WANTS TO SAY HI!"))
			clientMessage := processTier1Packet(packetMessage)
			fmt.Println(Info(clientMessage))
			greetings := sayHiToClient()
			clientConnection.Write([]byte(strings.Trim(strconv.QuoteToASCII(greetings), "\"")))
		}
		if packetMessage == "STOP" {
			fmt.Println("Client connection has exited")
			break
		}
	}
	clientConnection.Close()
}

// func toProtoBuf(inputStruct interface{}) proto.Message {
// 	switch s := inputStruct.(type) {
// 	case User:
// 		return &models.User{
// 			ObjectID: s.ObjectID.String(),
// 			Uuid:     s.Account_id.String(),
// 			UserId:   s.User_id,
// 			Password: s.Password,
// 			Active:   int32(s.Active),
// 			Logins:   int32(s.Logins),
// 		}

// 	case Profile:
// 		var items ItemRange = s.Items

// 		return &models.Profile{
// 			ObjectID:   s.ObjectID.String(),
// 			Uuid:       s.Account_id.String(),
// 			Name:       s.Name,
// 			Level:      int32(s.Level),
// 			Age:        int32(s.Age),
// 			Title:      s.Title,
// 			CurrentExp: float32(s.Current_EXP),
// 			TotalExp:   float32(s.Total_EXP),
// 			MaxExp:     float32(s.Max_EXP),
// 			RaceId:     s.Race_id,
// 			RaceName:   s.Race_name,
// 			ClassId:    s.Class_id,
// 			ClassName:  s.Class_name,
// 			LastPosition: &models.Position{
// 				PosX: float32(s.LastPosition.Position_x),
// 				PosY: float32(s.LastPosition.Position_y),
// 				PosZ: float32(s.LastPosition.Position_z),
// 			},
// 			LastRegion:  s.LastRegion,
// 			LastLevel:   s.LastLevel,
// 			Items:       toProtoBuf(s.Items),
// 			Purse:       &models.Purse{},
// 			Loadout:     &models.Loadout{},
// 			Stats:       &models.Stats{},
// 			BaseStats:   &models.Stats{},
// 			SpellIndex:  []string{},
// 			Description: "",
// 		}

// 	case Loadout:
// 		return &models.Loadout{
// 			head:        s.Head,
// 			body:        s.Body,
// 			feet:        s.Feet,
// 			weapon:      s.Weapon,
// 			accessory_1: s.Accessory_1,
// 			accessory_2: s.Accessory_2,
// 			accessory_3: s.Accessory_3,
// 		}

// 	case Purse:
// 		return &models.Purse{
// 			bits: s.Bits,
// 		}

// 	case Item:
// 		return &models.Item{
// 			ObjectID:    s.ObjectID.String(),
// 			ItemId:      s.Item_id,
// 			ItemType:    s.Item_type,
// 			ItemSubtype: s.Item_subtype,
// 			Entity:      s.Entity,
// 			Name:        s.Name,
// 			Num:         s.Num,
// 			Description: s.Description,
// 			Stats:       toProtoBuf(s.Stats),
// 			BaseValue:   0,
// 		}

// 	case Stats:
// 		return &models.Stats{
// 			Strength:     float32(s.Strength),
// 			Intelligence: float32(s.Intelligence),
// 			Dexterity:    float32(s.Dexterity),
// 			Charisma:     float32(s.Charisma),
// 			Luck:         float32(s.Luck),
// 			Health:       float32(s.Health),
// 			Mana:         float32(s.Mana),
// 			Attack:       float32(s.Attack),
// 			MagicAttack:  float32(s.MagicAttack),
// 			Defense:      float32(s.Defense),
// 			MagicDefense: float32(s.MagicDefense),
// 			Armor:        float32(s.Armor),
// 			Evasion:      float32(s.Evasion),
// 			Accuracy:     float32(s.Accuracy),
// 			Agility:      float32(s.Agility),
// 			FireRes:      float32(s.FireRes),
// 			WaterRes:     float32(s.WaterRes),
// 			EarthRes:     float32(s.EarthRes),
// 			WindRes:      float32(s.WindRes),
// 			IceRes:       float32(s.IceRes),
// 			EnergyRes:    float32(s.EnergyRes),
// 			NatureRes:    float32(s.NatureRes),
// 			PoisonRes:    float32(s.PoisonRes),
// 			MetalRes:     float32(s.MetalRes),
// 			LightRes:     float32(s.LightRes),
// 			DarkRes:      float32(s.DarkRes),
// 		}

// 	case ItemRange:
// 		return &models.ItemRange{
// 			Collection: s.Collection,
// 		}

// 	case ShopItem:
// 		return &models.ShopItem{
// 			Uuid:     s.Item_uuid.String(),
// 			ShopItem: toProtoBuf(s.Item),
// 			Price:    s.Price,
// 		}

// 	case Spell:
// 		return &models.Spell{
// 			objectID:       s.ObjectID,
// 			spell_id:       s.Spell_id,
// 			name:           s.Name,
// 			mana_cost:      s.Mana_cost,
// 			spell_type:     s.Spell_type,
// 			targetable:     s.Targetable,
// 			spell:          s.Spell,
// 			damage:         s.Damage,
// 			element:        s.Element,
// 			level:          s.Level,
// 			spell_duration: s.Spell_duration,
// 			init_block:     s.Init_block,
// 			block_count:    s.Block_count,
// 			effect:         toProtoBuf(s.Effect),
// 		}

// 	case Effect:
// 		return &models.Effect{
// 			name:             s.Name,
// 			effect_id:        s.Effect_id,
// 			element:          s.Element,
// 			effect_type:      s.Effect_type,
// 			buff_element:     s.Buff_element,
// 			debuff_element:   s.Debuff_element,
// 			damage_per_cycle: s.Damage_per_cycle,
// 			lifetime:         s.Lifetime,
// 			ticks_left:       s.Ticks_left,
// 			scalar:           s.Scalar,
// 			description:      s.Description,
// 			effector:         s.Effector,
// 		}

// 	case Client:
// 		return &models.Client{
// 			uuid:           s.Account_id,
// 			connect_time:   s.ConnectTime,
// 			udp_addr:       s.UDPAddress,
// 			broadcast_addr: s.BroadcastAddress,
// 			position:       toProtoBuf(s.Position),
// 		}

// 	case Position:
// 		return &models.Position{
// 			pos_x: s.Position_x,
// 			pos_y: s.Position_y,
// 			pos_z: s.Position_z,
// 		}

// 	case BattlePacket:
// 		var protoMonsters []*models.Monster
// 		for _, monster := range *s.Monsters {
// 			protoMonster := toProtoBuf(monster)
// 			protoMonsters = append(protoMonsters, protoMonster)
// 		}
// 		return &models.BattlePacket{
// 			battle_id:        s.BattleID,
// 			player_profile:   toProtoBuf(s.PlayerProfile),
// 			monsters:         protoMonsters,
// 			monster_quantity: s.MonsterQuantity,
// 		}

// 	case LoginSecretPacket:
// 		return &models.LoginSecretPacket{
// 			user_data:      toProtoBuf(s.User),
// 			player_profile: toProtoBuf(s.Profile),
// 			region_data:    toProtoBuf(s.Region),
// 		}

// 	case Region:
// 		return &models.Region{
// 			objectID:    s.ObjectID,
// 			region_id:   s.RegionID,
// 			region_name: s.RegionName,
// 			levels:      s.Levels,
// 		}

// 	case Level:
// 		return &models.Level{
// 			objectID:   s.ObjectID,
// 			level_id:   s.LevelID,
// 			level_name: s.LevelName,
// 			zip:        s.ZIP,
// 			monsters:   s.Monsters,
// 			residents:  s.Residents,
// 		}

// 	case Resident:
// 		return &models.Resident{
// 			objectID: s.ObjectID,
// 			npc_id:   s.NpcID,
// 			NpcName:  s.NpcName,
// 			dialogue: s.Dialogue,
// 		}

// 	case ShopKeeper:
// 		var protoCatalogue []*models.ShopItem
// 		for _, item := range *&s.Catalogue {
// 			protoItem := toProtoBuf(item)
// 			protoCatalogue = append(protoCatalogue, protoItem)
// 		}
// 		return &models.ShopKeeper{
// 			objectID:  s.ObjectID,
// 			npc_id:    s.NpcID,
// 			catalogue: protoCatalogue,
// 			purse:     toProtoBuf(s.Purse),
// 		}

// 	case Monster:
// 		var protoActions []*models.Spell
// 		for _, action := range *s.Actions {
// 			protoAction := toProtoBuf(action)
// 			protoActions = append(protoActions, protoAction)
// 		}
// 		return &models.Monster{
// 			objectID:        s.ObjectID,
// 			mob_id:          s.MobID,
// 			monster_type:    s.MonsterType,
// 			gold_gain:       s.GoldGain,
// 			experience_gain: s.ExperienceGain,
// 			profile:         toProtoBuf(s.profile),
// 			stats:           toProtoBuf(s.Stats),
// 			actions:         protoActions,
// 			element:         s.Element,
// 			regions:         s.Regions,
// 		}

// 	case RegionData:
// 		return &models.RegionData{
// 			region:     toProtoBuf(s.Region),
// 			level_data: toProtoBuf(s.LevelData),
// 		}

// 	case LevelData:
// 		var protoResidents []*models.Resident
// 		for _, resident := range *s.Residents {
// 			protoResident := toProtoBuf(resident)
// 			protoResidents = append(protoResidents, protoResident)
// 		}
// 		return &models.LevelData{
// 			level:     toProtoBuf(s.Level),
// 			residents: protoResidents,
// 		}

// 	case BattleSession:
// 		var protoMonsters []*models.Monster
// 		for _, monster := range *s.Monsters {
// 			protoMonster := toProtoBuf(monster)
// 			protoMonsters = append(protoMonsters, protoMonster)
// 		}
// 		return &models.Sessions{
// 			battle_id:     s.BattleID,
// 			status:        s.Status,
// 			monsters:      protoMonsters,
// 			reward_matrix: s.RewardMatrix,
// 			reward:        toProtoBuf(reward),
// 		}

// 	case Reward:
// 		return &models.Sessions{
// 			gold:      s.Gold,
// 			exp:       s.Exp,
// 			total_exp: s.TotalExp,
// 		}

// 	case Packet:
// 		return &models.Sessions{
// 			uuid:         s.PacketID,
// 			packet_code:  s.PacketCode,
// 			chain:        s.Chain,
// 			service_type: s.ServiceType,
// 			content:      s.Content,
// 		}

// 	default:
// 		fmt.Println("Unsupported protobuf message type")
// 		return nil

// 	}
// }

func sayHiToClient() string {
	return "Hello Client - FROM SERVER"
}
func createMultiDeliveryPacket(requestIDSTR string, packetCode string, serviceType string, contentJSON []byte) Packet {
	var packet Packet
	requestID, _ := uuid.Parse(requestIDSTR)
	packet.PacketID = requestID
	packet.PacketCode = packetCode
	packet.Chain = true
	packet.ServiceType = serviceType
	cleanedJSON := "@\"" + string(contentJSON) + "\"@"
	packet.Content = string(cleanedJSON)
	return packet
}
func createSimpleDeliveryPacket(requestIDSTR string, packetCode string, serviceType string, content string) Packet {
	var packet Packet
	requestID, _ := uuid.Parse(requestIDSTR)
	packet.PacketID = requestID
	packet.PacketCode = packetCode
	packet.Chain = false
	packet.ServiceType = serviceType
	packet.Content = content
	return packet
}
func addPacketToCache(accountID uuid.UUID, packet Packet) {
	if _, ok := playerPacketCache[accountID]; !ok {
		var newCache PlayerPacketCache
		newCache.PacketCache = map[uuid.UUID]Packet{}
		playerPacketCache[accountID] = newCache
	}
	if playerPacketCache[accountID].PacketCache == nil {
		b := make(map[uuid.UUID]Packet)
		b[packet.PacketID] = packet
		playerPacketCache[accountID].PacketCache[packet.PacketID] = b[packet.PacketID]
	}
	playerPacketCache[accountID].PacketCache[packet.PacketID] = packet
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
func processTier3Packet(packetMessage string) (string, string, string) {
	item1 := strings.Split(packetMessage, "?")[0]
	item2 := strings.Split(packetMessage, "?")[1]
	item3 := strings.Split(packetMessage, "?")[2]
	return item1, item2, item3
}
func processTier4Packet(packetMessage string) (string, string, string, string) {
	item1 := strings.Split(packetMessage, "?")[0]
	item2 := strings.Split(packetMessage, "?")[1]
	item3 := strings.Split(packetMessage, "?")[2]
	item4 := strings.Split(packetMessage, "?")[3]
	return item1, item2, item3, item4
}
func processTier5Packet(packetMessage string) (string, string, string, string, string) {
	item1 := strings.Split(packetMessage, "?")[0]
	item2 := strings.Split(packetMessage, "?")[1]
	item3 := strings.Split(packetMessage, "?")[2]
	item4 := strings.Split(packetMessage, "?")[3]
	item5 := strings.Split(packetMessage, "?")[4]
	return item1, item2, item3, item4, item5
}
func getArrayFromString(packetMessage string) []int {
	var itemArray []int
	arrayWithoutBrackets := strings.Split(strings.Split(packetMessage, "[")[1], "]")[0]
	if !strings.Contains(arrayWithoutBrackets, ",") {
		item, _ := strconv.Atoi(arrayWithoutBrackets)
		itemArray = append(itemArray, item)
	} else {
		itemArrayStrings := strings.Split(arrayWithoutBrackets, ",")
		for _, element := range itemArrayStrings {
			item, _ := strconv.Atoi(element)
			itemArray = append(itemArray, item)
		}
	}
	return itemArray
}
func writeResponse(accountID uuid.UUID, requestID string, packet Packet, clientConnection net.Conn, resend bool) {
	packetJSON, _ := json.Marshal(packet)
	clientResponse := packet.PacketCode + "?" + string(packetJSON)
	if !resend {
		addPacketToCache(accountID, packet)
	}
	clientConnection.Write([]byte(strings.Trim(strconv.QuoteToASCII(clientResponse), "\"")))
}
func chainWriteResponse(accountID uuid.UUID, requestID string, packet Packet, byteLimiter int, clientConnection net.Conn, resend bool) {
	packetJSON, _ := json.Marshal(packet)
	packetData := "?" + string(packetJSON)
	if !resend {
		addPacketToCache(accountID, packet)
	}
	base := strings.Replace(packet.PacketCode, "#", "", -1)
	totalByteData := []byte(strings.Trim(strconv.QuoteToASCII(packetData), "\""))
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
		fmt.Println(Info("("+packet.ServiceType+") "+"Sent message back to client : ", clientResponse))
		clientConnection.Write([]byte(clientResponse))
	}
	fmt.Println(Info("Size of ", packet.ServiceType, " data in bytes : ", len(totalByteData)))
	fmt.Println(Info("Size of remaining ", packet.ServiceType, " in bytes : ", len(totalByteData)%byteLimiter))
	fmt.Println(Info("Size of ", packet.ServiceType, " partitions : ", dataPartitions))
}

func createBattle(monsters *[]Monster, quantity int) *BattleSession {
	var battle BattleSession
	battle.BattleID = uuid.New()
	battle.Status = 0
	var createdMonsters []Monster
	i := 0
	var reward Reward
	for i = 0; i < quantity; i++ {
		rand.Seed(time.Now().UnixNano())
		min := 0
		max := len(*monsters) - 1
		randomIndex := rand.Intn(max-min+1) + min
		monster := (*monsters)[randomIndex]
		createdMonsters = append(createdMonsters, monster)
		battle.Monsters = &createdMonsters
		battle.RewardMatrix = append(battle.RewardMatrix, 0)
		reward.Gold += float64(monster.GoldGain)
		reward.Exp += float64(monster.ExperienceGain)
	}
	sessions.Battles[battle.BattleID] = battle
	return &battle
}
func tcpListener(PORT string, cxt context.Context, mongoClient *mongo.Client) {
	listenerConnection, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(Failure(err))
		return
	}
	defer listenerConnection.Close()
	for {
		clientConnection, err := listenerConnection.Accept()
		if err != nil {
			fmt.Println(Failure(err))
			return
		}
		go handleTCPConnection(clientConnection, cxt, mongoClient)
	}
}
func udpListener(PORT string, cxt context.Context, mongoClient *mongo.Client) {
	buffer := make([]byte, 1024)
	serverAddress, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(Failure(err))
		return
	}
	listenerConnection, err := net.ListenUDP("udp4", serverAddress)
	if err != nil {
		fmt.Println(Failure(err))
		return
	}
	defer listenerConnection.Close()
	for {
		n, clientAddress, err := listenerConnection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(Failure(err))
			return
		}
		incomingPacket := string(buffer[0:n])
		fmt.Println(Success(clientAddress))
		fmt.Println(IncomingPacket(incomingPacket))
	}
}
func handleLogin(username string, password string, mongoClient *mongo.Client) (string, bool) {
	player, playerFound := getUser(username, mongoClient)
	if playerFound {
		if validateUser(player, mongoClient) {
			fmt.Println(Success("Login successful!"))
			playerJSON, _ := json.Marshal(player)
			response := fmt.Sprintf("Login successful;%v", string(playerJSON))
			return response, true
		} else {
			fmt.Println(Failure("Login failed!"))
			return "Login Failed;" + username + ";0", false
		}
	} else {
		fmt.Println(Failure("Login failed!"))
		return "Login failed;" + username + ";0", false
	}
}
func handleRegistration(username string, password string, mongoClient *mongo.Client) (string, bool, uuid.UUID) {
	if !lookForUser(username, mongoClient) {
		fmt.Println(Info("Password : ", password))
		accountID := createUser(username, password, mongoClient)
		return "Account created", true, accountID
	} else {
		fmt.Println(Warn(username, " is not available"))
		return "Username is not available", false, uuid.New()
	}
}
func createUser(username string, password string, mongoClient *mongo.Client) uuid.UUID {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database := mongoClient.Database("player")
	users := database.Collection("users")
	var newUser User
	newUser.ObjectID = primitive.NewObjectID()
	newUser.Account_id = uuid.New()
	newUser.User_id = username
	newUser.Password = password
	newUser.Active = 0
	newUser.Logins = 0
	createResult, err := users.InsertOne(cxt, newUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Success("New user added to db : ", createResult.InsertedID))
	return newUser.Account_id
}
func lookForUser(username string, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database := mongoClient.Database("player")
	users := database.Collection("users")
	filterCursor, err := users.Find(cxt, bson.M{"user_id": username})
	if err != nil {
		fmt.Println(Failure(err))
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
func validateUser(player *User, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	users := database.Collection("users")
	_, err := users.UpdateOne(
		cxt,
		bson.M{"uuid": player.Account_id},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "active", Value: 1}}},
			{Key: "$inc", Value: bson.D{{Key: "logins", Value: 1}}},
		}, options.Update().SetUpsert(true))
	if err != nil {
		fmt.Println(Failure(err))
		fmt.Println(Failure("Error with incrementing user login amount!"))
		return false
	}
	return true
}

//	func createShopKeeper(npcID string, mongoClient *mongo.Client) {
//		cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//		defer cancel()
//		if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
//			panic(err)
//		}
//		database := mongoClient.Database("world")
//		shopkeepers := database.Collection("shopkeepers")
//		var freshShopKeeper models.ShopKeeper
//		freshShopKeeper.ObjectID = primitive.NewObjectID().String()
//		freshShopKeeper.NpcId = npcID
//		freshShopKeeper.Catalogue = make([]*models.ShopItem)
//		var freshPurse Purse
//		freshPurse.Bits = 0
//		freshShopKeeper.Purse = freshPurse
//		createResult, err := shopkeepers.InsertOne(cxt, freshShopKeeper)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Println(Success("New shopkeeper added to db : ", createResult.InsertedID))
//	}
//
//	func getShopKeepersForServer(mongoClient *mongo.Client) {
//		cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//		defer cancel()
//		if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
//			panic(err)
//		}
//		database := mongoClient.Database("world")
//		shopkeepers := database.Collection("shopkeepers")
//		filterCursor, err := shopkeepers.Find(cxt, bson.D{})
//		var itemResult []models.ShopKeeper
//		if err = filterCursor.All(cxt, &itemResult); err != nil {
//			log.Fatal(err)
//		}
//		for _, shopkeeper := range itemResult {
//			ALLshopkeepers[shopkeeper.NpcId] = shopkeeper
//		}
//	}
func addCatalogueItem(itemID string, npcID string, price float64, mongoClient *mongo.Client) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	catalogueItem, foundItem := getItem(itemID, mongoClient)
	if foundItem {
		var freshShopItem ShopItem
		freshShopItem.Item_uuid = uuid.New()
		freshShopItem.Item = catalogueItem
		freshShopItem.Price = price
		database := mongoClient.Database("world")
		shopkeepers := database.Collection("shopkeepers")
		entryArea := "catalogue"
		match := bson.M{"npcID": npcID}
		change := bson.M{"$push": bson.M{entryArea: freshShopItem}}
		updateResponse, err := shopkeepers.UpdateOne(cxt, match, change)
		fmt.Println(Info(updateResponse))
		if err != nil {
			fmt.Println(Failure(err))
		} else {
			fmt.Println(Success("Item added successfully!"))
		}
	}
}
func getRegion(regionID string, mongoClient *mongo.Client) *Region {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	regions := database.Collection("regions")
	filterCursor, err := regions.Find(cxt, bson.M{"regionID": regionID})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []Region
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	if len(filterResult) == 1 {
		return &filterResult[0]
	}
	var dummyRegion Region
	return &dummyRegion
}
func getLevel(levelID string, mongoClient *mongo.Client) *Level {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	levels := database.Collection("levels")
	filterCursor, err := levels.Find(cxt, bson.M{"levelID": levelID})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []Level
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	if len(filterResult) == 1 {
		return &filterResult[0]
	}
	var dummyLevel Level
	return &dummyLevel
}
func getNPC(npcID string, mongoClient *mongo.Client) *Resident {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	npcs := database.Collection("npcs")
	filterCursor, err := npcs.Find(cxt, bson.M{"npcID": npcID})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []Resident
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	if len(filterResult) == 1 {
		return &filterResult[0]
	}
	var dummyNPC Resident
	return &dummyNPC
}
func getNPCs(npcIDs []string, mongoClient *mongo.Client) *[]Resident {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	npcs := database.Collection("npcs")
	filterCursor, err := npcs.Find(cxt, bson.M{"npcID": bson.M{"$in": npcIDs}})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []Resident
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	return &filterResult
}
func getMonster(monsterID string, mongoClient *mongo.Client) *Monster {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	monsters := database.Collection("monsters")
	filterCursor, err := monsters.Find(cxt, bson.M{"mobID": monsterID})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []Monster
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	if len(filterResult) == 1 {
		return &filterResult[0]
	}
	var dummyMonster Monster
	return &dummyMonster
}
func getMonsters(monsterIDs []string, mongoClient *mongo.Client) *[]Monster {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	monsters := database.Collection("monsters")
	filterCursor, err := monsters.Find(cxt, bson.M{"mobID": bson.M{"$in": monsterIDs}})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []Monster
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	return &filterResult
}
func getLevels(levelIDs []string, mongoClient *mongo.Client) *[]Level {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	levels := database.Collection("levels")
	filterCursor, err := levels.Find(cxt, bson.M{"levelID": bson.M{"$in": levelIDs}})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []Level
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	return &filterResult
}
func getUser(userID string, mongoClient *mongo.Client) (*User, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	users := database.Collection("users")
	filterCursor, err := users.Find(cxt, bson.M{"user_id": userID})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []User
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	if len(filterResult) == 1 {
		return &filterResult[0], true
	}
	var dummyUser User
	return &dummyUser, false
}
func updateUserLastPosition(target_uuid uuid.UUID, lastPosition *Position, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	user := database.Collection("profiles")
	match := bson.M{"uuid": target_uuid}
	change := bson.M{"$set": bson.D{
		{Key: "last_position", Value: lastPosition},
	}}
	_, err := user.UpdateOne(cxt, match, change)
	if err != nil {
		fmt.Println(Failure(err))
		return false
	}
	return true
}
func createProfile(userID string, accountID uuid.UUID, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profile := database.Collection("profiles")

	var newProfile Profile
	var defaultPosition Position
	newProfile.ObjectID = primitive.NewObjectID()
	newProfile.Account_id = accountID
	newProfile.Name = userID
	newProfile.Level = 1
	newProfile.Age = 0
	newProfile.Title = "Just a newbie"
	newProfile.Current_EXP = 0
	newProfile.Total_EXP = 0
	newProfile.Max_EXP = 100
	newProfile.Race_id = "human"
	newProfile.Race_name = "Human"
	newProfile.Class_id = "stranger"
	newProfile.Class_name = "Stranger"
	newProfile.LastPosition = defaultPosition
	newProfile.LastLevel = "00001"
	newProfile.LastRegion = "001"

	newProfile.Stats.Health = 100
	newProfile.Stats.Mana = 100
	newProfile.Stats.Attack = 1
	newProfile.Stats.MagicAttack = 1
	newProfile.Stats.Defense = 1
	newProfile.Stats.MagicDefense = 1
	newProfile.Stats.Accuracy = 1
	newProfile.Stats.Agility = 1
	newProfile.Items.Collection = make([]string, 0)

	var newPurse Purse
	newPurse.Bits = 0
	newProfile.Purse = newPurse

	var newLoadout Loadout
	newLoadout.Head = "EMPTY"
	newLoadout.Body = "EMPTY"
	newLoadout.Feet = "EMPTY"
	newLoadout.Weapon = "EMPTY"
	newLoadout.Accessory_1 = "EMPTY"
	newLoadout.Accessory_2 = "EMPTY"
	newLoadout.Accessory_3 = "EMPTY"
	newProfile.Loadout = newLoadout

	newProfile.SpellIndex = make([]string, 0)
	newProfile.BaseStats = newProfile.Stats

	insertResult, err := profile.InsertOne(cxt, newProfile)
	if err != nil {
		fmt.Println(Failure(err))
		return false
	}
	fmt.Println(Success("Fresh profile created for user: ", userID, " insertID: ", insertResult.InsertedID))
	return true
}
func getProfile(accountID uuid.UUID, mongoClient *mongo.Client) (*Profile, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profile := database.Collection("profiles")
	filterCursor, err := profile.Find(cxt, bson.M{"uuid": accountID})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []Profile
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	if len(filterResult) == 1 {
		return &filterResult[0], true
	}
	return nil, true
}
func updateProfile_EXP(accountID uuid.UUID, streamed_exp float64, mongoClient *mongo.Client) float64 {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	profile, profileFound := getProfile(accountID, mongoClient)
	newTotalExp := float64(profile.Total_EXP) + streamed_exp
	if profileFound {
		fmt.Println(Success("Found player profile to update"))
		database := mongoClient.Database("player")
		profiles := database.Collection("profiles")
		match := bson.M{"uuid": accountID}
		totalEXP := float64(profile.Current_EXP) + streamed_exp
		fmt.Println(Info(totalEXP, " (TotalEXP) = ", profile.Current_EXP, " (current exp) + ", streamed_exp, " (streamed exp)"))
		if totalEXP >= float64(profile.Max_EXP) {
			bufferEXP := 0.0
			levelUpperLimit := 0
			levelUpperLimitEXP := profile.Max_EXP
			bufferEXP = float64(profile.Max_EXP)
			if totalEXP > bufferEXP {
				for totalEXP > bufferEXP {
					levelUpperLimit++
					levelUpperLimitEXP += 50.0
					bufferEXP += float64(profile.Max_EXP) + float64(levelUpperLimit*50.0)
				}
			}
			newCurrentEXP := float64(levelUpperLimitEXP) - (bufferEXP - totalEXP)
			newLevel := int(profile.Level) + levelUpperLimit
			newMaxEXP := levelUpperLimitEXP
			change := bson.D{
				{Key: "$set", Value: bson.D{{Key: "profile.level", Value: newLevel}, {Key: "profile.current_exp", Value: newCurrentEXP}, {Key: "profile.max_exp", Value: newMaxEXP}, {Key: "profile.total_exp", Value: newTotalExp}}},
			}
			_, err := profiles.UpdateOne(cxt, match, change)
			if err != nil {
				fmt.Println(Failure(err))
				return 0
			}
			return newTotalExp
		} else if totalEXP < float64(profile.Max_EXP) {
			newCurrentEXP := totalEXP
			change := bson.D{
				{Key: "$set", Value: bson.D{{Key: "profile.current_exp", Value: newCurrentEXP}, {Key: "profile.total_exp", Value: newTotalExp}}},
			}
			_, err := profiles.UpdateOne(cxt, match, change)
			if err != nil {
				fmt.Println(Failure(err))
				return 0
			}
			return newTotalExp
		}
	}
	fmt.Println(Failure("Did not find player profile to update!"))
	return 0
}
func addBits(accountID uuid.UUID, streamed_bits float64, add bool, mongoClient *mongo.Client) float64 {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	profile, profileFound := getProfile(accountID, mongoClient)
	newTotalBits := 0.0
	if add {
		newTotalBits = float64(profile.Purse.Bits) + streamed_bits
	} else {
		newTotalBits = streamed_bits
	}
	if profileFound {
		database := mongoClient.Database("player")
		inventory := database.Collection("profiles")
		entryArea := "purse.bits"
		match := bson.M{"uuid": accountID}
		change := bson.M{"$set": bson.D{{Key: entryArea, Value: newTotalBits}}}
		_, err := inventory.UpdateOne(cxt, match, change)
		if err != nil {
			fmt.Println(Failure(err))
			return 0
		} else {
			fmt.Println(Success("Updated player bits!"))
		}
	}
	return newTotalBits
}
func getBits(accountID uuid.UUID, mongoClient *mongo.Client) float32 {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	profile, _ := getProfile(accountID, mongoClient)
	return float32(profile.Purse.Bits)
}
func addSpell(accountID uuid.UUID, spellID string, mongoClient *mongo.Client) string {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profiles := database.Collection("profiles")
	retrievedSpell, spellFound := getSpell(spellID, mongoClient)
	if spellFound {
		entryArea := "spell_index"
		match := bson.M{"uuid": accountID}
		change := bson.M{"$push": bson.M{entryArea: retrievedSpell.Spell_id}}
		updateResponse, err := profiles.UpdateOne(cxt, match, change)
		fmt.Println(Info(updateResponse))
		if err != nil {
			fmt.Println(Failure(err))
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
		fmt.Println(Failure(err))
		panic(err)
	}
	var spellResult []Spell
	if err = filterCursor.All(cxt, &spellResult); err != nil {
		log.Fatal(err)
	}
	if len(spellResult) == 1 {
		return &spellResult[0], true
	}
	return &spellResult[0], true
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
		fmt.Println(Failure(err))
		panic(err)
	}
	var itemResult []Item
	if err = filterCursor.All(cxt, &itemResult); err != nil {
		log.Fatal(err)
	}
	if len(itemResult) == 0 {
		fmt.Println(Warn("Item not found!"))
		emptyItem := new(Item)
		return *emptyItem, false
	} else if len(itemResult) == 0 {
		fmt.Println(Info("Item retrieved : ", itemResult[0].Item_id))
		return itemResult[0], true
	}
	return itemResult[0], true
}
func equipItem(accountID uuid.UUID, itemID string, mongoClient *mongo.Client) string {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profile := database.Collection("profiles")
	retrievedItem, itemFound := getItem(itemID, mongoClient)
	if itemFound {
		entryArea := retrievedItem.Item_type
		match := bson.M{"uuid": accountID}
		change := bson.M{"$set": bson.M{entryArea: retrievedItem.Item_id}}
		updateResponse, err := profile.UpdateOne(cxt, match, change)
		fmt.Println(Info(updateResponse))
		if err != nil {
			fmt.Println(Failure(err))
			return "EQUIP$0"
		}
		updateStats(accountID, retrievedItem, "add", mongoClient)
		return "EQUIP$1"
	}
	return "EQUIP$0"
}
func unequipItem(accountID uuid.UUID, itemID string, mongoClient *mongo.Client) string {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	loadout := database.Collection("profile")
	retrievedItem, itemFound := getItem(itemID, mongoClient)
	if itemFound {
		entryArea := retrievedItem.Item_type
		match := bson.M{"uuid": accountID}
		change := bson.M{"$set": bson.M{entryArea: "EMPTY"}}
		updateResponse, err := loadout.UpdateOne(cxt, match, change)
		fmt.Println(Info(updateResponse))
		if err != nil {
			fmt.Println(Failure(err))
			return "UNEQUIP$0"
		}
		updateStats(accountID, retrievedItem, "remove", mongoClient)
		return "UNEQUIP$1"
	}
	return "UNEQUIP$0"
}
func updateStats(accountID uuid.UUID, item Item, operation string, mongoClient *mongo.Client) *Profile {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profiles := database.Collection("profiles")
	profile, _ := getProfile(accountID, mongoClient)
	originalStats := profile.Stats
	updatedStats := updateStatsByItem(&originalStats, &item, operation)
	profile.Stats = *updatedStats
	match := bson.M{"uuid": accountID}
	change := bson.M{"$set": bson.D{{Key: "stats", Value: profile.Stats}}}
	updateResponse, err := profiles.UpdateOne(cxt, match, change)
	fmt.Println(Info(updateResponse))
	if err != nil {
		fmt.Println(Failure(err))
	}
	return profile
}
func addInventoryItem(accountID uuid.UUID, itemID string, mongoClient *mongo.Client) string {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profiles := database.Collection("profiles")
	retrievedItem, itemFound := getItem(itemID, mongoClient)
	if itemFound {
		entryArea := "items.collection"
		match := bson.M{"uuid": accountID}
		change := bson.M{"$push": bson.M{entryArea: retrievedItem.Item_id}}
		updateResponse, err := profiles.UpdateOne(cxt, match, change)
		fmt.Println(Info(updateResponse))
		if err != nil {
			fmt.Println(Failure(err))
			return "Item addition failed!"
		}
		return "Item added successfully!"
	}
	return "Item does not exist!"
}
func performTrade(accountID uuid.UUID, npcID string, playerBasket []Item, shopBasket []uuid.UUID, totalPlayerValue float64, totalShopValue float64, mongoClient *mongo.Client) bool {
	//access player purse
	//write func getBits()
	//access shopkeeper purse
	shopkeeper := ALLshopkeepers[npcID]
	shopkeeperBits := shopkeeper.Purse.Bits
	playerBits := getBits(accountID, mongoClient)
	//calculate total value of baskets
	pBasketValue := 0.0
	sBasketValue := 0.0
	for _, item := range playerBasket {
		//
		pBasketValue += item.BaseValue
	}
	for _, uuid := range shopBasket {
		for _, item := range ALLshopkeepers[npcID].Catalogue {
			if item.Item_uuid == uuid {
				sBasketValue += item.Price
			}
		}
	}
	pTradePower := pBasketValue + float64(playerBits)
	sTradePower := sBasketValue + shopkeeperBits
	if sTradePower > pTradePower && sTradePower >= sBasketValue {
		//buy player items first and then sell to player
		//if shop can afford player basket
		if sTradePower >= pBasketValue && pBasketValue > 0 {
			shopkeeperBits -= pBasketValue
			playerBits += float32(pBasketValue)
			pBasketValue = 0
			shopkeeperBits += sBasketValue
			playerBits -= float32(sBasketValue)
			sBasketValue = 0
			//update shopkeepers and player
			shopkeeper.Purse.Bits = shopkeeperBits
			ALLshopkeepers[npcID] = shopkeeper
			//create setBits function or modify addBits() with extra param
			addBits(accountID, float64(playerBits), false, mongoClient)
			return true
		} else {
			return false
		}
	} else if pTradePower > sTradePower && pTradePower >= pBasketValue {
		//buy shopkeeper items first and then sell to shopkeeper
		//if player can afford shopkeeper basket
		if pTradePower >= sBasketValue && sBasketValue > 0 {
			shopkeeperBits += sBasketValue
			playerBits -= float32(sBasketValue)
			sBasketValue = 0
			shopkeeperBits -= pBasketValue
			playerBits += float32(pBasketValue)
			pBasketValue = 0
			//update shopkeepers and player
			shopkeeper.Purse.Bits = shopkeeperBits
			ALLshopkeepers[npcID] = shopkeeper
			//create setBits function or modify addBits() with extra param
			addBits(accountID, float64(playerBits), false, mongoClient)
			return true
		} else {
			return false
		}
	}
	return false
}
func removeInventoryItem(accountID uuid.UUID, itemID string, mongoClient *mongo.Client) string {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profiles := database.Collection("profiles")
	item, itemFound := getItem(itemID, mongoClient)
	if itemFound {
		entryArea := "items.collection"
		match := bson.M{"uuid": accountID}
		change := bson.M{"$pull": bson.M{entryArea: bson.M{"$in": item.Item_id}}}
		updateResponse, err := profiles.UpdateOne(cxt, match, change)
		fmt.Println(Info(updateResponse))
		if err != nil {
			fmt.Println(Failure(err))
			return "Item deletion failed!"
		}
		return "Item deletion successfully!"
	}
	return "Item does not exist!"
}
func getItemsGlobalAndCache(mongoClient *mongo.Client) []Item {
	//get all items from world/items
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	items := database.Collection("items")
	//filterCursor, err := regions.Find(cxt, bson.M{"regionID": regionID})
	filterCursor, err := items.Find(cxt, bson.D{})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []Item
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	//cache the items to a global map
	for _, item := range filterResult {
		MASTER_ITEM_TABLE[item.Item_id] = item
	}
	return filterResult
}
func getSpellsGlobalAndCache(mongoClient *mongo.Client) []Spell {
	//get all items from world/items
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	items := database.Collection("spells")
	filterCursor, err := items.Find(cxt, bson.D{})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []Spell
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	//cache the items to a global map
	for _, spell := range filterResult {
		MASTER_SPELL_TABLE[spell.Spell_id] = spell
	}
	return filterResult
}
func getMonstersGlobalAndCache(mongoClient *mongo.Client) []Monster {
	//get all items from world/items
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	items := database.Collection("monsters")
	filterCursor, err := items.Find(cxt, bson.D{})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []Monster
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	//cache the items to a global map
	for _, monster := range filterResult {
		MASTER_MONSTER_TABLE[monster.MobID] = monster
	}
	return filterResult
}
func getLevelsGlobalAndCache(mongoClient *mongo.Client) []Level {
	//get all items from world/items
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("world")
	items := database.Collection("levels")
	filterCursor, err := items.Find(cxt, bson.D{})
	if err != nil {
		fmt.Println(Failure(err))
		panic(err)
	}
	var filterResult []Level
	if err = filterCursor.All(cxt, &filterResult); err != nil {
		log.Fatal(err)
	}
	//cache the items to a global map
	for _, level := range filterResult {
		MASTER_LEVEL_TABLE[level.LevelID] = level
	}
	return filterResult
}
func updateStatsByItem(originalStats *Stats, item *Item, operation string) *Stats {
	op := 0.0
	if operation == "add" {
		op = 1.0
	} else if operation == "remove" {
		op = -1.0
	}
	originalStats.Strength += (op * item.Stats.Strength)
	originalStats.Intelligence += (op * item.Stats.Intelligence)
	originalStats.Dexterity += (op * item.Stats.Dexterity)
	originalStats.Charisma += (op * item.Stats.Charisma)
	originalStats.Luck += (op * item.Stats.Luck)
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
func getSpellFromCache(spell_id string) Spell {
	fmt.Println(Internal("Spell retrieved from cache!"))
	return MASTER_SPELL_TABLE[spell_id]
}
func getLevelFromCache(level_id string) Level {
	fmt.Println(Internal("Level retrieved from cache!"))
	return MASTER_LEVEL_TABLE[level_id]
}

func main() {
	//mongoDB specs
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://admin:w1583069@cluster0.8acnf.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
	playerPacketCache = map[uuid.UUID]PlayerPacketCache{}
	//initialize mongoDB client
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(cxt, opts)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(Success("Connected to MongoDB cluster0..."))
	}

	//-------------------------------------------------------------------------------------------------------------------
	//TODO - Get a copy of the master Item and Spell collection from the database when the server is connected and online
	//-------------------------------------------------------------------------------------------------------------------

	//createShopKeeper("NPC1", mongoClient)
	//addCatalogueItem("WizardHat", "NPC1", 7, mongoClient)
	//createShopKeeper("NPC0", mongoClient)
	//addCatalogueItem("WizardRobe", "NPC0", 23, mongoClient)
	// addInventoryItem("asd", "WizardHat", mongoClient)
	// addInventoryItem("asd", "WizardHat", mongoClient)

	//getShopKeepersForServer(mongoClient)
	// get all global spell and item data and cache it during server runtime
	getItemsGlobalAndCache(mongoClient)
	getSpellsGlobalAndCache(mongoClient)
	getMonstersGlobalAndCache(mongoClient)
	getLevelsGlobalAndCache(mongoClient)
	//disconnect mongoDB client on return
	defer func() {
		if err = mongoClient.Disconnect(cxt); err != nil {
			panic(err)
		}
	}()
	sessions.Battles = make(map[uuid.UUID]BattleSession)
	// add to sync.WaitGroup
	fmt.Println(Success("SERVER RUNNING on Port 20001"))
	wg.Add(1)
	go tcpListener(":20001", cxt, mongoClient)
	wg.Add(2)
	go udpListener(":26950", cxt, mongoClient)
	wg.Wait()
}
