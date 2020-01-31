package predictor

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"taskAssignmentForEdge/common"
	"taskAssignmentForEdge/master/histogram"
	"taskAssignmentForEdge/taskmgt"
)

//Feature constants
const (
	FeatureType_Job_Name = iota
	FeatureType_logic_Job_Name
	FeatureType_User_Name
	FeatureType_All
)

const (
	FeatureVal_All = iota
)

const (
	Max_Predict_Unit_Point_Num = 40
	Max_Recent_Eelement_Num = 20
	Max_Recent_Eelement_Num_For_Wait = 10
)

type FeaturePair struct {
	FeatureType int
	FeatureVal interface{}
}

//For runtime
type PredictUnit struct {
	Feature FeaturePair
	histog *histogram.Histogram

	//estimators
	recentQueue []float64
	avg float64               //average of all averages
	median float64			  //median of recent elements
	ewma float64              //Exponential Weighted Moving Average
	recentAvg float64		  //average of recent elements

	//mean absolute error, need to normalized( mae/(max-min)) when predicting
 	maeofavg float64
	maeofmedian float64
	maeofewma float64
	maeofrecentAvg float64

	maxv float64
	minv float64

	totalRecordNum int

	rwlock sync.RWMutex
}

func NewPredictUnit(fpair FeaturePair) *PredictUnit {
	preUnit := new(PredictUnit)
	preUnit.Feature = fpair
	preUnit.histog = histogram.NewHistogram(Max_Predict_Unit_Point_Num)
	//preUnit.recentQueue = queue.New()
	preUnit.recentQueue = []float64{}
	return preUnit
}

func (preUnit *PredictUnit) Update(val int64) {

	preUnit.rwlock.Lock()
	defer preUnit.rwlock.Unlock()

	valf := float64(val)

	//update histogram
	preUnit.histog.Update(valf, 1)

	//update max and min
	if preUnit.totalRecordNum == 0 {
		preUnit.maxv = valf
		preUnit.minv = valf
	} else {
		if preUnit.maxv < valf {
			preUnit.maxv = valf
		}
		if preUnit.minv > valf {
			preUnit.minv = valf
		}
	}

	//update MAE of the last prediction of esimator pair
	avg := float64(preUnit.totalRecordNum)/float64(preUnit.totalRecordNum+1)*preUnit.avg + valf/float64(preUnit.totalRecordNum+1)
	if preUnit.totalRecordNum > 0 {
		//avg
		preUnit.maeofavg = preUnit.maeofavg/avg*preUnit.avg + math.Abs(valf-preUnit.avg)/avg
		//median
		preUnit.maeofmedian = preUnit.maeofmedian/avg*preUnit.avg + math.Abs(valf-preUnit.median)/avg
		//rolling
		preUnit.maeofewma = preUnit.maeofewma/avg*preUnit.avg + math.Abs(valf-preUnit.ewma)/avg
		//recentavg
		preUnit.maeofrecentAvg = preUnit.maeofrecentAvg/avg*preUnit.avg + math.Abs(valf-preUnit.recentAvg)/avg
	}

	//update four esimator pair
	//avg
	preUnit.avg = avg

	//maintain Recent_Val_Max_NUM elements
	preUnit.recentQueue = append(preUnit.recentQueue, valf)
	if len(preUnit.recentQueue) > Max_Recent_Eelement_Num {
		preUnit.recentQueue = preUnit.recentQueue[1:]
	}

	//rolling
	const alpha float64 = 0.5
	if preUnit.totalRecordNum == 0 {
		preUnit.ewma = valf
	} else {
		preUnit.ewma = alpha*preUnit.ewma + (1-alpha)*valf
	}

	//recentavg
	var sum float64 = 0.0
	for _, v := range preUnit.recentQueue {
		sum += float64(v)
	}
	preUnit.recentAvg = sum / float64(len(preUnit.recentQueue))

	//recent median
	cpyRecentQ := make([]float64, len(preUnit.recentQueue))
	copy(cpyRecentQ, preUnit.recentQueue)
	sort.Float64s(cpyRecentQ)
	preUnit.median = float64(cpyRecentQ[len(cpyRecentQ)/2])

	//update the total number of records at last
	preUnit.totalRecordNum++
}

func (preUnit *PredictUnit) ShowHistogram( ) {
	fmt.Println(preUnit.histog)
}

type EstimatorResult struct {
	estimate float64
	nmae float64
}

type PredictResult struct {
	Feature FeaturePair
	EtmRes []EstimatorResult
}

func (preUnit *PredictUnit) Predict() PredictResult {
	preUnit.rwlock.RLock()
	defer preUnit.rwlock.RUnlock()

	diff := 0.0
	if common.IsValueEqual(preUnit.minv, preUnit.maxv) {
		diff = 0.0000001
	} else {
		diff = preUnit.maxv - preUnit.minv
	}

	result := []EstimatorResult{{preUnit.avg, preUnit.maeofavg/diff},{preUnit.median, preUnit.maeofmedian/diff},
		{preUnit.ewma, preUnit.maeofewma/diff},{preUnit.recentAvg, preUnit.maeofrecentAvg/diff}}

	return PredictResult{preUnit.Feature, result}
}

//used for default feature
func (preUnit *PredictUnit) SinglePredict() (estimate float64) {
	preUnit.rwlock.RLock()
	defer preUnit.rwlock.RUnlock()

	diff := 0.0
	if common.IsValueEqual(preUnit.minv, preUnit.maxv) {
		diff = 0.0000001
	} else {
		diff = preUnit.maxv - preUnit.minv
	}

	result := []EstimatorResult{{preUnit.avg, preUnit.maeofavg/diff},{preUnit.median, preUnit.maeofmedian/diff},
		{preUnit.ewma, preUnit.maeofewma/diff},{preUnit.recentAvg, preUnit.maeofrecentAvg/diff}}

	//find the minimum
	minnmae := result[0].nmae
	estimate = result[0].estimate
	for i:=1; i < len(result); i++ {
		if result[i].nmae < minnmae {
			minnmae = result[i].nmae
			estimate = result[i].estimate
		}
	}

	return estimate
}

func CreateFeaturePairs(task *taskmgt.TaskEntity) []FeaturePair{
	if task == nil {
		return nil
	}
	fpair := []FeaturePair{{FeatureType_Job_Name, task.TaskName},
		{FeatureType_logic_Job_Name, task.LogicName},
		{ FeatureType_User_Name, task.Username}}
	return fpair
}

//For waitint time, without histogram
/*
type PredictUnitSimple struct {
	//estimators
	recentQueue []float64
	avg float64               //average of all averages
	median float64			  //median of recent elements
	ewma float64              //Exponential Weighted Moving Average
	recentAvg float64		  //average of recent elements

	maeofavg float64
	maeofmedian float64
	maeofewma float64
	maeofrecentAvg float64

	totalRecordNum int
}*/

/*
func NewPredictUnitSimple() *PredictUnitSimple {
	preUnit := new(PredictUnitSimple)
	preUnit.recentQueue = []float64{}
	return preUnit
}

func (preUnit *PredictUnitSimple) Update(val int64) {

	valf := float64(val)

	//update MAE of the last prediction of esimator pair
	avg := float64(preUnit.totalRecordNum)/float64(preUnit.totalRecordNum+1)*preUnit.avg + valf/float64(preUnit.totalRecordNum+1)
	if preUnit.totalRecordNum > 0 {
		//avg
		preUnit.maeofavg = preUnit.maeofavg/avg*preUnit.avg + math.Abs(valf-preUnit.avg)/avg
		//median
		preUnit.maeofmedian = preUnit.maeofmedian/avg*preUnit.avg + math.Abs(valf-preUnit.median)/avg
		//rolling
		preUnit.maeofewma = preUnit.maeofewma/avg*preUnit.avg + math.Abs(valf-preUnit.ewma)/avg
		//recentavg
		preUnit.maeofrecentAvg = preUnit.maeofrecentAvg/avg*preUnit.avg + math.Abs(valf-preUnit.recentAvg)/avg
	}

	//update four esimator pair
	//avg
	preUnit.avg = avg

	//maintain Recent_Val_Max_NUM elements
	preUnit.recentQueue = append(preUnit.recentQueue, valf)
	if len(preUnit.recentQueue) > Max_Recent_Eelement_Num_For_Wait {
		preUnit.recentQueue = preUnit.recentQueue[1:]
	}

	//rolling
	const alpha float64 = 0.5
	if preUnit.totalRecordNum == 0 {
		preUnit.ewma = valf
	} else {
		preUnit.ewma = alpha*preUnit.ewma + (1-alpha)*valf
	}

	//recentavg
	var sum float64 = 0.0
	for _, v := range preUnit.recentQueue {
		sum += float64(v)
	}
	preUnit.recentAvg = sum / float64(len(preUnit.recentQueue))

	//recent median
	cpyRecentQ := make([]float64, len(preUnit.recentQueue))
	copy(cpyRecentQ, preUnit.recentQueue)
	sort.Float64s(cpyRecentQ)
	preUnit.median = float64(cpyRecentQ[len(cpyRecentQ)/2])

	//update the total number of records at last
	preUnit.totalRecordNum++
}

func (preUnit *PredictUnitSimple) Predict()  float64 {
	result := []EstimatorResult{{preUnit.avg, preUnit.maeofavg},{preUnit.median, preUnit.maeofmedian},
		{preUnit.ewma, preUnit.maeofewma},{preUnit.recentAvg, preUnit.maeofrecentAvg}}
}*/