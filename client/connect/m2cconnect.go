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
		clt.recvPreLock.Lock()
		defer clt.recvPreLock.Unlock()
		for i := 0; i < len(in.InfoGp); i++ {
			log.Printf("Get result of pretrain task(Id:%d, cnt:%d) from master", in.InfoGp[i].TaskId, clt.recvPretrainTaskCnt)
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
		clt.recvEvalLock.Lock()
		defer clt.recvEvalLock.Unlock()
		for i := 0; i < len(in.InfoGp);i++ {
			log.Printf("Get result of evaluation task(Id:%d,cnt:%d) from master", in.InfoGp[i].TaskId, clt.recvEvalTaskCnt)
			taskmgt.TranslateProtoTaskToRecord(in.InfoGp[i], clt.EvalTaskInfoGrp[clt.recvEvalTaskCnt])
			clt.recvEvalTaskCnt++
		}
		if clt.recvEvalTaskCnt >=  clt.EvalSamplesNum {
			//save the file
			newfile, err := os.Create("./EvalTaskResult.csv")
			if err != nil {
				log.Printf("Error: cannot create EvalTaskResult.csv")
			} else {
				defer newfile.Close()
				writer := csv.NewWriter(newfile)
				if werr := writer.WriteAll(clt.EvalTaskInfoGrp); werr != nil {
					log.Printf("Error: cannot write to EvalTaskResult.csv")
				} else {
					log.Printf("Success!Evaluation is finished")
					clt.EvaluationStatus = EvaluationStatus_Finish
					os.Exit(0)
				}
			}
		}
	}

	return &pb.TaskSubmitResResp{Reply:true}, nil
}