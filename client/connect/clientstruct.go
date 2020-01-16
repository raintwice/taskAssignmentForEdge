package connect

import (
	"google.golang.org/grpc"
	"taskAssignmentForEdge/taskmgt"
)
/*
type JobSubmitConfig struct {
	EvaluationDir    string
	jobArrivalRate   float64
	EvaluationStatus int
	PretrainNum 	 int
	EvalSamplesNum    int
}*/

const (
	EvaluationStatus_Pretrain = iota
	EvaluationStatus_Evaluation
	EvaluationStatus_Finish
)

type Client struct {
//attribute for sending
	EvaluationDir    string
	jobArrivalRate   float64  //x samples per min
	EvaluationStatus int
	PretrainNum 	 int
	EvalSamplesNum    int

	SentTaskCnt    int

//attribute for receiving
	recvPretrainTaskCnt 	int
	PreTrainTaskInfoGrp   [][]string
	recvEvalTaskCnt     int
	EvalTaskInfoGrp       [][]string

	conn *grpc.ClientConn  //call master
}

func NewClient(dir string, rate float64, pretainNum int, evalNum int) (* Client) {
	clt := &Client{}
	clt.EvaluationDir = dir
	clt.PretrainNum = pretainNum
	clt.jobArrivalRate = rate
	clt.EvalSamplesNum = evalNum

	clt.EvaluationStatus = EvaluationStatus_Pretrain
	clt.recvPretrainTaskCnt = 0
	clt.PreTrainTaskInfoGrp = make([][]string, clt.PretrainNum)
	for i := 0; i < clt.PretrainNum; i++ {
		clt.PreTrainTaskInfoGrp[i] = make([]string, taskmgt.TaskAttributeNum)
	}
	clt.recvEvalTaskCnt = 0
	clt.EvalTaskInfoGrp = make([][]string, clt.EvalSamplesNum)
	for i := 0; i < clt.EvalSamplesNum ; i++ {
		clt.EvalTaskInfoGrp[i] = make([]string, taskmgt.TaskAttributeNum)
	}
	return clt
}
