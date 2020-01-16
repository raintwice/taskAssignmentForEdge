package predictor

import (
	"taskAssignmentForEdge/common"
	"taskAssignmentForEdge/master/histogram"
	"fmt"
	"log"
)

//one default predictUnit contains all values, several predictUnit per one feature-val
type RunTimePredictor struct {
	UnitList *common.ListWithHash
	MachineType int      //indicate the node which owns this runtimePredictor
	Cnt int              //
}

func NewRunTimePredictor(machineType int) *RunTimePredictor{
	predictor := &RunTimePredictor{common.NewListWithHash(), machineType, 0}

	//default predictUnit for every feature
	fpair := FeaturePair{FeatureType_All, FeatureVal_All}
	newUnit := NewPredictUnit(fpair)
	predictor.UnitList.AddUnique(fpair, newUnit)

	return predictor
}

func (predictor *RunTimePredictor) GetDefaultUnit() (defaultUnit *PredictUnit) {
	fpair := FeaturePair{FeatureType_All, FeatureVal_All}
	ev := predictor.UnitList.Find(fpair)
	if ev == nil {

		fmt.Printf("Error, Cannot find default PredictUnit")
		defaultUnit = NewPredictUnit(fpair)
		predictor.UnitList.AddUnique(fpair, defaultUnit)
	} else{
		defaultUnit = ev.(*PredictUnit)
	}

	return defaultUnit
}

func (predictor *RunTimePredictor) Update(FeaturePairList []FeaturePair, val int64) {
	predictor.GetDefaultUnit().Update(val)
	for _, fpair := range FeaturePairList {
		if ev := predictor.UnitList.Find(fpair); ev != nil {
			unit := ev.(*PredictUnit)
			unit.Update(val)
		} else {
			newUnit := NewPredictUnit(fpair)
			predictor.UnitList.AddUnique(fpair, newUnit)
			newUnit.Update(val)
		}
	}
}

func (predictor *RunTimePredictor) Predict(FeaturePairList []FeaturePair)  (bestEstimate float64, bestFec FeaturePair){
	results := []PredictResult{}
	for _, fpair := range FeaturePairList {
		if ev := predictor.UnitList.Find(fpair); ev != nil {
			unit := ev.(*PredictUnit)
			results = append(results, unit.Predict())
		}
	}
	if len(results) == 0 {
		//return 0.0, FeaturePair{}
		results = append(results, predictor.GetDefaultUnit().Predict())
	}

	//fmt.Println(results)
	minnmae := results[0].EtmRes[0].nmae
	bestEstimate = results[0].EtmRes[0].estimate
	bestFec = results[0].Feature
	for i := 0;i < len(results);i++ {
		for j := 0; j < len(results[i].EtmRes); j++ {
			if results[i].EtmRes[i].nmae < minnmae {
				minnmae = results[i].EtmRes[i].nmae
				bestEstimate = results[i].EtmRes[i].estimate
				bestFec = results[i].Feature
			}
		}
	}
	return bestEstimate, bestFec
}

func (predictor *RunTimePredictor) ShowPredictorUnit() {
	lh := predictor.UnitList.GetList()
	for e := lh.Front(); e != nil ; e = e.Next() {
		fpair := e.Value.(*PredictUnit).Feature
		log.Printf("Feature:%d, value:%v\n", fpair.FeatureType, fpair.FeatureVal)
	}
}

func (predictor *RunTimePredictor) FindHistogByFec(fec FeaturePair) *histogram.Histogram{
	if ev := predictor.UnitList.Find(fec); ev != nil {
		unit := ev.(*PredictUnit)
		return unit.histog
	} else {
		return nil
	}
}

type ConnPredictor struct {
	histog *histogram.Histogram
	GroupIndex int
}

func NewConnPredictor(groupIndex int) *ConnPredictor {
	predictor := new(ConnPredictor)
	predictor.histog = histogram.NewHistogram(Max_Predict_Unit_Point_Num)
	predictor.GroupIndex = groupIndex
	return predictor
}

func (predictor *ConnPredictor) GetHistogram()  *histogram.Histogram{
	return predictor.histog
}

func (predictor *ConnPredictor) Update(val int64) {
	predictor.histog.Update(float64(val), 1)
}
