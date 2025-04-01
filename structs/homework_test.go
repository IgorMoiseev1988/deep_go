package main

import (
	"math"
	"testing"
	"unsafe"
//	"fmt"
	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

/*****************************************\
* Option setters                          *
\*****************************************/
func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		var nameLen uint16
		for i := 0; i < len(name); i++ {
			nameLen++
			person.name[i] = name[i]
		}
		person.healthAndLen &= 0x03FF /* clear nameLen */
		person.healthAndLen |= uint16(nameLen << 10)
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
		person.manaAndFlags &= 0xFC00 /* clear old value */
		person.manaAndFlags |= uint16(mana & 0x03FF)
	}
}

func WithHealth(health int) func (*GamePerson) {
	return func(person *GamePerson) {
		person.healthAndLen &= 0xFC00 /* clear old value */
		person.healthAndLen |= uint16(health & 0x03FF)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rspStrExpLvl &= 0x0FFF /* clear old value */
		person.rspStrExpLvl |= uint16(respect & 0x000F) << 12
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rspStrExpLvl &= 0xF0FF /* clear old value */
		person.rspStrExpLvl |= uint16(strength & 0x000F) << 8
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rspStrExpLvl &= 0xFF0F /* clear old value */
		person.rspStrExpLvl |= uint16(experience & 0x000F) << 4
	}
}

func WithLevel(level int) func (*GamePerson) {
	return func(person *GamePerson) {
		person.rspStrExpLvl &= 0xFFF0 /* clear old value */
		person.rspStrExpLvl |= uint16(level & 0x000F)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaAndFlags |= 0x4000
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaAndFlags |= 0x2000
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaAndFlags |= 0x1000
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaAndFlags &= 0xF3FF /* clear old value */
		person.manaAndFlags |= uint16(personType & 0x0003) << 10
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
	return unsafe.String(unsafe.SliceData(p.name[:]), p.healthAndLen >> 10)
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
	return int(p.manaAndFlags & 0x03FF)
}

func (p *GamePerson) Health() int {
	return int(p.healthAndLen & 0x03FF)
}

func (p *GamePerson) Respect() int {
	return int((p.rspStrExpLvl & 0xF000) >> 12)
}

func (p *GamePerson) Strength() int {
	return int((p.rspStrExpLvl & 0x0F00) >> 8)
}

func (p *GamePerson) Experience() int {
	return int((p.rspStrExpLvl & 0x00F0) >> 4)
}

func (p *GamePerson) Level() int {
	return int(p.rspStrExpLvl & 0x000F)
}

func (p *GamePerson) HasHouse() bool {
	return (p.manaAndFlags & 0x4000) > 0
}
	
func (p *GamePerson) HasGun() bool {
	return (p.manaAndFlags & 0x2000) > 0
}

func (p *GamePerson) HasFamily() bool {
	return (p.manaAndFlags & 0x1000) > 0
}

func (p *GamePerson) Type() int {
	return int((p.manaAndFlags & 0x0C00) >> 10)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))
	
	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = WarriorGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

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
		
	
