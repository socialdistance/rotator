package alghoritms

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMakeChoice(t *testing.T) {
	stats := []Bandit{
		{ID: 1, Trials: 100, Reward: 1},
		{ID: 2, Trials: 100, Reward: 1},
		{ID: 3, Trials: 100, Reward: 1},
		{ID: 9, Trials: 100, Reward: 1},
	}

	t.Run("any id", func(t *testing.T) {
		choice, err := ChooseAlgorithm(stats, 400)
		require.NoError(t, err)
		require.Contains(t, []int{1, 2, 3, 9}, choice)
	})

	// id = 1 has a lot of reward
	stats[0] = Bandit{ID: 1, Trials: 100, Reward: 99}
	t.Run("id = 1", func(t *testing.T) {
		choice, err := ChooseAlgorithm(stats, 400)
		require.NoError(t, err)
		require.Equal(t, 1, choice)
	})

	// id = 9 is largest
	stats[3] = Bandit{ID: 9, Trials: 100, Reward: 100}
	t.Run("id = 9", func(t *testing.T) {
		choice, err := ChooseAlgorithm(stats, 400)
		require.NoError(t, err)
		require.Equal(t, 9, choice)
	})
}
