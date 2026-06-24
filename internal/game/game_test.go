package game

import (
	"testing"
)

func TestStartGame(t *testing.T) {
	tests := []struct {
		name		string
		players 	map[string]*Player
		phase		Phase
		words 		[]string
		expectErr	bool
		expectPhase	Phase
	}{
		{
			name: "Fewer that 3 players fail",
			players: map[string]*Player{"p1": {ID: "p1"}, "p2": {ID: "p2"}},
			phase:	PhaseLobby,
			words:	[]string{"apple"},
			expectErr: true,
			expectPhase: PhaseLobby,
		},
		{
			name: "Non-Lobby phase fail",
			players: map[string]*Player{"p1": {ID: "p1"}, "p2": {ID: "p2"}, "p3": {ID: "p3"}},
			phase: PhaseVoting,
			words: []string{"apple"},
			expectErr: true,
			expectPhase: PhaseVoting,
		},
		{
			name: "Success with 3 players",
			players: map[string]*Player{"p1": {ID: "p1"}, "p2": {ID: "p2"}, "p3": {ID: "p3"}},
			phase: PhaseLobby,
			words: []string{"apple"},
			expectErr: false,
			expectPhase: PhaseWordDistribution,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g:= &Game{Players: tt.players, CurrentPhase: tt.phase}
			err := g.StartGame(tt.words)

			if(err != nil) != tt.expectErr {
				t.Errorf("func StartGame() error: %v, expectErr: %v", err, tt.expectErr)
			}
			if g.CurrentPhase != tt.expectPhase {
				t.Errorf("Phase = %v, wanted = %v", g.CurrentPhase, tt.expectPhase)
			}
		})
	}
}

func TestWordForPlayer(t *testing.T) {
	g := &Game{
		SecretWord: "ocean",
		Players: map[string]*Player{
			"human": {ID: "human",IsImposter: false},
			"imposter":{ ID: "imposter", IsImposter: true},
		},
	}
	t.Run("Imposter gets an empty string", func(t *testing.T) {
		got := g.WordForPlayer("imposter")
		if got != "" {
			t.Errorf("Imposter should get empty string, got :%s", got)
		}
	})

	t.Run("normal human gets the secret word", func(t *testing.T) {
		got := g.WordForPlayer("human")
		if got != "ocean" {
			t.Errorf("Human shoud get the secrectword ocean, got %s", got)
		}
	})
}