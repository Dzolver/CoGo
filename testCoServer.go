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

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	ObjectID     primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	Account_id   uuid.UUID          `json:"uuid" bson:"uuid,omitempty"`
	User_id      string             `json:"user_id" default:"" bson:"user_id, omitempty"`
	Password     string             `json:"password" default:"" bson:"password, omitempty"`
	Active       int                `json:"active" default:"" bson:"active,omitempty"`
	Logins       int                `json:"logins" default:"" bson:"logins, omitempty"`
	LastPosition Position           `json:"last_position" bson:"last_position,omitempty"`
	LastRegion   string             `json:"last_region" default:"" bson:"last_region,omitempty"`
	LastLevel    string             `json:"last_level" default:"" bson:"last_level,omitempty"`
}
type PlayerInventory struct {
	ObjectID   primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	Account_id uuid.UUID          `json:"uuid" bson:"uuid,omitempty"`
	Purse      Purse              `json:"purse" bson:"purse, omitempty"`
	Equipment  Equipment          `json:"equipment" bson:"equipment, omitempty"`
	Items      Items              `json:"items" bson:"items, omitempty"`
}
type PlayerLoadout struct {
	ObjectID    primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	Account_id  uuid.UUID          `json:"uuid" bson:"uuid,omitempty"`
	Head        Item               `json:"head" bson:"Head, omitempty"`
	Body        Item               `json:"body" bson:"Body, omitempty"`
	Feet        Item               `json:"feet" bson:"Feet, omitempty"`
	Weapon      Item               `json:"weapon" bson:"Weapon, omitempty"`
	Accessory_1 Item               `json:"accessory_1" bson:"Accessory_1, omitempty"`
	Accessory_2 Item               `json:"accessory_2" bson:"Accessory_2, omitempty"`
	Accessory_3 Item               `json:"accessory_3" bson:"Accessory_3, omitempty"`
}

type Purse struct {
	Bits float64 `default:"0" json:"bits" bson:"bits, omitempty"`
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
	Account_id  uuid.UUID          `json:"uuid" bson:"uuid,omitempty"`
	Spell_index []Spell            `json:"spell_index" bson:"spell_index, omitempty"`
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
	Account_id    uuid.UUID          `json:"uuid" bson:"uuid,omitempty"`
	PlayerProfile PlayerProfile      `json:"profile" default:"" bson:"profile, omitempty"`
	Stats         Stats              `json:"stats" default:"" bson:"stats, omitempty"`
	BaseStats     Stats              `json:"base_stats" default:"" bson:"base_stats, omitempty"`
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
	BattleID        uuid.UUID         `json:"battle_id" default:"" bson:"battle_id"`
	Vital           *PlayerVital      `json:"vital" default:"" bson:"vital, omitempty"`
	Inventory       *PlayerInventory  `json:"inventory" default:"" bson:"inventory, omitempty"`
	SpellIndex      *PlayerSpellIndex `json:"spellIndex" default:"" bson:"spellIndex, omitempty"`
	Loadout         *PlayerLoadout    `json:"loadout" default:"" bson:"loadout, omitempty"`
	Monsters        *[]Monster        `json:"monsters" default:"" bson:"monsters, omitempty"`
	MonsterQuantity int               `json:"monster_quantity" default:"0" bson:"monster_quantity, omitempty"`
}
type LoginSecretPacket struct {
	User       *User             `json:"user_data" default:"" bson:"user_data,omitempty"`
	Inventory  *PlayerInventory  `json:"inventory_data" default:"" bson:"inventory_data,omitempty"`
	SpellIndex *PlayerSpellIndex `json:"spell_index_data" default:"" bson:"spell_index_data,omitempty"`
	Loadout    *PlayerLoadout    `json:"loadout_data" default:"" bson:"loadout_data,omitempty"`
	Vital      *PlayerVital      `json:"vital_data" default:"" bson:"vital_data,omitempty"`
	Region     *RegionData       `json:"region_data" default:"" bson:"region_data,omitempty"`
}
type Map struct {
	ConnectedClients map[uuid.UUID]Client `json:"connected_clients" bson:"connected_clients, omitempty"`
}
type Party struct {
	Party_id uuid.UUID   `json:"party_id" bson:"party_id,omitempty"`
	Members  []uuid.UUID `json:"members" bson:"members,omitempty"`
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
type Monster struct {
	ObjectID       primitive.ObjectID `json:"objectID" bson:"_id, omitempty"`
	MobID          string             `json:"mob_id" default:"" bson:"mobID,omitempty"`
	MonsterType    string             `json:"monster_type" default:"" bson:"monsterType,omitempty"`
	GoldGain       int                `json:"gold_gain" default:"" bson:"goldGain,omitempty"`
	ExperienceGain int                `json:"experience_gain" default:"" bson:"experienceGain,omitempty"`
	MobVitals      *MobVitals         `json:"mob_vitals" default:"" bson:"mobVitals,omitempty"`
}
type MobVitals struct {
	Profile        *MobProfile `json:"profile" default:"" bson:"profile,omitempty"`
	Stats          *MobStats   `json:"stats" default:"" bson:"stats,omitempty"`
	AttackActions  *[]Spell    `json:"attack_actions" default:"" bson:"attackActions,omitempty"`
	DefenseActions *[]Spell    `json:"defense_actions" default:"" bson:"defenseActions,omitempty"`
}
type MobProfile struct {
	Name         string   `json:"name" default:"" bson:"name, omitempty"`
	Level        int      `json:"level" default:"0" bson:"level, omitempty"`
	Element      string   `json:"element" default:"" bson:"element, omitempty"`
	Regions      []string `json:"regions" bson:"regions, omitempty"`
	Description  string   `json:"description" default:"" bson:"description, omitempty"`
	Strength     float64  `json:"strength" default:"0" bson:"strength, omitempty"`
	Intelligence float64  `json:"intelligence" default:"0" bson:"intelligence, omitempty"`
	Dexterity    float64  `json:"dexterity" default:"0" bson:"dexterity, omitempty"`
	Charisma     float64  `json:"charisma" default:"0" bson:"charisma, omitempty"`
	Luck         float64  `json:"luck" default:"0" bson:"luck, omitempty"`
}
type MobStats struct {
	Health       float64 `json:"health" default:"0" bson:"health"`
	MaxHealth    float64 `json:"maxHealth" default:"0" bson:"maxHealth"`
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

var count = 0
var PACKET_SIZE = 10240
var wg sync.WaitGroup
var connectedUsers = make(map[string]bool)
var portNumbers = make(map[string]int)
var portNumbersReversed = make(map[string]string)
var portIndex = 1
var mapInstance Map
var sessions Sessions
var (
	Info           = Teal
	IncomingPacket = Magenta
	Warn           = Yellow
	Fata           = Red
	Success        = Green
	Failure        = Yellow
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
		//fmt.Println(netData)
		if packetCode == "BR#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Read for Battle packet received!"))
			accountIDSTR, levelID := processTier2Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			var freshBattlePacket BattlePacket
			vital, _ := getVital(accountID, mongoClient)
			inventory, _ := getInventory(accountID, mongoClient)
			spellIndex, _ := getSpellIndex(accountID, mongoClient)
			loadout, _ := getLoadout(accountID, mongoClient)
			level := getLevel(levelID, mongoClient)
			monsters := getMonsters(level.Monsters, mongoClient)
			//create BattleSession out of this information and add to the list of sessions
			freshBattlePacket.MonsterQuantity = 1
			battle := createBattle(monsters, freshBattlePacket.MonsterQuantity)

			freshBattlePacket.BattleID = battle.BattleID
			freshBattlePacket.Vital = vital
			freshBattlePacket.Inventory = inventory
			freshBattlePacket.SpellIndex = spellIndex
			freshBattlePacket.Loadout = loadout
			freshBattlePacket.Monsters = monsters

			battlePacketJSON, _ := json.Marshal(freshBattlePacket)
			battlePlayerData := "?" + string(battlePacketJSON)
			chainWriteResponse(packetCode, battlePlayerData, byteLimiter, clientConnection, "BATTLE")
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
			accountIDSTR, battleIDSTR, rewardMatrixSTR := processTier3Packet(packetMessage)
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
				addBits(accountID, gold, mongoClient)
				entry.Reward.TotalExp = updateProfile_EXP(accountID, entry.Reward.Exp, mongoClient)
				sessions.Battles[battleID] = entry
				updateStatus = "True"
			}
			inventory, _ := getInventory(accountID, mongoClient)
			inventoryJSON, _ := json.Marshal(inventory)
			battleFinishData := "?" + strconv.FormatFloat(exp, 'f', -1, 64) + "|" + strconv.FormatFloat(gold, 'f', -1, 64) + "|" + updateStatus + "|" + string(inventoryJSON)
			chainWriteResponse(packetCode, battleFinishData, byteLimiter, clientConnection, "BATTLEFINISH")
		}
		if packetCode == "HB#" {
			// fmt.Println("Heartbeat packet received!")
			accountID, x, y, z := processTier4Packet(packetMessage)
			target_uuid, _ := uuid.Parse(accountID)
			var lastPosition Position
			lastPosition.Position_x, _ = strconv.ParseFloat(x, 64)
			lastPosition.Position_y, _ = strconv.ParseFloat(y, 64)
			lastPosition.Position_z, _ = strconv.ParseFloat(z, 64)
			updateUserLastPosition(target_uuid, lastPosition, mongoClient)
		}
		//Inventory add
		if packetCode == "IA#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Add Inventory packet received!"))
			accountIDSTR, itemID := processTier2Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			clientResponse += addInventoryItem(accountID, itemID, mongoClient)
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "ID#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Delete Inventory packet received!"))
		}
		if packetCode == "IR#" {
			fmt.Println(IncomingPacket("Read Inventory packet received!"))
			accountIDSTR := processTier1Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			inventory, _ := getInventory(accountID, mongoClient)
			inventoryJSON, _ := json.Marshal(inventory)
			inventoryData := "?" + string(inventoryJSON)
			chainWriteResponse(packetCode, inventoryData, byteLimiter, clientConnection, "INVENTORY")
		}
		//Inventory update
		if packetCode == "IU#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Update Inventory packet received!"))
			accountIDSTR, itemID := processTier2Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			clientResponse += "?" + addInventoryItem(accountID, itemID, mongoClient)
			writeResponse(clientResponse, clientConnection)
		}
		//Login
		if packetCode == "L0#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Login packet received!"))
			fmt.Println(Info(packetMessage))
			username, password := processTier2Packet(packetMessage)
			loginResponse, valid := handleLogin(username, password, mongoClient)
			if valid {
				//Login success
				clientResponse = "LS#"
				var freshlsp LoginSecretPacket
				User, _ := getUser(username, mongoClient)
				Inventory, _ := getInventory(User.Account_id, mongoClient)
				Loadout, _ := getLoadout(User.Account_id, mongoClient)
				SpellIndex, _ := getSpellIndex(User.Account_id, mongoClient)
				Vital, _ := getVital(User.Account_id, mongoClient)
				freshlsp.User = User
				freshlsp.Inventory = Inventory
				freshlsp.Loadout = Loadout
				freshlsp.SpellIndex = SpellIndex
				freshlsp.Vital = Vital

				var freshRegionData *RegionData = new(RegionData)
				Region := getRegion(User.LastRegion, mongoClient)
				freshRegionData.Region = Region

				var freshLevelData *LevelData = new(LevelData)
				Level := getLevel(User.LastLevel, mongoClient)
				Residents := getNPCs(Level.Residents, mongoClient)
				freshLevelData.Level = Level
				freshLevelData.Residents = Residents

				freshRegionData.LevelData = freshLevelData

				freshlsp.Region = freshRegionData

				freshlspJSON, _ := json.Marshal(freshlsp)
				lspData := "?" + string(freshlspJSON)
				loginResponse = string(freshlspJSON)
				chainWriteResponse("LS#", lspData, byteLimiter, clientConnection, "LOGIN SECRET PACKET")
			} else if !valid {
				//Login fail
				clientResponse = "LF#"
				clientResponse += "?" + loginResponse
				fmt.Println(Info(clientResponse))
				writeResponse(clientResponse, clientConnection)
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
			accountIDSTR, itemID := processTier2Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			equipFeedback := equipItem(accountID, itemID, mongoClient)
			clientResponse += equipFeedback
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "LL#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Level packet received!"))
			levelID := processTier1Packet(packetMessage)
			var freshLevel LevelData
			level := getLevel(levelID, mongoClient)
			NPC := getNPCs(level.Residents, mongoClient)
			freshLevel.Level = level
			freshLevel.Residents = NPC
			levelDataJSON, _ := json.Marshal(freshLevel)
			levelData := "?" + string(levelDataJSON)
			chainWriteResponse(packetCode, levelData, byteLimiter, clientConnection, "LEVEL")
		}
		if packetCode == "LR#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Read loadout packet received!"))
			accountIDSTR := processTier1Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			loadout, _ := getLoadout(accountID, mongoClient)
			loadoutJSON, _ := json.Marshal(loadout)
			loadoutData := "?" + string(loadoutJSON)
			chainWriteResponse(packetCode, loadoutData, byteLimiter, clientConnection, "LOADOUT")
		}
		if packetCode == "LU#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Update EXP packet received!"))
			accountIDSTR, streamedEXPString := processTier2Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			streamedEXP, _ := strconv.ParseFloat(streamedEXPString, 64)
			newTotalEXP := updateProfile_EXP(accountID, streamedEXP, mongoClient)
			newTotalEXPSTR := strconv.FormatFloat(newTotalEXP, 'E', -1, 64)
			clientResponse += "?" + newTotalEXPSTR
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "LUE#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Unequip from loadout"))
			accountIDSTR, itemID := processTier2Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			unequipFeedback := unequipItem(accountID, itemID, mongoClient)
			clientResponse += unequipFeedback
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "ML#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Map load packet received"))
			mapJSON, _ := json.Marshal(mapInstance.ConnectedClients)
			mapData := "?" + string(mapJSON)
			chainWriteResponse(packetCode, mapData, byteLimiter, clientConnection, "MAP")
		}
		//Register
		if packetCode == "R0#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Register packet received!"))
			fmt.Println(Info(packetMessage))
			username, password := processTier2Packet(packetMessage)
			registerResponse, valid, accountID := handleRegistration(username, password, mongoClient)
			if valid {
				//Register success
				createInventory(accountID, mongoClient)
				createLoadout(accountID, mongoClient)
				createSpellIndex(accountID, mongoClient)
				createVital(username, accountID, mongoClient)
				addSpell(accountID, "Fireball", mongoClient)
				addSpell(accountID, "Scorch", mongoClient)
				clientResponse = "RS#"
			} else if !valid {
				//Register fail
				clientResponse = "RF#"
			}
			clientResponse += "?" + registerResponse
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "RLL#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Region and Level packet received!"))
			fmt.Println(Info(packetMessage))
			regionID, levelID := processTier2Packet(packetMessage)
			var freshRegionData RegionData
			var freshLevelData LevelData
			region := getRegion(regionID, mongoClient)
			level := getLevel(levelID, mongoClient)
			NPC := getNPCs(level.Residents, mongoClient)
			freshRegionData.Region = region
			freshLevelData.Level = level
			freshLevelData.Residents = NPC
			freshRegionData.LevelData = &freshLevelData
			regionDataJSON, _ := json.Marshal(freshRegionData)
			regionData := "?" + string(regionDataJSON)
			chainWriteResponse(packetCode, regionData, byteLimiter, clientConnection, "REGION")
		}

		if packetCode == "SR#" {
			fmt.Println(IncomingPacket("Read Spell Index packet received!"))
			accountIDSTR := processTier1Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			spellIndex, _ := getSpellIndex(accountID, mongoClient)
			spellIndexJSON, _ := json.Marshal(spellIndex)
			spellIndexData := "?" + string(spellIndexJSON)
			chainWriteResponse(packetCode, spellIndexData, byteLimiter, clientConnection, "SPELLINDEX")
		}
		if packetCode == "SU#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Update Spell Index packet received!"))
			accountIDSTR, spellID := processTier2Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			clientResponse += "?" + addSpell(accountID, spellID, mongoClient)
			writeResponse(clientResponse, clientConnection)
		}
		if packetCode == "VR#" {
			clientResponse = packetCode
			fmt.Println(IncomingPacket("Load vital packet received!"))
			accountIDSTR := processTier1Packet(packetMessage)
			accountID, _ := uuid.Parse(accountIDSTR)
			vital, _ := getVital(accountID, mongoClient)
			vitalJSON, _ := json.Marshal(vital)
			vitalData := "?" + string(vitalJSON)
			chainWriteResponse(packetCode, vitalData, byteLimiter, clientConnection, "PROFILE")
		}
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
		fmt.Println(Info("("+serviceType+") "+"Sent message back to client : ", clientResponse))
		clientConnection.Write([]byte(clientResponse))
	}
	fmt.Println(Info("Size of ", serviceType, " data in bytes : ", len(totalByteData)))
	fmt.Println(Info("Size of remaining ", serviceType, " in bytes : ", len(totalByteData)%byteLimiter))
	fmt.Println(Info("Size of ", serviceType, " partitions : ", dataPartitions))
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
		count++
	}
}
func handleNewUDPConnection(accountID string, broadcastAddress *net.UDPAddr, clientAddress *net.UDPAddr) bool {
	fmt.Println("Handling new UDP Connection!")
	connected := false
	target_uuid, _ := uuid.Parse(accountID)
	if _, existing := mapInstance.ConnectedClients[target_uuid]; !existing {
		//add new client connection to client map instance
		var freshClient Client
		target_uuid, _ := uuid.Parse(accountID)
		freshClient.Account_id = target_uuid
		freshClient.ConnectTime = time.Now()
		freshClient.UDPAddress = clientAddress
		freshClient.BroadcastAddress = broadcastAddress
		mapInstance.ConnectedClients[target_uuid] = freshClient
		for key, element := range mapInstance.ConnectedClients {
			fmt.Println("uuid:", key, "=>", "client broadcast address:", element.BroadcastAddress)
		}
		connected = true
	}
	return connected
}
func handleUDPConnection(netData string, clientAddress *net.UDPAddr, listenerConnection *net.UDPConn, mongoClient *mongo.Client) {
	packetCode, packetMessage := packetDissect(netData)
	if packetCode != "" {
		//fmt.Println("UDP Net data:" + netData)
		//fmt.Println("Packetcode : " + packetCode)
		//fmt.Println("Packet msg : " + packetMessage)
		if packetCode == "UDPC#" {
			clientResponse := packetCode
			fmt.Println("Start UDP stream packet received!")
			accountID, broadcastAddr := processTier2Packet(packetMessage)
			broadcastAddress, _ := net.ResolveUDPAddr("udp", broadcastAddr)
			fmt.Println("User joined : " + accountID)
			connected := handleNewUDPConnection(accountID, broadcastAddress, clientAddress)
			if connected {
				broadcastNewPlayer(accountID)
			}
			clientResponse += "?CONNECTED TO SERVER:" + strconv.FormatBool(connected)
			data := []byte(clientResponse)
			listenerConnection.WriteToUDP(data, clientAddress)
		}
		if packetCode == "M1#" {
			accountID, x, y, z := processTier4Packet(packetMessage)
			target_uuid, _ := uuid.Parse(accountID)
			if player, existing := mapInstance.ConnectedClients[target_uuid]; existing {
				player.Position.Position_x, _ = strconv.ParseFloat(x, 64)
				player.Position.Position_y, _ = strconv.ParseFloat(y, 64)
				player.Position.Position_z, _ = strconv.ParseFloat(z, 64)
				//updateUserLastPosition(target_uuid, player.Position, mongoClient)
				mapInstance.ConnectedClients[target_uuid] = player
				// clientResponse := "Got your M1# message!"
				// data := []byte(clientResponse)
				// listenerConnection.WriteToUDP(data, clientAddress)
			}
			//broadcastSend <- mapInstance
		}
	}
}
func broadcastMap() {
	local, _ := net.ResolveUDPAddr("udp4", ":6666")
	broadcastAddress, _ := net.ResolveUDPAddr("udp", "255.255.255.255"+":23456")
	connection, _ := net.DialUDP("udp", local, broadcastAddress)
	defer connection.Close()
	mapData := ""
	// for _, client := range mapInstance.ConnectedClients {
	// 	clientJSON, _ := json.Marshal(client)
	// 	movementData += "?" + string(clientJSON)
	// }
	mapJSON, _ := json.Marshal(mapInstance.ConnectedClients)
	mapData += string(mapJSON)
	broadcastData := "BRO#"
	broadcastData += mapData
	//fmt.Println("MAP COUNT : ", len(mapInstance.ConnectedClients))
	//fmt.Println("MAP DATA : " + mapData)
	_, err := connection.Write([]byte(mapData))
	if err != nil {
		panic(err)
	}
}
func broadcastNewPlayer(accountID string) {
	local, _ := net.ResolveUDPAddr("udp4", ":6666")
	broadcastAddress, _ := net.ResolveUDPAddr("udp", "255.255.255.255"+":23456")
	connection, _ := net.DialUDP("udp", local, broadcastAddress)
	defer connection.Close()
	broadcastData := "NP#"
	broadcastData += accountID
	fmt.Println("New Player : " + broadcastData)
	_, err := connection.Write([]byte(broadcastData))
	if err != nil {
		panic(err)
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
		//fmt.Println(n, " UDP client address : ", clientAddress)
		handleUDPConnection(incomingPacket, clientAddress, listenerConnection, mongoClient)
		if len(mapInstance.ConnectedClients) > 0 {
			broadcastMap()
		}
		// if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
		// 	fmt.Println("UDP client has exited!")
		// 	count--
		// 	break
		// }
		// data := []byte("UDP server acknowledges!")
		// v, err := listenerConnection.WriteToUDP(data, clientAddress)
		// if err != nil {
		// 	fmt.Println(Failure(err))
		// 	return
		// }
		// fmt.Print(v)
		// count++
		// processLimitChan <- count
	}
}
func handleLogin(username string, password string, mongoClient *mongo.Client) (string, bool) {
	player, playerFound := getUser(username, mongoClient)
	if playerFound {
		if validateUser(player, mongoClient) {
			fmt.Println(Success("Login successful!"))
			playerJSON, _ := json.Marshal(player)
			//positionJSON, _ := json.Marshal(player.LastPosition)
			//response := fmt.Sprintf("Login successful;%v;%v;%v;%v", player.Account_id, player.User_id, player.Logins, string(positionJSON))
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
	var freshUser User
	freshUser.ObjectID = primitive.NewObjectID()
	freshUser.Account_id = uuid.New()
	freshUser.User_id = username
	freshUser.Password = password
	freshUser.Active = 0
	freshUser.Logins = 0
	freshUser.LastPosition = Position{}
	freshUser.LastLevel = "00001"
	freshUser.LastRegion = "001"
	createResult, err := users.InsertOne(cxt, freshUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Success("New user added to db : ", createResult.InsertedID))
	return freshUser.Account_id
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
			{"$set", bson.D{{"active", 1}}},
			{"$inc", bson.D{{"logins", 1}}},
		}, options.Update().SetUpsert(true))
	if err != nil {
		fmt.Println(Failure(err))
		fmt.Println(Failure("Error with incrementing user login amount!"))
		return false
	}
	return true
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
func updateUserLastPosition(target_uuid uuid.UUID, lastPosition Position, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	user := database.Collection("users")
	match := bson.M{"uuid": target_uuid}
	change := bson.M{"$set": bson.D{
		{"last_position", lastPosition},
	}}
	_, err := user.UpdateOne(cxt, match, change)
	//fmt.Printf("Updated %v Documents!\n", updateResponse.ModifiedCount)
	if err != nil {
		fmt.Println(Failure(err))
		return false
	}
	return true
}
func createVital(userID string, accountID uuid.UUID, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profile := database.Collection("vital")

	var freshVital PlayerVital
	freshVital.ObjectID = primitive.NewObjectID()
	freshVital.Account_id = accountID
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
		fmt.Println(Failure(err))
		return false
	}
	fmt.Println(Success("Fresh Vital created for user: ", userID, " insertID: ", insertResult.InsertedID))
	return true
}
func getVital(accountID uuid.UUID, mongoClient *mongo.Client) (*PlayerVital, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profile := database.Collection("vital")
	filterCursor, err := profile.Find(cxt, bson.M{"uuid": accountID})
	if err != nil {
		fmt.Println(Failure(err))
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
func updateProfile_EXP(accountID uuid.UUID, streamed_exp float64, mongoClient *mongo.Client) float64 {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	playerVital, profileFound := getVital(accountID, mongoClient)
	newTotalExp := playerVital.PlayerProfile.Total_EXP + streamed_exp
	if profileFound {
		fmt.Println(Success("Found player profile to update"))
		database := mongoClient.Database("player")
		profile := database.Collection("vital")
		match := bson.M{"uuid": accountID}
		totalEXP := playerVital.PlayerProfile.Current_EXP + streamed_exp
		fmt.Println(Info(totalEXP, " (TotalEXP) = ", playerVital.PlayerProfile.Current_EXP, " (current exp) + ", streamed_exp, " (streamed exp)"))
		if totalEXP >= playerVital.PlayerProfile.Max_EXP {
			bufferEXP := 0.0
			levelUpperLimit := 0
			levelUpperLimitEXP := playerVital.PlayerProfile.Max_EXP
			bufferEXP = playerVital.PlayerProfile.Max_EXP
			for totalEXP > bufferEXP {
				levelUpperLimitEXP += float64(levelUpperLimit * 50.0)
				bufferEXP += playerVital.PlayerProfile.Max_EXP + float64(levelUpperLimit*50.0)
				levelUpperLimit++
			}
			newCurrentEXP := levelUpperLimitEXP - (bufferEXP - totalEXP)
			newLevel := playerVital.PlayerProfile.Level + levelUpperLimit
			newMaxEXP := levelUpperLimitEXP
			change := bson.D{
				{"$set", bson.D{{"profile.level", newLevel}, {"profile.current_exp", newCurrentEXP}, {"profile.max_exp", newMaxEXP}, {"profile.total_exp", newTotalExp}}},
			}
			_, err := profile.UpdateOne(cxt, match, change)
			if err != nil {
				fmt.Println(Failure(err))
				return 0
			}
			return newTotalExp
		} else if totalEXP < playerVital.PlayerProfile.Max_EXP {
			newCurrentEXP := totalEXP
			change := bson.D{
				{"$set", bson.D{{"profile.current_exp", newCurrentEXP}, {"profile.total_exp", newTotalExp}}},
			}
			_, err := profile.UpdateOne(cxt, match, change)
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
func addBits(accountID uuid.UUID, streamed_bits float64, mongoClient *mongo.Client) float64 {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	playerInventory, inventoryFound := getInventory(accountID, mongoClient)
	newTotalBits := playerInventory.Purse.Bits + streamed_bits
	if inventoryFound {
		database := mongoClient.Database("player")
		inventory := database.Collection("inventory")
		entryArea := "purse.bits"
		match := bson.M{"uuid": accountID}
		change := bson.M{"$set": bson.D{{entryArea, newTotalBits}}}
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
func createSpellIndex(accountID uuid.UUID, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	spellIndex := database.Collection("spellIndex")

	var freshSpellIndex PlayerSpellIndex
	freshSpellIndex.ObjectID = primitive.NewObjectID()
	freshSpellIndex.Account_id = accountID
	freshSpellIndex.Spell_index = make([]Spell, 0)

	insertResult, err := spellIndex.InsertOne(cxt, freshSpellIndex)
	if err != nil {
		fmt.Println(Failure(err))
		return false
	}
	fmt.Println("Fresh Spell index created! insertID: ", insertResult.InsertedID)
	return true
}
func getSpellIndex(accountID uuid.UUID, mongoClient *mongo.Client) (*PlayerSpellIndex, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	spellIndex := database.Collection("spellIndex")
	filterCursor, err := spellIndex.Find(cxt, bson.M{"uuid": accountID})
	if err != nil {
		fmt.Println(Failure(err))
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
func addSpell(accountID uuid.UUID, spellID string, mongoClient *mongo.Client) string {
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
		match := bson.M{"uuid": accountID}
		change := bson.M{"$push": bson.M{entryArea: retrievedSpell}}
		updateResponse, err := spellIndex.UpdateOne(cxt, match, change)
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
	if len(spellResult) == 0 {
		fmt.Println(Warn("Item not found!"))
		emptyItem := new(Spell)
		return emptyItem, false
	}
	retrievedItem := spellResult[0]
	return &retrievedItem, true
}
func getInventory(accountID uuid.UUID, mongoClient *mongo.Client) (*PlayerInventory, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	inventory := database.Collection("inventory")
	filterCursor, err := inventory.Find(cxt, bson.M{"uuid": accountID})
	if err != nil {
		fmt.Println(Failure(err))
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
func createInventory(accountID uuid.UUID, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	inventory := database.Collection("inventory")

	var freshInventory PlayerInventory
	freshInventory.ObjectID = primitive.NewObjectID()
	freshInventory.Account_id = accountID

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
		fmt.Println(Failure(err))
		return false
	}
	fmt.Println(Success("Fresh Inventory created for user! insertID: ", insertResult.InsertedID))
	return true
}
func createLoadout(accountID uuid.UUID, mongoClient *mongo.Client) bool {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	inventory := database.Collection("loadout")

	var freshLoadout PlayerLoadout
	freshLoadout.ObjectID = primitive.NewObjectID()
	freshLoadout.Account_id = accountID

	freshLoadout.Head = Item{}
	freshLoadout.Body = Item{}
	freshLoadout.Feet = Item{}
	freshLoadout.Weapon = Item{}
	freshLoadout.Accessory_1 = Item{}
	freshLoadout.Accessory_2 = Item{}
	freshLoadout.Accessory_3 = Item{}

	insertResult, err := inventory.InsertOne(cxt, freshLoadout)
	if err != nil {
		fmt.Println(Failure(err))
		return false
	}
	fmt.Println(Success("Fresh Loadout created for user! insertID: ", insertResult.InsertedID))
	return true
}
func getLoadout(accountID uuid.UUID, mongoClient *mongo.Client) (*PlayerLoadout, bool) {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	profile := database.Collection("loadout")
	filterCursor, err := profile.Find(cxt, bson.M{"uuid": accountID})
	if err != nil {
		fmt.Println(Failure(err))
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
		return emptyItem, false
	}
	fmt.Println(Info("Item retrieved : ", itemResult[0].Item_id))
	retrievedItem := itemResult[0]
	return &retrievedItem, true
}
func equipItem(accountID uuid.UUID, itemID string, mongoClient *mongo.Client) string {
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
		match := bson.M{"uuid": accountID}
		change := bson.M{"$set": bson.M{entryArea: retrievedItem}}
		updateResponse, err := loadout.UpdateOne(cxt, match, change)
		fmt.Println(Info(updateResponse))
		if err != nil {
			fmt.Println(Failure(err))
			return "EQUIP$0"
		}
		updateVital(accountID, retrievedItem, "add", mongoClient)
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
	loadout := database.Collection("loadout")
	retrievedItem, itemFound := getItem(itemID, mongoClient)
	if itemFound {
		entryArea := retrievedItem.Item_type
		match := bson.M{"uuid": accountID}
		change := bson.M{"$set": bson.M{entryArea: Item{}}}
		updateResponse, err := loadout.UpdateOne(cxt, match, change)
		fmt.Println(Info(updateResponse))
		if err != nil {
			fmt.Println(Failure(err))
			return "UNEQUIP$0"
		}
		updateVital(accountID, retrievedItem, "remove", mongoClient)
		return "UNEQUIP$1"
	}
	return "UNEQUIP$0"
}
func updateVital(accountID uuid.UUID, item *Item, operation string, mongoClient *mongo.Client) *PlayerVital {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	vital := database.Collection("vital")
	retrievedVital, _ := getVital(accountID, mongoClient)
	originalStats := retrievedVital.Stats
	updatedStats := updateStatsByItem(&originalStats, item, operation)
	retrievedVital.Stats = *updatedStats
	match := bson.M{"uuid": accountID}
	change := bson.M{"$set": bson.D{{"profile", retrievedVital.PlayerProfile}, {"stats", retrievedVital.Stats}}}
	updateResponse, err := vital.UpdateOne(cxt, match, change)
	fmt.Println(Info(updateResponse))
	if err != nil {
		fmt.Println(Failure(err))
	}
	return retrievedVital
}
func addInventoryItem(accountID uuid.UUID, itemID string, mongoClient *mongo.Client) string {
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
		match := bson.M{"uuid": accountID}
		change := bson.M{"$push": bson.M{entryArea: retrievedItem}}
		updateResponse, err := inventory.UpdateOne(cxt, match, change)
		fmt.Println(Info(updateResponse))
		if err != nil {
			fmt.Println(Failure(err))
			return "Item addition failed!"
		}
		return "Item added successfully!"
	}
	return "Item does not exist!"
}
func removeInventoryItem(accountID uuid.UUID, itemID string, mongoClient *mongo.Client) string {
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(cxt, readpref.Primary()); err != nil {
		panic(err)
	}
	database := mongoClient.Database("player")
	inventory := database.Collection("inventory")
	retrievedItem, itemFound := getItem(itemID, mongoClient)
	if itemFound {
		entityType := strings.ToLower(retrievedItem.Entity)
		entryArea := entityType + "." + retrievedItem.Item_type + ".collection"
		match := bson.M{"uuid": accountID}
		arr := []Item{*retrievedItem}
		change := bson.M{"$pull": bson.M{entryArea: bson.M{"$in": arr}}}
		updateResponse, err := inventory.UpdateOne(cxt, match, change)
		fmt.Println(Info(updateResponse))
		if err != nil {
			fmt.Println(Failure(err))
			return "Item deletion failed!"
		}
		return "Item deletion successfully!"
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
	//createPorts()
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println(Warn("Please provide port"))
		return
	}
	//mongoDB specs
	uri := "mongodb+srv://admin:london1234@cluster0.8acnf.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	mapInstance = Map{}
	mapInstance.ConnectedClients = make(map[uuid.UUID]Client)
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
	sessions.Battles = make(map[uuid.UUID]BattleSession)
	// add to sync.WaitGroup
	wg.Add(1)
	go tcpListener(PORT, cxt, mongoClient)
	wg.Add(2)
	go udpListener(":26950", cxt, mongoClient)
	wg.Wait()
}
