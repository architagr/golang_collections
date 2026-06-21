package main

import (
	"fmt"
)

func getAvatarUrl(characterId int) string {
	// Simulate getting URL
	return fmt.Sprintf("https://base-url/%d", characterId)
}

type CharacterInfo struct {
	id            int
	characterName string
	avatarUrl     string
}

func main() {
	charactersOfGOT := []CharacterInfo{
		{id: 1, characterName: "Jon Snow"},
		{id: 2, characterName: "Daenerys Targaryen"},
		{id: 3, characterName: "Arya Stark"},
		{id: 4, characterName: "Tyrion Lannister"},
	}

	for _, charInfo := range charactersOfGOT {
		charInfo.avatarUrl = getAvatarUrl(charInfo.id)
	}

	for _, charInfo := range charactersOfGOT {
		fmt.Printf("Id: %d Name: %s, avatar: %s\n", charInfo.id, charInfo.characterName, charInfo.avatarUrl)
	}
}
