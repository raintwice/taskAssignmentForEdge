package connect

import (
	"encoding/csv"
	"google.golang.org/grpc"
	"net"
	"os"
	"strconv"
	"sync"
	"log"
	"taskAssignmentForEdge/common"
	pb "taskAssignmentForEdge/proto"
	"golang.org/x/net/context"
	"taskAssignmentForEdge/taskmgt"
)

//client as recevier
func (clt *Client) StartRecvResultServer(wg *sync.WaitGroup) {
	lis, err := net.Listen("tcp", ":"+ strconv.Itoa(common.ClientPort))
	if err != nil {
		log.Fatal("Failed to start client server for receving results: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMaster2ClientConnServer(s, clt)
	s.Serve(lis)

	wg.Done()
}

func (clt *Client) ReturnSubmittedTasks(ctx context.Context, in *pb.TaskSubmitResReq) (*pb.TaskSubmitResResp, error) {
	//update attribute
	if  clt.EvaluationStatus == EvaluationStatus_Pretrain {
		for i := 0; i < len(in.InfoGp); i++ {
			log.Printf("Get result of task(Id:%d) from master", in.InfoGp[i].TaskId)
			taskmgt.TranslateProtoTaskToRecord(in.InfoGp[i],clt.PreTrainTaskInfoGrp[clt.recvPretrainTaskCnt])
			clt.recvPretrainTaskCnt++
		}

		if clt.recvPretrainTaskCnt >= clt.PretrainNum {
			clt.EvaluationStatus = EvaluationStatus_Evaluation
			//save the file
			newfile, err := os.Create("./PretrainTaskResult.csv")
			if err != nil {
				log.Printf("Error: cannot create PretrainTaskResult.csv")
			} else {
				defer newfile.Close()
				writer := csv.NewWriter(newfile)
				if werr := writer.WriteAll(clt.PreTrainTaskInfoGrp); werr != nil {
					log.Printf("Error: cannot write to  PretrainTaskResult.csv")
				}
			}
		}
	} else if clt.EvaluationStatus == EvaluationStatus_Evaluation {
		for i := 0; i < len(in.InfoGp);i++ {
			log.Printf("Get result of task(Id:%d) from master", in.InfoGp[i].TaskId)
			taskmgt.TranslateProtoTaskToRecord(in.InfoGp[i], clt.EvalTaskInfoGrp[clt.recvEvalTaskCnt])
			clt.recvEvalTaskCnt++
		}
		if clt.recvEvalTaskCnt >=  clt.EvalSamplesNum {
			clt.EvaluationStatus = EvaluationStatus_Finish
			//save the file
			newfile, err := os.Create("./EvalTaskResult.csv")
			if err != nil {
				log.Printf("Error: cannot create EvalTaskResult.csv")
			} else {
				defer newfile.Close()
				writer := csv.NewWriter(newfile)
				if werr := writer.WriteAll(clt.EvalTaskInfoGrp); werr != nil {
					log.Printf("Error: cannot write to EvalTaskResult.csv")
				}
			}
		}
	}

	return &pb.TaskSubmitResResp{Reply:true}, nil
}