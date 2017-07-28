package main

import (
	"fmt"
	"time"
)

const (
	Townhall = "townhall"
	Storage  = "storage"
	Houses   = "houses"
	Farm     = "farm"
	Sawmill  = "sawmill"
	Mine     = "mine"
	Barracks = "barracks"
)

type Coef struct {
	Gold  int
	Wood  int
	Stone int
}

var coefs = map[string]Coef{
	Townhall: Coef{Gold: 500, Wood: 22, Stone: 200},
	Storage:  Coef{Gold: 200, Wood: 100, Stone: 100},
	Houses:   Coef{Gold: 200, Wood: 100, Stone: 100},
	Farm:     Coef{Gold: 100, Wood: 50, Stone: 50},
	Sawmill:  Coef{Gold: 100, Wood: 50, Stone: 50},
	Mine:     Coef{Gold: 100, Wood: 50, Stone: 50},
	Barracks: Coef{Gold: 200, Wood: 100, Stone: 100},
}

type State struct {
	Buildings map[string]int
}

type Upgrade struct {
	Type     string
	Gold     int
	Stone    int
	Wood     int
	Duration time.Duration
}

func (u *Upgrade) String() string {
	return fmt.Sprintf("[%s] Type: %s, Gold: %d, Stone: %d, Wood: %d", u.Duration, u.Type, u.Gold, u.Stone, u.Wood)
}

func (u *Upgrade) totalCost() int {
	return u.Gold + u.Stone*2 + u.Wood*2
}

var initialState = &State{
	Buildings: map[string]int{
		Townhall: 130,
		Storage:  177,
		Houses:   185,
		Farm:     151,
		Sawmill:  82,
		Mine:     82,
		Barracks: 176,
	},
}

func main() {
	state := initialState
	var totalDuration time.Duration
	for state.Buildings[Barracks] < 500 {
		up := state.rushUpgrade()
		state.apply(up)
		fmt.Println(up)
		totalDuration += up.Duration
	}
	fmt.Printf("Total duration: %s (%d days)\n", totalDuration, totalDuration/(24*time.Hour))
}

func (s *State) apply(up *Upgrade) {
	s.Buildings[up.Type]++
}

func (s *State) balancedUpgrade() *Upgrade {
	return &Upgrade{}
}

func (s *State) rushUpgrade() *Upgrade {
	if s.Buildings[Barracks] < s.Buildings[Houses] {
		return s.calcUpgrade(Barracks)
	}
	if up, ok := s.storageFitUpgrade(Houses); ok {
		return up
	}
	return s.calcUpgrade(Storage)
}

func (s *State) storageFitUpgrade(upType string) (*Upgrade, bool) {
	up := s.calcUpgrade(upType)
	storageLvl := s.Buildings[Storage]
	storageCap := (storageLvl*50 + 1000) * storageLvl
	if up.Wood > storageCap || up.Stone > storageCap {
		return nil, false
	}
	return up, true
}

func (s *State) calcUpgrade(upType string) *Upgrade {
	level := s.Buildings[upType]
	k := (level + 1) * (level + 2) / 2
	up := &Upgrade{
		Type:  upType,
		Gold:  k * coefs[upType].Gold,
		Wood:  k * coefs[upType].Wood,
		Stone: k * coefs[upType].Stone,
	}
	up.Duration = time.Minute * time.Duration((up.totalCost() / s.gpm()))
	return up
}

func (s *State) gpm() int {
	gold := s.Buildings[Houses] * (10 + s.Buildings[Townhall]*2)
	food := (s.Buildings[Farm] - s.Buildings[Houses]) * 20
	wood := s.Buildings[Sawmill] * 20
	stone := s.Buildings[Mine] * 20
	return gold + food + wood + stone
}
