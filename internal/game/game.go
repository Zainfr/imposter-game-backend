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
		return errors.New("The game's current phase is not lobby.")
	}

	if len(g.Players) < 3 {
		return errors.New("Players must atleast be 3 to start the game.")
	}

	secretWord := words[rand.Intn(len(words))]
	g.SecretWord = secretWord

	ids := make([]string, 0, len(g.Players))
	for id := range g.Players {
		ids = append(ids, id)
	}

	imposterId := ids[rand.Intn(len(ids))]
	g.Players[imposterId].IsImposter = true

	g.CurrentPhase = PhaseWordDistribution
	
	return nil

}