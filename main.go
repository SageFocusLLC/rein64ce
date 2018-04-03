package main

import (
	"fmt"
	"os"

	"gonum.org/v1/gonum/mat"

	"github.com/SageFocusLLC/gophernet"
)

//has to be mupen64plus 64 bit linux, with default input plugin
func main() {
	args := os.Args[1:]
	mapPath := "./statemap.json"
	env := NewEnvironment(args, mapPath)

	env.Init()

	episodeLength := 1200 
	episodeProgress := 0.0

	//create neural net
	nnConfig := gonet.NeuralNetConfig{
		InputNeurons:  6,
		OutputNeurons: 8,
		HiddenNeurons: 10,
		NumEpochs:     1,
		LearningRate:  0.01,
	}

	agent := NewAgent(nnConfig)

	//agent.LoadNN()

	agent.SetTau(0.99)

	stateArr := []float64{0.01, 0.01, 0.01, 0.01, 0.01, 0.01}
	env.GetState(stateArr)

	reward := 0.0
	epoch := 0
	action := uint64(0x00)
	//actionP := uint64(0x00)
	mapPositionVec(stateArr[0:3])

	stateP := mat.NewDense(1, len(stateArr), nil)
	state := mat.NewDense(1, len(stateArr), nil)

	//actionMem := mat.NewDense(1000, 1, nil)
	//rewardMem := mat.NewDense(1000, 1, nil)

	

	for epoch < 100 {
		fmt.Println("Epoch ", epoch)

		env.LoadGame()

		endstate := false
		step := 1

		for step < episodeLength && endstate != true {
			episodeProgress = float64(step) / float64(episodeLength + 1)
			fmt.Println(stateArr)
			//action
			state.SetRow(0, stateArr)

			//greedy
			//action = agent.GetActionGreedy(state)

			//e greedy exploration
			//action = agent.GetActionEGreedy(state)

			//boltzmann
			action = agent.GetActionBoltzmann(state)

			//_ = env.GameStepTrain()
			//actionP = env.GameStepTrain()
			env.GameStep(action)

			//observation
			env.GetState(stateArr)

			mapPositionVec(stateArr[0:3])

			stateArr[3] = stateArr[3] / 65536.0
			stateArr[4] = stateArr[4] * 0.01

			stateArr[5] = episodeProgress

			//reward
			reward, endstate = getReward(stateArr, epoch, step)

			stateP.SetRow(0, stateArr)

			//scale reward
			reward = reward * 0.05

			agent.GiveReward(state, stateP, reward)
			
			step += 1
		}
		//decrease temp
		curTemp := agent.GetTau()
		if curTemp > 0.002 {
			agent.SetTau(curTemp * 0.9)
		}
		
		fmt.Println(agent.Q)
		fmt.Println(curTemp)
		agent.SaveNN()

		epoch += 1
	}
	fmt.Println("done :D")
}

func getReward(stateArr []float64, epoch int, step int) (float64, bool) {
	reward := 0.0
	endstate := false
/*
//stairs
	//left the room
	if stateArr[2] > 0.7{
		reward = -1.0
		fmt.Println("left the room")
		endstate = true
	}	
	//made it
	if stateArr[1] > 0.8{
		reward = 2.0
		fmt.Println("OMG")
		endstate = true
	}
	//more height -> more reward
	reward += stateArr[1] - 1.0
*/
//slide

	reward = -0.3
	if stateArr[1] > 0.27055 && stateArr[1] < 0.28270 {
		if stateArr[0] > 0.12314 && stateArr[0] < 0.24100 {
			if stateArr[2] < 0.82216 && stateArr[2] > 0.70840 {
				reward = 5.0
				fmt.Println("WIN")
				fmt.Println(stateArr)
				endstate = true
			}
		}
	}

	//fell off the level
	if stateArr[1] < 0.26 {
		reward = -2.0
		fmt.Println("Fellllllllllllllllllllllllllll")
		endstate = true
	}

	//reward = stateArr[1]*-0.8 + stateArr[3]*0.2
	//avoid start area
	if stateArr[0] > 0.64 && stateArr[1] > 0.76{
		reward += -0.5
	} else {
		reward += -0.05
	}
	//reward = reward * 0.3 - 0.2

	//slight reward for moving down
	reward += 0.01 * (1.0 - stateArr[1])

	return reward, endstate

}

func mapPositionVec(fArr []float64) {
	for i, v := range fArr {
		fArr[i] = mapPositionVal(v)
	}
}

func mapPositionVal(x float64) float64 {
	return (x / 20000.0) + 0.5
}
