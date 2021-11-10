package main

import (
	"encoding/json"
	"os"
)

const source = "members.json"

type Member string

type Members struct {
	LeadingParts []Member `json:"leadingParts"`
	Participants []Member `json:"participants"`

	All []Member `json:"-"`
}

func (ms Members) LeadingPartsCount() int {
	return len(ms.LeadingParts)
}

func (ms Members) ParticipantsCount() int {
	return len(ms.Participants)
}

func (ms Members) AllMembersCount() int {
	return len(ms.All)
}

func LoadMembers(ms *Members) error {
	f, err := os.Open(source)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(ms); err != nil {
		return err
	}

	ms.All = append(ms.All, ms.LeadingParts...)
	ms.All = append(ms.All, ms.Participants...)

	return nil
}
