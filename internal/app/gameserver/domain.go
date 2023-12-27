package gameserver

import (
	"net"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
type ShopItem struct {
	Item_uuid uuid.UUID `json:"uuid" bson:"uuid,omitempty"`
	Item      Item      `json:"shop_item" bson:"shop_item"`
	Price     float64   `json:"price" bson:"price"`
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
