package dispatch

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/master/predictor"
	"taskAssignmentForEdge/taskmgt"
	"time"
)

const (
	Dispatcher_RR = iota
	Dispatcher_SCT_Greedy  //greedy, shortest complete time
	Dispatcher_SCT_Genetic
	Dispatcher_STT_Greedy  //shortest total time, greedy
	Dispatcher_STT_Genetic	//shortest total time, genetic algorithm
	Dispatcher_SCT_RC       //FemtoClouds risk-controlled
	Dispatcher_SCT_Genetic_ExtraWait //shortest complete time with extra queueing time by the same node, genetic algorithm
	Dispatcher_STT_Genetic_ExtraWait //shortest total time with extra queueing time by the same node, genetic algorithm
)

var DispatcherList = []string {"Dispatcher_RR", "Dispatcher_SCT_Greedy","Dispatcher_SCT_Genetic",
	"Dispatcher_STT_Greedy", "Dispatcher_STT_Genetic", "Dispatcher_SCT_EW_Genetic","Dispatcher_STT_EW_Genetic" }

type Dispatcher interface {
	//EnqueueTask(tq *taskmgt.TaskQueue, task *taskmgt.TaskEntity )
	//EnqueueNode(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity )

	MakeDispatchDecision(tq *taskmgt.TaskQueue, nq *nodemgt.NodeQueue) (isNeedAssign bool)
	//UpdateDispatcherStat(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity)

	//Node(nq *nodemgt.NodeQueue, node *nodemgt.NodeEntity)
	//CheckNode(nq *nodemgt.NodeQueue)
}

func NewDispatcher(index int)  Dispatcher{
	if index < len(DispatcherList) {
		log.Printf("Using %s algorithm", DispatcherList[index])
	}
	switch index {
	case Dispatcher_RR:
		return NewRRDispatcher()
	case Dispatcher_SCT_Greedy:
		return NewSctGreedyDispatcher()
	case Dispatcher_SCT_Genetic:
		return NewSctGADispatcher()
	case Dispatcher_STT_Greedy:
		return NewSttGreedyDispatcher()
	case Dispatcher_STT_Genetic:
		return NewSttGADispatcher()
	case Dispatcher_SCT_RC:
		return  NewRcDispatcher()
	case Dispatcher_SCT_Genetic_ExtraWait:
		return NewSctBetterGADispatcher()
	case Dispatcher_STT_Genetic_ExtraWait:
		return NewSttBetterGADispatcher()
	default:
		return nil
	}
	return nil
}

func AssignTaskToNode(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity) {
	task.NodeId.IP = node.NodeId.IP
	task.NodeId.Port = node.NodeId.Port
	task.Status = taskmgt.TaskStatusCode_Assigned
	log.Printf("Predict: task %d, trans:%d ms, wait:%d ms, exec:%d ms, extra:%d ms in Node(%s:%d)", task.TaskId,
		task.PredictTransTime/1e3, task.PredictWaitTime/1e3, task.PredictExecTime/1e3, task.PredictExtraTime/1e3,
		task.NodeId.IP, task.NodeId.Port)
	node.TqPrepare.EnqueueTask(task)
}

func GetPredictTransTime(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity) int64{
	if task == nil || node == nil {
		return 0
	}
	transTime := int64(task.DataSize/node.Bandwidth*1e6*8)
	return transTime
}

func GetPredictWaitTime(node *nodemgt.NodeEntity) int64 {
	if node == nil {
		return 0
	}
	return int64(float64(node.WaitQueueLen)*node.RunTimePredict.GetDefaultUnit().SinglePredict())
}

func GetPredictAvgExecTime(node *nodemgt.NodeEntity) int64 {
	if node == nil {
		return 0
	}
	return int64(node.RunTimePredict.GetDefaultUnit().SinglePredict())
}

func GetPredictExecTime(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity) int64 {
	fpairList := predictor.CreateFeaturePairs(task)
	execTime, _ := node.RunTimePredict.Predict(fpairList)
	return int64(execTime)
}

//return execute time and extra time
func GetPredictTimes(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity) (int64, int64){
	fpairList := predictor.CreateFeaturePairs(task)
	execTime, bestFec := node.RunTimePredict.Predict(fpairList)
	execTimeHisto := node.RunTimePredict.FindHistogByFec(bestFec)
	if execTimeHisto == nil {  //原来一个任务都没有, 否则是FeatureVal_All对应的直方图
		return int64(execTime), 0
	} else {
		lastTimeForNode := time.Now().UnixNano()/1e3 - node.JoinTST
		fracNode := node.ConnPredict.GetHistogram().CalFractionFromOffset(float64(lastTimeForNode))
		preTimeForTask := GetPredictTransTime(task, node) + GetPredictWaitTime(node)
		fracTask := node.RunTimePredict.FindHistogByFec(bestFec).CalFractionAfterShifting(float64(preTimeForTask))
		_, extraTime := fracNode.JointProbaByEqualOrLess(fracTask)
		return  int64(execTime), int64(extraTime)
	}
}

func GetUsualPredictTimes(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity) (transTime int64, waitTime int64, execTime int64){
	transTime = GetPredictTransTime(task, node)
	waitTime = GetPredictWaitTime(node)
	execTime = GetPredictExecTime(task, node)
	return
}

func GetAllPredictTimes(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity) (transTime int64, waitTime int64, execTime int64, extraTime int64) {
	transTime = GetPredictTransTime(task, node)
	waitTime = GetPredictWaitTime(node)
	execTime, extraTime = GetPredictTimes(task, node)
	return
}

//add up extra time when isAddExtra is true
func GetTotalTime(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity, isAddExtra bool )  (totalTime int64){
	transTime := GetPredictTransTime(task, node)
	waitTime := GetPredictWaitTime(node)
	if isAddExtra {
		execTime, extraTime := GetPredictTimes(task, node)
		totalTime = transTime + waitTime + execTime + extraTime
	} else {
		execTime := GetPredictExecTime(task, node)
		totalTime = transTime + waitTime + execTime
	}
	return
}

const (
	IteratorNum = 500
	CpRate = 0.3
	ChromosomeNum = 20
	ContinuousDiffCnt = 30
	FitnessDelta = 1e-10
)

//return the best chromosome
//first select and then cross&mutation
func GaAlgorithm(tasks []*taskmgt.TaskEntity, nodes []*nodemgt.NodeEntity, iteratorNum int, chromosomeNum int, isAddExtra bool) []int {
	generation := createInitGeneration( len(tasks), len(nodes), chromosomeNum)
	/*
	fmt.Printf("First Generation:\n")
	PrintGeneration(generation)*/
	fitnessList := calFitnessList(generation, tasks, nodes, isAddExtra)
	bestChromo, bestFitness := getMaxChromosome(generation, fitnessList)
	lastBestFitness := bestFitness
	slightDiffcnt := 0
	//fitnesshistory := make([]float64, IteratorNum)

	rand.Seed(time.Now().UnixNano())

	iterCnt := 0
	for ; iterCnt < iteratorNum ; iterCnt++ {

		//select
		//fitnessList = calFitness(generation, tasks, nodes, isAddExtra)
		chosenGeneration := selectChromosome(generation, fitnessList)

		//create new generation
		fitnessList = calFitnessList(chosenGeneration, tasks, nodes, isAddExtra)

		//first, copy part of max N old generation to be new generation
		copyNum := (int)(CpRate*float64(chromosomeNum))
		generation = getMaxNChromosome(chosenGeneration, fitnessList, copyNum)

		//second, create the rest part of new generation by cross and mutation
		childGeneration := crossAndMutation(chosenGeneration, fitnessList, chromosomeNum - copyNum, len(tasks), len(nodes))
		generation = append(generation, childGeneration...)

		fitnessList = calFitnessList(generation, tasks, nodes, isAddExtra)
		curBestchromo, curBestFitness := getMaxChromosome(generation, fitnessList)

		//check
		/*
		fmt.Printf("Generation %d:\n", i)
		PrintGeneration(generation)
		checkFitness, checkTime := TestChromosome(curBestchromo, tasks, nodes, isAddExtra)
		fmt.Printf("current best chromosome, totalTime %d, curfitness %.10f :", checkTime, checkFitness)
		PrintChromosome(curBestchromo)*/
		/*
		for cnt := 0; cnt < len(curBestchromo); cnt ++ {
			task := tasks[cnt]
			node := nodes[curBestchromo[cnt]]
			predTransTime, predWaitTime, predExecTime, predExtraTime := GetAllPredictTimes(task, node)
			fmt.Printf("trans:%d, wait:%d, exec:%d, extra:%d, size:%f, bw:%f\n", predTransTime, predWaitTime, predExecTime, predExtraTime, task.DataSize, node.Bandwidth)
		}
		*/
		/*
		fmt.Printf("last best chromosome, lastBestFitness %.10f :", bestFitness)
		PrintChromosome(bestChromo)*/
		/*
		for cnt := 0; cnt < len(curBestchromo); cnt ++ {
			task := tasks[cnt]
			node := nodes[bestChromo[cnt]]
			predTransTime, predWaitTime, predExecTime, predExtraTime := GetAllPredictTimes(task, node)
			fmt.Printf("trans:%d, wait:%d, exec:%d, extra:%d, size:%f, bw:%f\n", predTransTime, predWaitTime, predExecTime, predExtraTime, task.DataSize, node.Bandwidth)
		}
		*/
		//end check
		lastBestFitness = bestFitness
		if curBestFitness > bestFitness {
			bestFitness = curBestFitness
			bestChromo = curBestchromo
		}

		if bestFitness - lastBestFitness < FitnessDelta {
			slightDiffcnt ++
			if slightDiffcnt >= ContinuousDiffCnt {
				break
			}
		} else {
			slightDiffcnt = 0
		}

		//fitnesshistory[iterCnt] = bestFitness
	}

	/*
	fmt.Printf("trainning history: ")
	for j := 0;j < iterCnt-29;j++ {
		fmt.Printf("%.10f ", fitnesshistory[j])
		if j%10 == 0 {
			fmt.Printf("\n")
		}
	}*/

	//fmt.Printf("Total iteration:%d\n", iterCnt)
	return bestChromo
}

//return the best chromosome
//select max n and cross&mutation
func GaAlgorithm2(tasks []*taskmgt.TaskEntity, nodes []*nodemgt.NodeEntity, iteratorNum int, chromosomeNum int, isAddExtra bool) []int {
	generation := createInitGeneration( len(tasks), len(nodes), chromosomeNum)
	fmt.Printf("First Generation:\n")
	PrintGeneration(generation)
	fitnessList := calFitnessList(generation, tasks, nodes, isAddExtra)
	bestChromo, bestFitness := getMaxChromosome(generation, fitnessList)

	slightDiffcnt := 0

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < iteratorNum ; i++ {

		//first, copy part of max N old generation to be new generation
		copyNum := (int)(CpRate*float64(chromosomeNum))
		generation = getMaxNChromosome(generation, fitnessList, copyNum)

		//second, create the rest part of new generation by cross and mutation
		childGeneration := crossAndMutation(generation, fitnessList, chromosomeNum - copyNum, len(tasks), len(nodes))
		generation = append(generation, childGeneration...)

		fitnessList = calFitnessList(generation, tasks, nodes, isAddExtra)
		curBestchromo, curBestFitness := getMaxChromosome(generation, fitnessList)

		//check
		fmt.Printf("Generation %d:\n", i)
		PrintGeneration(generation)
		checkFitness, checkTime := TestChromosome(curBestchromo, tasks, nodes, isAddExtra)
		fmt.Printf("current best chromosome, totalTime %d, curfitness %.10f :", checkTime, checkFitness)
		/*
			PrintChromosome(curBestchromo)
			for cnt := 0; cnt < len(curBestchromo); cnt ++ {
				task := tasks[cnt]
				node := nodes[curBestchromo[cnt]]
				predTransTime, predWaitTime, predExecTime, predExtraTime := GetAllPredictTimes(task, node)
				fmt.Printf("trans:%d, wait:%d, exec:%d, extra:%d, size:%f, bw:%f\n", predTransTime, predWaitTime, predExecTime, predExtraTime, task.DataSize, node.Bandwidth)
			}
		*/
		fmt.Printf("last best chromosome, lastBestFitness %.10f :", bestFitness)
		PrintChromosome(bestChromo)
		/*
			for cnt := 0; cnt < len(curBestchromo); cnt ++ {
				task := tasks[cnt]
				node := nodes[bestChromo[cnt]]
				predTransTime, predWaitTime, predExecTime, predExtraTime := GetAllPredictTimes(task, node)
				fmt.Printf("trans:%d, wait:%d, exec:%d, extra:%d, size:%f, bw:%f\n", predTransTime, predWaitTime, predExecTime, predExtraTime, task.DataSize, node.Bandwidth)
			}
		*/
		//end check

		if curBestFitness >= bestFitness {  //
			bestFitness = curBestFitness
			bestChromo = curBestchromo
			if math.Abs(curBestFitness - bestFitness) < FitnessDelta {
				slightDiffcnt ++
				if slightDiffcnt >= ContinuousDiffCnt {
					return  bestChromo
				}
			} else {
				slightDiffcnt = 0
			}
		} else {
			slightDiffcnt = 0
		}

	}
	return bestChromo
}

//
func GaAlgorithmBetter(tasks []*taskmgt.TaskEntity, nodes []*nodemgt.NodeEntity, iteratorNum int, chromosomeNum int, isAddExtra bool) []int {
	fmt.Printf("start ga GaAlgorithmBetter, tasknum:%d, nodenum:%d \n", len(tasks), len(nodes))
	generation := createInitGeneration(len(tasks), len(nodes), chromosomeNum)
	/*
		fmt.Printf("First Generation:\n")
		PrintGeneration(generation)*/
	fitnessList := calFitnessListBetter(generation, tasks, nodes, isAddExtra)
	bestChromo, bestFitness := getMaxChromosome(generation, fitnessList)
	lastBestFitness := bestFitness
	slightDiffcnt := 0
	//fitnesshistory := make([]float64, IteratorNum)

	rand.Seed(time.Now().UnixNano())

	iterCnt := 0
	for ; iterCnt < iteratorNum ; iterCnt++ {

		//select
		//fitnessList = calFitness(generation, tasks, nodes, isAddExtra)
		chosenGeneration := selectChromosome(generation, fitnessList)

		//create new generation
		fitnessList = calFitnessListBetter(chosenGeneration, tasks, nodes, isAddExtra)

		//first, copy part of max N old generation to be new generation
		copyNum := (int)(CpRate*float64(chromosomeNum))
		generation = getMaxNChromosome(chosenGeneration, fitnessList, copyNum)

		//second, create the rest part of new generation by cross and mutation
		childGeneration := crossAndMutation(chosenGeneration, fitnessList, chromosomeNum - copyNum, len(tasks), len(nodes))
		generation = append(generation, childGeneration...)

		fitnessList = calFitnessListBetter(generation, tasks, nodes, isAddExtra)
		curBestchromo, curBestFitness := getMaxChromosome(generation, fitnessList)

		//check
		/*
			fmt.Printf("Generation %d:\n", i)
			PrintGeneration(generation)
			checkFitness, checkTime := TestChromosome(curBestchromo, tasks, nodes, isAddExtra)
			fmt.Printf("current best chromosome, totalTime %d, curfitness %.10f :", checkTime, checkFitness)
			PrintChromosome(curBestchromo)*/
		/*
			for cnt := 0; cnt < len(curBestchromo); cnt ++ {
				task := tasks[cnt]
				node := nodes[curBestchromo[cnt]]
				predTransTime, predWaitTime, predExecTime, predExtraTime := GetAllPredictTimes(task, node)
				fmt.Printf("trans:%d, wait:%d, exec:%d, extra:%d, size:%f, bw:%f\n", predTransTime, predWaitTime, predExecTime, predExtraTime, task.DataSize, node.Bandwidth)
			}
		*/
		/*
			fmt.Printf("last best chromosome, lastBestFitness %.10f :", bestFitness)
			PrintChromosome(bestChromo)*/
		/*
			for cnt := 0; cnt < len(curBestchromo); cnt ++ {
				task := tasks[cnt]
				node := nodes[bestChromo[cnt]]
				predTransTime, predWaitTime, predExecTime, predExtraTime := GetAllPredictTimes(task, node)
				fmt.Printf("trans:%d, wait:%d, exec:%d, extra:%d, size:%f, bw:%f\n", predTransTime, predWaitTime, predExecTime, predExtraTime, task.DataSize, node.Bandwidth)
			}
		*/
		//end check
		lastBestFitness = bestFitness
		if curBestFitness > bestFitness {
			bestFitness = curBestFitness
			bestChromo = curBestchromo
		}

		if bestFitness - lastBestFitness < FitnessDelta {
			slightDiffcnt ++
			if slightDiffcnt >= ContinuousDiffCnt {
				break
			}
		} else {
			slightDiffcnt = 0
		}

		//fitnesshistory[iterCnt] = bestFitness
	}

	/*
		fmt.Printf("trainning history: ")
		for j := 0;j < iterCnt-29;j++ {
			fmt.Printf("%.10f ", fitnesshistory[j])
			if j%10 == 0 {
				fmt.Printf("\n")
			}
		}*/

	//fmt.Printf("Total iteration:%d\n", iterCnt)
	return bestChromo
}

func createInitGeneration(taskNum int, nodeNum int, chromosomeNum int) [][]int{
	generation := make([][]int, chromosomeNum)
	for i := 0; i < chromosomeNum; i++ {
		chromosome := make([]int, taskNum)
		for j := 0; j < taskNum; j++ {
			chromosome[j] = rand.Intn(nodeNum)
		}
		generation[i] = chromosome
	}
	return generation
}

func Fitness(time int64) float64{
	fitness := 1/float64(time)*1e6
	return fitness*fitness
}

func calFitnessList(generation [][]int, tasks []*taskmgt.TaskEntity, nodes []*nodemgt.NodeEntity, isAddExtra bool) []float64 {
	if generation == nil || tasks == nil || nodes == nil {
		return  nil
	}
	fitnessList := make([]float64, len(generation))
	for i := 0; i < len(fitnessList); i++ {
		var totalTime int64 = 0
		chromosome := generation[i]
		for j := 0; j < len(chromosome) ;j++ {
			totalTime += GetTotalTime(tasks[j], nodes[chromosome[j]], isAddExtra)
		}

		fitnessList[i] = Fitness(totalTime)
	}
	return fitnessList
}

//add up extra queueing time
func calFitnessListBetter(generation [][]int, tasks []*taskmgt.TaskEntity, nodes []*nodemgt.NodeEntity, isAddExtra bool) []float64 {
	if generation == nil || tasks == nil || nodes == nil {
		return  nil
	}
	fitnessList := make([]float64, len(generation))
	for i := 0; i < len(fitnessList); i++ {
		var totalTime int64 = 0
		chromosome := generation[i]
		//fmt.Println(chromosome)
		for j := 0; j < len(chromosome) ;j++ {
			totalTime += GetTotalTime(tasks[j], nodes[chromosome[j]], isAddExtra)
		}
		//add up extra queueing time caused by the same node
		if len(chromosome) > 1 {
			nodeCntmap := make(map[int]int)
			for _, v := range chromosome {
				nodeCntmap[v]++
			}
			for key, value := range nodeCntmap {
				//fmt.Printf("key:%d, value:%d\n", key, value)
				node := nodes[key]
				if value > 1 {
					extraWaitTime := int64(value-1) * GetPredictAvgExecTime(node)
					totalTime += extraWaitTime
					//fmt.Printf("extra wait time (%d) by the same node(%s:%d, cnt:%d)\n", extraWaitTime, node.NodeId.IP, node.NodeId.Port, value)
				}
			}
		}

		fitnessList[i] = Fitness(totalTime)
	}
	return fitnessList
}

func calSelectProbability(fitnessList []float64) []float64{
	if fitnessList == nil {
		return  nil
	}
	probabilityList := make([]float64, len(fitnessList))
	var sum float64 = 0
	for i := 0; i < len(fitnessList); i++ {
		sum += fitnessList[i]
	}
	for i :=0; i < len(fitnessList); i++ {
		probabilityList[i] = fitnessList[i]/sum
	}
	return  probabilityList
}

func selectOneIndexByRoulette(probabilityList []float64) int {
	if probabilityList == nil {
		return 0
	}
	probaChosen := rand.Float64()
	var sum float64 = 0
	for i := 0; i < len(probabilityList); i++ {
		sum += probabilityList[i]
		if sum >= probaChosen {
			return  i
		}
	}
	return 0
}

func selectChromosome(lastGeneration [][]int, fitnessList []float64)  [][]int{
	if lastGeneration == nil || fitnessList == nil {
		return  nil
	}
	chromosomeNum := len(lastGeneration)
	probabilityList := calSelectProbability(fitnessList)
	newGeneration := make([][]int, chromosomeNum)
	for i := 0; i < chromosomeNum; i++ {
		index := selectOneIndexByRoulette(probabilityList)
		newGeneration[i] = lastGeneration[index]
	}
	return  newGeneration
}

func getMaxChromosome(generation [][]int, fitnessList []float64 ) (bestChromo []int, bestFitness float64){
	if generation == nil {
		return  nil, 0
	}

	maxIndex := 0
	maxFitness := fitnessList[0]
	for i := 1; i < len(fitnessList); i++ {
		if fitnessList[i] > maxFitness {
			maxIndex = i
			maxFitness = fitnessList[i]
		}
	}
	return generation[maxIndex], fitnessList[maxIndex]
}

type ChromosomeType struct {
	fitness float64
	chromosome []int
}

func getMaxNChromosome(generation [][]int, fitnessList []float64, num int) [][]int {
	chromoNum := len(generation)
	chromosomeGrp := make([]ChromosomeType, chromoNum)
	for i := 0; i < chromoNum; i++ {
		chromosomeGrp[i] = ChromosomeType{
			fitness:    fitnessList[i],
			chromosome: generation[i],
		}
	}

	//decrease
	sort.Slice(chromosomeGrp, func(i, j int) bool {
		return chromosomeGrp[i].fitness > chromosomeGrp[j].fitness
	})

	maxGeneration := make([][]int, num)
	for i := 0; i < num; i++ {
		maxGeneration[i] = chromosomeGrp[i].chromosome
	}
	return  maxGeneration
}

func crossAndMutation(generation [][]int, fitnessList []float64, crossNum int, taskNum int, NodeNum int) [][]int{
	if generation == nil || fitnessList == nil {
		return  nil
	}
	probabilityList := calSelectProbability(fitnessList)

	childGeneration := make([][]int, crossNum)
	for i := 0; i < crossNum; i++ {
		chromosome := make([]int, 0)
		//cross
		fatherIndex := selectOneIndexByRoulette(probabilityList)
		motherIndex := selectOneIndexByRoulette(probabilityList)
		headLen := rand.Intn(taskNum)
		chromosome = append(chromosome, generation[fatherIndex][0:headLen]...)
		chromosome = append(chromosome, generation[motherIndex][headLen:]...)

		//mutation
		mutationCnt := taskNum/3 + 1
		for j := 0; j < mutationCnt; j++ {
			mutationIndex := rand.Intn(taskNum)
			mutaionNodeVal := rand.Intn(NodeNum)
			chromosome[mutationIndex] = mutaionNodeVal
		}

		childGeneration[i] = chromosome
	}
	return  childGeneration
}

func TestChromosome(chromosome []int, tasks []*taskmgt.TaskEntity, nodes []*nodemgt.NodeEntity, isAddExtra bool)  (fitness float64, totalTime int64) {
	if chromosome == nil || tasks == nil || nodes == nil {
		return 0, 0
	}

	for i := 0; i < len(chromosome); i++ {
		totalTime += GetTotalTime(tasks[i], nodes[chromosome[i]], isAddExtra)
	}
	fitness = Fitness(totalTime)

	return
}

func PrintChromosome(chromo []int) {
	//fmt.Println(chromo)
	fmt.Printf("[")
	for i := 0; i < len(chromo); i++ {
		fmt.Printf("%d,",chromo[i])
	}
	fmt.Printf("]\n")
}

func PrintGeneration(generation [][]int) {
	for i := 0; i < len(generation); i++ {
		PrintChromosome(generation[i])
	}
}