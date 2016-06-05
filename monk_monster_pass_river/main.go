package main

import (
	"fmt"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	init_status := Status{
		Monk:      [2]int{3, 0},
		Monster:   [2]int{3, 0},
		Direction: true,
		History:   map[string]bool{},
	}

	possible_solutions := init_status.FindSolutions()
	log.Println("find", len(possible_solutions), "solutions")
	for i, solution := range possible_solutions {
		status := Status{
			Monk:      [2]int{3, 0},
			Monster:   [2]int{3, 0},
			Direction: true,
		}
		log.Println("=========== solution", i, "===========")
		for _, pass := range solution {
			status.TakeAction(pass)
			fmt.Println(pass, "\t\t\t", status)
		}
	}
}

type BoatAction struct {
	Monk, Monster int
}
type ActionList []BoatAction
type Status struct {
	Monk      [2]int
	Monster   [2]int
	History   map[string]bool
	Direction bool
}

func (s *Status) TakeAction(b BoatAction) {
	if s.Direction {
		// from left to right
		s.Monk[0] -= b.Monk
		s.Monk[1] += b.Monk

		s.Monster[0] -= b.Monster
		s.Monster[1] += b.Monster
	} else {
		// from left to right
		s.Monk[1] -= b.Monk
		s.Monk[0] += b.Monk

		s.Monster[1] -= b.Monster
		s.Monster[0] += b.Monster
	}
	s.Direction = !s.Direction
}

func (s Status) Check() bool {
	if (s.Monk[0] > 0 && s.Monster[0] > s.Monk[0]) ||
		(s.Monk[1] > 0 && s.Monster[1] > s.Monk[1]) ||
		s.Monk[0] < 0 || s.Monk[1] < 0 ||
		s.Monster[0] < 0 || s.Monster[1] < 0 {
		return false
	}
	return true
}
func (s *Status) FindSolutions() []ActionList {
	// left to right
	key := s.String()
	ret := []ActionList{}
	if met, exists := s.History[key]; exists && met {
		log.Println("met before:", s)
		return ret
	}
	s.History[key] = true
	side := 0
	if !s.Direction {
		// right to left
		side = 1
	}
	for i := 0; i <= s.Monk[side]; i++ {
		for j := 0; j <= s.Monster[side]; j++ {
			action := BoatAction{
				Monk:    i,
				Monster: j,
			}
			if action.Check() {
				log.Println(action, "IS valid Action!")
				s.TakeAction(action)
				if s.Check() {
					log.Println(s, "IS valid Status!")
					log.Println("After Routin", action, "status:", s)
					if s.End() {
						ret = append(ret, ActionList{action})
					} else {
						more_actions := s.FindSolutions()
						if len(more_actions) > 0 {
							for idx, mas := range more_actions {
								more_actions[idx] = append(ActionList{action}, mas...)
							}
							ret = append(ret, more_actions...)
						}
					}
				} else {
					log.Println(s, "is NOT valid Status!")
				}
				// 原样运回去：
				s.TakeAction(action)
				log.Println("Revert Routin", action, "status:", s)
			} else {
				log.Println(action, "is NOT valid Action!")
			}
		}
	}
	s.History[key] = false
	return ret
}

func (s Status) End() bool {
	if s.Monk[0] == 0 &&
		s.Monk[1] == 3 &&
		s.Monster[0] == 0 &&
		s.Monster[1] == 3 &&
		!s.Direction {
		log.Println("finish job!")
		return true
	}
	return false
}

func (s Status) String() string {
	if s.Direction {
		return fmt.Sprintf("boat, %d %d | %d %d",
			s.Monk[0], s.Monster[0], s.Monk[1], s.Monster[1])
	} else {
		return fmt.Sprintf("%d %d | %d %d, boat",
			s.Monk[0], s.Monster[0], s.Monk[1], s.Monster[1])
	}
}
func (b BoatAction) Check() bool {
	// 存在monk且monster多于monk时，无效！
	if (b.Monk > 0 && b.Monster > b.Monk) ||
		(b.Monk == 0 && b.Monster == 0) ||
		(b.Monk+b.Monster > 2) {
		return false
	}
	return true
}
func (b BoatAction) String() string {
	return fmt.Sprintf("Monk[%d] Monster[%d]", b.Monk, b.Monster)
}
