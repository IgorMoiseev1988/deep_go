package main

import (
	"math"
	"testing"
	"unsafe"
//	"fmt"
	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

const (
	NameLenBitMask     uint16 = 0xFC00
	NameLenBitShift    uint16 = 0x000A
	ManaBitMask        uint16 = 0x03FF
	HealthBitMask      uint16 = 0x03FF
	RespectBitMask     uint16 = 0xF000
	RespectBitShift    uint16 = 0x000C
	StrengthBitMask    uint16 = 0x0F00
	StrengthBitShift   uint16 = 0x0008
	ExperienceBitMask  uint16 = 0x00F0
	ExperienceBitShift uint16 = 0x0004
	LevelBitMask       uint16 = 0x000F
	HouseBitMask       uint16 = 0x4000
	GunBitMask         uint16 = 0x2000
	FamilyBitMask      uint16 = 0x1000
	PersonTypeBitMask  uint16 = 0x0C00
	PersonTypeBitShift uint16 = 0x000A
	PersonTypeBitSize  uint16 = 0x0003
)

/*****************************************\
* Option setters                          *
\*****************************************/
func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		nameLen := uint16(len(name))
		for i := uint16(0); i < nameLen; i++ {
			person.name[i] = name[i]
		}
		person.healthAndLen &= (^NameLenBitMask) /* clear nameLen */
		person.healthAndLen |= uint16(nameLen << NameLenBitShift)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)	
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaAndFlags &= (^ManaBitMask) /* clear old value */
		person.manaAndFlags |= uint16(mana) & ManaBitMask
	}
}

func WithHealth(health int) func (*GamePerson) {
	return func(person *GamePerson) {
		person.healthAndLen &= (^HealthBitMask) /* clear old value */
		person.healthAndLen |= uint16(health) & HealthBitMask
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rspStrExpLvl &= (^RespectBitMask) /* clear old value */
		person.rspStrExpLvl |= uint16(uint8(respect)) << RespectBitShift
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rspStrExpLvl &= (^StrengthBitMask) /* clear old value */
		person.rspStrExpLvl |= uint16(uint8(strength)) << StrengthBitShift
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rspStrExpLvl &= (^ExperienceBitMask) /* clear old value */
		person.rspStrExpLvl |= uint16(uint8(experience)) << ExperienceBitShift
	}
}

func WithLevel(level int) func (*GamePerson) {
	return func(person *GamePerson) {
		person.rspStrExpLvl &= (^LevelBitMask) /* clear old value */
		person.rspStrExpLvl |= uint16(uint8(level))
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaAndFlags |= HouseBitMask
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaAndFlags |= GunBitMask
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaAndFlags |= FamilyBitMask
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaAndFlags &= (^PersonTypeBitMask) /* clear old value */
		person.manaAndFlags |= (uint16(personType) & PersonTypeBitSize) << PersonTypeBitShift
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)


/*****************************************\
* GamePerson has size equal 64 bytes:     *
*       0 1 2 3 4 5 6 7 8 9 A B C D E F   *
*      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+  *
* 0x00 |   X   |   Y   |   Z   | Gold  |  *
*      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+  *
* 0x10 |                               |  *
*      +        Name (42 bytes)        +  *
* 0x20 |                               |  *
*      +                   +-+-+-+-+-+-+  *
* 0x30 |                   |   Other   |  *
*      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+  *
*                                         *
* Other bits map:                         *
*             0              1            *
*      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+  *
*  0x0 | |H|G|F| T |       Mana        |  *
*      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+  *
*  0x2 | len(name) |      Health       |  *
*      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+  *
*  0x4 |Respect|  Str  |  Exp  | Level |  *
*      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+  *     
*                                         *
* H - HasHouse                            *
* G - HasGun                              *
* F - HasFamily                           *
* T - PersonType                          *
\*****************************************/

type GamePerson struct {
	x, y, z int32
	gold uint32
	name [42]byte
	manaAndFlags uint16
	healthAndLen uint16
	rspStrExpLvl uint16
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	
	for _, option := range options {
		option(&person)
	}
	return person
}

/*****************************************\
* Getters                                 *
\*****************************************/
func (p *GamePerson) Name() string {
	return unsafe.String(unsafe.SliceData(p.name[:]), p.healthAndLen >> NameLenBitShift)
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.manaAndFlags & ManaBitMask)
}

func (p *GamePerson) Health() int {
	return int(p.healthAndLen & HealthBitMask)
}

func (p *GamePerson) Respect() int {
	return int((p.rspStrExpLvl & RespectBitMask) >> RespectBitShift)
}

func (p *GamePerson) Strength() int {
	return int((p.rspStrExpLvl & StrengthBitMask) >> StrengthBitShift)
}

func (p *GamePerson) Experience() int {
	return int((p.rspStrExpLvl & ExperienceBitMask) >> ExperienceBitShift)
}

func (p *GamePerson) Level() int {
	return int(p.rspStrExpLvl & LevelBitMask)
}

func (p *GamePerson) HasHouse() bool {
	return (p.manaAndFlags & HouseBitMask) > 0
}
	
func (p *GamePerson) HasGun() bool {
	return (p.manaAndFlags & GunBitMask) > 0
}

func (p *GamePerson) HasFamily() bool {
	return (p.manaAndFlags & FamilyBitMask) > 0
}

func (p *GamePerson) Type() int {
	return int((p.manaAndFlags & PersonTypeBitMask) >> PersonTypeBitShift)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))
	
	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = WarriorGamePersonType
	const gold = math.MaxInt32 - 1
	const mana = 1000
	const health = 999
	const respect = 9
	const strength = 8
	const experience = 7
	const level = 6

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
		
	
