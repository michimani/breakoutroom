package main

import (
	"fmt"
	"math"
)

const roomCount = 5

type Room struct {
	Name    string
	Members []Member
}

type Rooms [roomCount]Room

var members Members

func init() {
	if err := LoadMembers(&members); err != nil {
		panic(err)
	}
}

const seed = 0

func main() {
	allPattern := []Rooms{}
	offset := 1

	// シャッフル回数は部屋の数
	for s := 0; s < roomCount; s++ {
		down := false
		memberIdx := 0
		rooms := Rooms{}

		// 主役から先に部屋に配置
		// 何回目のシャッフルかによって開始位置をずらす
		// for _, l := range members.LeadingParts[s*offset:] {
		for _, l := range members.LeadingParts {
			roomNo := (memberIdx + seed) % roomCount
			if down {
				roomNo = int(math.Abs(float64(((memberIdx + seed) % roomCount) - roomCount + 1)))
			}
			rooms[roomNo].Members = append(rooms[roomNo].Members, l)
			memberIdx++
			if roomNo == 0 {
				down = false
			} else if roomNo == (roomCount - 1) {
				down = true
			}
		}

		// ずらした分を配置
		// if s > 0 {
		// 	for _, mm := range members.LeadingParts[:s*offset] {
		// 		roomNo := (memberIdx + seed) % roomCount
		// 		if down {
		// 			roomNo = int(math.Abs(float64(((memberIdx + seed) % roomCount) - roomCount + 1)))
		// 		}
		// 		rooms[roomNo].Members = append(rooms[roomNo].Members, mm)
		// 		memberIdx++
		// 		if roomNo == 0 {
		// 			down = false
		// 		} else if roomNo == (roomCount - 1) {
		// 			down = true
		// 		}
		// 	}
		// }

		// 参加者を部屋に配置
		// 何回目のシャッフルかによって開始位置をずらす
		memberIdx = 0
		down = false
		for _, m := range members.Participants[s*offset:] {
			roomNo := (memberIdx + seed) % roomCount
			if down {
				roomNo = int(math.Abs(float64(((memberIdx + seed) % roomCount) - roomCount + 1)))
			}
			rooms[roomNo].Members = append(rooms[roomNo].Members, m)
			memberIdx++
			if roomNo == 0 {
				down = false
			} else if roomNo == (roomCount - 1) {
				down = true
			}
		}

		// ずらした分を配置
		if s > 0 {
			for _, mm := range members.Participants[:s*offset] {
				roomNo := (memberIdx + seed) % roomCount
				if down {
					roomNo = int(math.Abs(float64(((memberIdx + seed) % roomCount) - roomCount + 1)))
				}
				rooms[roomNo].Members = append(rooms[roomNo].Members, mm)
				memberIdx++
				if roomNo == 0 {
					down = false
				} else if roomNo == (roomCount - 1) {
					down = true
				}
			}
		}

		offset++
		allPattern = append(allPattern, rooms)
	}

	membersMatch := map[Member]map[Member]struct{}{}

	for n, p := range allPattern {
		fmt.Printf("------ %d 回目 ------ \n", n+1)
		for ri, r := range p {
			fmt.Printf("\tルーム %d\n", ri+1)
			for _, m := range r.Members {
				fmt.Printf("\t\t%s\n", m)
				for _, mm := range r.Members {
					if m != mm {
						if membersMatch[m] == nil {
							membersMatch[m] = map[Member]struct{}{}
						}
						membersMatch[m][mm] = struct{}{}
					}
				}
			}
		}
	}

	for _, m := range members.All {
		fmt.Printf("%s は %d 人と喋りました\n", m, len(membersMatch[m]))
	}
}
