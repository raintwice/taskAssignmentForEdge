package predictor

import (
    "taskAssignmentForEdge/common"
)

type RunTimePredictManager struct {
    PredictorList *common.ListWithHash
}

func NewRunTimePredictManager() *RunTimePredictManager{
    return &RunTimePredictManager{common.NewListWithHash()}
}

func (predictMng *RunTimePredictManager) FindRunTimePredictor(machineType int) *RunTimePredictor {
   ev := predictMng.PredictorList.Find(machineType)
   if ev != nil {
       return ev.(*RunTimePredictor)
   } else {
       return nil
   }
}

func (predictMng *RunTimePredictManager) GetRunTimePredictor(machineType int)  *RunTimePredictor {
    predictor := predictMng.FindRunTimePredictor(machineType)
    if predictor == nil {
        predictor = NewRunTimePredictor(machineType)
        predictMng.PredictorList.AddUnique(machineType, predictor)
    }
    //predictor.Cnt++
    return predictor
}

type ConnPredictManager struct {
    PredictorList *common.ListWithHash
}

func NewConnPredictManager() *ConnPredictManager{
    return &ConnPredictManager{common.NewListWithHash()}
}

func (predictMng *ConnPredictManager) FindConnPredictManager(groupIndex int) *ConnPredictor {
    ev := predictMng.PredictorList.Find(groupIndex)
    if ev != nil {
        return ev.(*ConnPredictor)
    } else {
        return nil
    }
}

func (predictMng *ConnPredictManager) GetConnPredictor(groupIndex int)  *ConnPredictor {
    predictor := predictMng.FindConnPredictManager(groupIndex)
    if predictor == nil {
        predictor = NewConnPredictor(groupIndex)
        predictMng.PredictorList.AddUnique(groupIndex, predictor)
    }
    return predictor
}