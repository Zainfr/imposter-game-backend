package game

import (
	"errors"
	"math/rand"
)

type Game struct {
	Players	map[string]*Player
	SecretWord string
	CurrentPhase Phase
	CurrentRound int
	TurnOrder	[]string
	CurrentIdx 	int
	Votes 	map[string]string
	
}

type Player struct{
	ID		string
	Name	string
	IsImposter bool
	HasVoted bool
	Clues	[]string
}

type Phase string 
const (
    PhaseLobby            Phase = "lobby"
    PhaseWordDistribution Phase = "word_distribution"
    PhaseClueRound        Phase = "clue_round"
    PhaseDiscussion       Phase = "discussion"
    PhaseVoting           Phase = "voting"
)

func (g *Game) StartGame(words []string) error {
    if g.CurrentPhase != PhaseLobby {
        return errors.New("game is not in lobby phase")
    }

    if len(g.Players) < 3 {
        return errors.New("at least 3 players are required to start")
    }

    g.SecretWord = words[rand.Intn(len(words))]

    ids := make([]string, 0, len(g.Players))
    for id := range g.Players {
        ids = append(ids, id)
    }
    imposterID := ids[rand.Intn(len(ids))]
    g.Players[imposterID].IsImposter = true

    g.CurrentPhase = PhaseWordDistribution
    return nil
}

func (g *Game) WordForPlayer(playerID string) string {
    player, ok := g.Players[playerID]
    if !ok {
        return ""
    }
    if player.IsImposter {
        return ""
    }
    return g.SecretWord
}

var WordList = []string{
	"Ocean",
	"Galaxy",
	"Mountain",
	"Library",
	"Desert",
	"Volcano",
	"Orchestra",
	"Rainforest",
	"Architecture",
	"Astronaut",
}