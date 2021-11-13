package main

import (
	"fmt"
	"math"
	"sort"
)

// 部屋の数は主役の数 + 1 以上にする
const roomCount = 6

const RoopLimit = 10000

type Room struct {
	Name    string
	Members []Member
}

type Rooms [roomCount]Room

type MemberMatchedScore map[Member]map[Member]int

type RemainScore struct {
	RoomNo int
	Score  int
}

func NewMemberMatchedScore() MemberMatchedScore {
	membersMatchScore := MemberMatchedScore{}
	for _, m := range members.All {
		membersMatchScore[m] = map[Member]int{}
	}
	return membersMatchScore
}

func (mms MemberMatchedScore) Score(source, target Member) int {
	if mms == nil {
		return 0
	}
	if _, ok := mms[source]; !ok {
		return 0
	}

	v, ok := mms[source][target]
	if !ok {
		return 0
	}

	return v
}

func (mms MemberMatchedScore) Up(m1, m2 Member) MemberMatchedScore {
	if mms == nil {
		return nil
	}

	if m1 == m2 {
		return mms
	}

	if _, ok := mms[m1]; !ok {
		mms[m1] = map[Member]int{}
	}
	if _, ok := mms[m2]; !ok {
		mms[m2] = map[Member]int{}
	}

	v, ok := mms[m1][m2]
	if !ok {
		mms[m1][m2] = 1
	} else {
		mms[m1][m2] = v + 1
	}

	v2, ok2 := mms[m2][m1]
	if !ok2 {
		mms[m2][m1] = 1
	} else {
		mms[m2][m1] = v2 + 1
	}

	return mms
}

var members Members
var maxRoomMember int

func init() {
	if err := LoadMembers(&members); err != nil {
		panic(err)
	}

	// 部屋の数 確認 (主役 + 1)
	if members.LeadingPartsCount()+1 > roomCount {
		panic(fmt.Errorf("部屋の数が `%d` で定義されています。部屋の数は `主役の数 + 1` = `%d` 以上にしてください。\n", roomCount, members.LeadingPartsCount()+1))
	}

	maxRoomMember = int(math.Ceil(float64(members.AllMembersCount()) / float64(roomCount)))
}

type Expect struct {
	Member Member
	Count  int
}

func main() {
	expect := &Expect{ // 特定の参加者に期待するマッチ人数
		Member: Member("主役1"),
		Count:  3,
	}

	found := false
	for i := 0; i < RoopLimit; i++ {
		matchRes, allPattern := generate()
		if expect == nil {
			showResult(allPattern)
			found = true
			break
		}
		if len(matchRes[expect.Member]) == expect.Count {
			showResult(allPattern)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("期待値を満たすパターンは見つかりませんでした。再実行すると見つかるかもしれません。")
	}
}

func generate() (MemberMatchedScore, []Rooms) {
	allPattern := []Rooms{}
	membersMatchScore := NewMemberMatchedScore()

	// シャッフル回数は部屋の数
	for s := 0; s < roomCount; s++ {
		rooms := Rooms{}
		participant := members.Participants

		// 各部屋の最初の一人を決める
		for i, l := range members.LeadingParts {
			rooms[i].Members = append(rooms[i].Members, l)
		}

		// 参加者の最初の一人はシャッフルごとに変える
		rooms[roomCount-1].Members = append(rooms[roomCount-1].Members, participant[s])
		currParticipant := remove(participant, s)
		participantMap := sliceToMap(currParticipant)

		// 参加者を各部屋に割り当てていく
		remains := map[Member][]RemainScore{} // 余ったやつは後で

		for p := range participantMap {
			added := false
			// 各部屋に対して、既に割り当てられた人とのマッチスコアが 0 なら割り当て
			// 0 でなければ部屋番号とスコアを保持して、あとから再割り当て
			for roomNo, r := range rooms {

				if len(rooms[roomNo].Members) == maxRoomMember {
					continue
				}

				roomMatchScore := 0 // ルーム内のメンバーとのスコアが一番高いもの
				for rmi, rm := range r.Members {
					memberMatchScore := membersMatchScore.Score(p, rm)
					// 一人目 (=主役) とのスコアが 0 ならその時点で入れる
					if rmi == 0 && memberMatchScore == 0 {
						rooms[roomNo].Members = append(rooms[roomNo].Members, p)
						added = true
					}

					// ルームとのスコアを足していく (余ったときに使う) (主役とのスコアは足さない)
					if rmi != 0 && roomMatchScore < memberMatchScore {
						roomMatchScore = memberMatchScore
					}
				}

				// 主役とはマッチしたが他の全員とはマッチしない場合は割り当てる
				if !added && roomMatchScore == 0 {
					rooms[roomNo].Members = append(rooms[roomNo].Members, p)
					added = true
				}

				// どこかのルームに追加されていればルームメンバーとのスコアを UP して break
				if added {
					for _, rm := range rooms[roomNo].Members {
						membersMatchScore = membersMatchScore.Up(p, rm)
					}
					break
				}

				// 追加されていなければ再割り当てリストに追加
				if v, ok := remains[p]; !ok || v == nil {
					remains[p] = []RemainScore{}
				}
				remains[p] = append(remains[p], RemainScore{RoomNo: roomNo, Score: roomMatchScore})
			}

			if added {
				remains[p] = nil
			}
		}

		// どこにも割り当てられなかった参加者は、ルームスコアが低いルームに割り当てていく
		for remainMember, remain := range remains {
			if remain == nil {
				continue
			}

			added := false

			sort.Slice(remain, func(i, j int) bool { return remain[i].Score < remain[j].Score })
			for _, rs := range remain {
				if len(rooms[rs.RoomNo].Members) < maxRoomMember {
					for _, rm := range rooms[rs.RoomNo].Members {
						membersMatchScore = membersMatchScore.Up(rm, remainMember)
					}
					rooms[rs.RoomNo].Members = append(rooms[rs.RoomNo].Members, remainMember)
					added = true
					break
				}
			}

			// それでも割り当てられない場合は、仕方ないので一番スコアの低いルームに割り当てる
			if !added {
				for _, rm := range rooms[remain[0].RoomNo].Members {
					membersMatchScore = membersMatchScore.Up(rm, remainMember)
				}
				rooms[remain[0].RoomNo].Members = append(rooms[remain[0].RoomNo].Members, remainMember)
			}
		}

		allPattern = append(allPattern, rooms)
	}

	// showResult(allPattern)

	// fmt.Println("---------------------------")

	// showMemberMatchScore(membersMatchScore)

	return membersMatchScore, allPattern
}

func showResult(r []Rooms) {
	membersMatch := MemberMatchedScore{}
	for n, rooms := range r {
		fmt.Printf("------ %d 回目 ------ \n", n+1)
		all := 0
		for roomNo, room := range rooms {
			fmt.Printf("\tルーム %d [%d 人]\n", roomNo+1, len(room.Members))
			all = all + len(room.Members)
			for _, m := range room.Members {
				fmt.Printf("\t\t%s\n", m)
				for _, mm := range room.Members {
					if m != mm {
						if _, ok := membersMatch[m]; !ok {
							membersMatch[m] = map[Member]int{}
						}
						membersMatch = membersMatch.Up(m, mm)
					}
				}
			}
		}
		fmt.Printf("------ %d 回目 参加人数 %d 人------ \n", n+1, all)
	}

	for _, m := range members.All {
		fmt.Printf("%s は %d 人と喋りました\n", m, len(membersMatch[m]))
	}
}

func showMemberMatchScore(mms MemberMatchedScore) {
	for _, m := range members.All {
		fmt.Printf("%s が喋った人のスコア (一緒のルームになった回数)\n", m)
		for mm, score := range mms[m] {
			fmt.Printf("\t%3d: %s\n", score, mm)
		}
	}

	for _, m := range members.All {
		fmt.Printf("%s は %d 人と喋りました\n", m, len(mms[m]))
	}
}

func remove(s []Member, idx int) []Member {
	tmp := []Member{}
	for i := 0; i < len(s); i++ {
		tmp = append(tmp, s[i])
	}
	return append(tmp[:idx], tmp[idx+1:]...)
}

func sliceToMap(s []Member) map[Member]struct{} {
	m := map[Member]struct{}{}
	for _, sm := range s {
		m[sm] = struct{}{}
	}

	return m
}
