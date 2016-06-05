package main

import (
	//"flag"
	"fmt"
	"log"
	"math"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	b := Bucket{
		Capacity: [3]int{8, 5, 3},
		Contain:  [3]int{8, 0, 0},
		History:  map[string]bool{},
	}

	result := b.Resolve()
	log.Println(">>>>>> find", len(result), "solutions:")
	for _, r := range result {
		b := Bucket{
			Capacity: [3]int{8, 5, 3},
			Contain:  [3]int{8, 0, 0},
		}
		log.Printf("================ solution: %d =============\n", len(r))
		for _, a := range r {
			b.TakeAction(&a)
			fmt.Println(a, b)
		}
	}
}

type Bucket struct {
	Capacity,
	Contain [3]int
	History map[string]bool
}
type Action struct {
	Amount,
	From,
	To int
}

func (a Action) String() string {
	return fmt.Sprintf("[Amout: %d, From: %d, To: %d]",
		a.Amount, a.From, a.To)
}

type ActionList []Action

func (b *Bucket) Resolve() []ActionList {
	ret := []ActionList{}
	var sign string = b.String()
	if processed, exists := b.History[sign]; exists && processed {
		log.Println("jump duplate case:", b)
		return ret
	} else {
		b.History[sign] = true
	}
	log.Println("enter next layer:", b)
	valid_actions := b.GetActions(func(a *Action) bool {
		b.TakeAction(a)
		defer b.RevertAction(*a)
		if b.End() {
			return true
		}
		a.Amount = 0
		return false
	})
	for _, a := range valid_actions {
		if a.Amount > 0 {
			ret = append(ret, ActionList{a})
			continue
		}
		// 这些都是需要进一步试探的action
		b.TakeAction(&a)
		if a.Amount <= 0 || b.Contain[a.From] < 0 || b.Contain[a.To] > b.Capacity[a.To] {
			log.Fatalln("invalid action!", a, b)
		}
		fmt.Printf("dump %d from %d to %d, result: \t%v\n",
			a.Amount, a.From, a.To, b)
		next_actions := b.Resolve()
		if len(next_actions) > 0 {
			for _, na := range next_actions {
				ret = append(ret, append(
					ActionList{a}, na...))
			}
		}
		b.RevertAction(a)
	}
	b.History[sign] = false
	return ret
}
func (b Bucket) String() string {
	return fmt.Sprintf("%v", b.Contain)
}
func (b Bucket) GetActions(finished func(a *Action) bool) []Action {
	const l = 3
	ret := []Action{}
	for i := 0; i < l; i++ {
		for j := 0; j < l; j++ {
			if i == j {
				continue
			}
			action := Action{0, i, j}
			if finished(&action) ||
				b.validAction(action) {
				ret = append(ret, action)
			}
		}
	}
	log.Println("valid actions:", ret)
	return ret
}

func (b *Bucket) TakeAction(action *Action) bool {
	if b.validAction(*action) {
		amount := int(math.Min(
			float64(b.Contain[action.From]),
			float64(b.Capacity[action.To]-b.Contain[action.To])))

		b.Contain[action.From] -= amount
		b.Contain[action.To] += amount
		action.Amount = amount

		log.Println("after dump water:", action, b)
		return true
	}
	return false
}
func (b *Bucket) RevertAction(action Action) {
	if action.Amount > 0 {
		b.Contain[action.From] += action.Amount
		b.Contain[action.To] -= action.Amount
		log.Println("revert bucket:", b)
	}
}
func (b Bucket) End() bool {
	return b.Contain[0] == 4 && b.Contain[1] == 4
}
func (b Bucket) validAction(a Action) bool {
	return (b.Contain[a.From] > 0 &&
		(b.Capacity[a.To]-b.Contain[a.To] > 0))
}
