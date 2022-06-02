package alghoritms

import (
	"fmt"
	"math"
)

type Bandit struct {
	ID     int
	Trials int
	Reward int
}

func ChooseAlgorithm(bandit []Bandit, allTrials int) (int, error) {
	if len(bandit) == 0 {
		return 0, fmt.Errorf("len slice is nil")
	}

	index := 0
	var statistic float64 = -1
	for i, v := range bandit {
		reward := float64(v.Reward) / float64(v.Trials)
		currentStatistic := reward + math.Sqrt(2*math.Log(float64(allTrials))/float64(v.Trials))
		if currentStatistic > statistic {
			index, statistic = i, currentStatistic
		}
	}

	return bandit[index].ID, nil
}
