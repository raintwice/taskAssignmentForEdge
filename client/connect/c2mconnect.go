package connect

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"taskAssignmentForEdge/common"
	pb "taskAssignmentForEdge/proto"
	"taskAssignmentForEdge/taskmgt"
	"time"
)

//client as recevier
/*
func (clt *Client) SubmitTask(taskgp []*taskmgt.TaskEntity) {

	conn, err := grpc.Dial("127.0.0.1"+":"+strconv.Itoa(common.MasterPortForClient), grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cannot not connect with master(%s:%d): %v", "", "127.0.0.1", common.MasterPortForClient, err)
	}
	defer conn.Close()
	c := pb.NewClient2MasterConnClient(conn)

	infogp := make([]*pb.TaskInfo, 0)
	for _, task := range taskgp {
		taskinfo := pb.TaskInfo{
			TaskName: task.TaskName,
		}
		infogp = append(infogp, &taskinfo)
	}

	taskReq := pb.TaskSubmitReq{TaskGp: infogp}

	r, err := c.SubmitTask(context.Background(), &taskReq)
	if err != nil {
		log.Fatalf("Could not submit task: %v", err)
	}
	if (r.Reply) {
		log.Printf("Successed to join")
	} else {
		log.Printf("Failed to join")
	}
}*/

func (clt *Client) SubmitOneTask(taskReq *pb.TaskSubmitReq ) {

	conn, err := grpc.Dial("127.0.0.1"+":"+strconv.Itoa(common.MasterPortForClient), grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cannot not connect with master(%s:%d): %v", "", "127.0.0.1", common.MasterPortForClient, err)
	}
	defer conn.Close()
	c := pb.NewClient2MasterConnClient(conn)

	r, err := c.SubmitTasks(context.Background(), taskReq)
	if err != nil {
		log.Fatalf("Could not submit task: %v", err)
	}
	if (r.Reply) {
		log.Printf("Successed to submit task %d", clt.SentTaskCnt)
		clt.SentTaskCnt++
	} else {
		log.Printf("Failed to submit task ")
	}
}

//possion process, get the interval, per minutes
func NextTime(ratePara float64) float64{
	return -math.Log(1.0 - rand.Float64())/ratePara
}

//filter_name = ['user name','cpu','ram', 'disk','runtime','job_name','ljobname','datasize','deadlineslack']
//build tasks from csv, and send them

func (clt *Client)buildSendOneTask(record []string) {
	task := taskmgt.TranslateRecordToProtoTask(record)
	//For simulation, reduce to 1/5
	task.RuntimePreSet = task.RuntimePreSet/5
	infogp := []*pb.TaskInfo{task}
	taskReq := pb.TaskSubmitReq{TaskGp: infogp}
	clt.SubmitOneTask(&taskReq)
}

func (clt *Client) ProduceTasks() {
	//log.Printf("Begin to produce tasks")
	files, _ := ioutil.ReadDir(clt.EvaluationDir)
	if len(files) == 0 {

	}
	//fmt.Println(files)
	for _, file := range files{
		//fmt.Println(file)
		fileOpen, err := os.Open(clt.EvaluationDir + "/" + file.Name())
		if err != nil {
			fmt.Printf("Cannot open file[%s], err: %v\n", file.Name(), err)
			continue
		}
		defer fileOpen.Close()
		reader := csv.NewReader(fileOpen)
		records, rerr := reader.ReadAll()
		if rerr != nil {
			fmt.Printf("Cannot read csv file[%s]\n", file.Name())
		}

		log.Printf("Begin to produce pretain tasks")
		//clt.EvaluationStatus = EvaluationStatus_Pretrain
		recordCnt := 0
		for ; recordCnt < clt.PretrainNum; recordCnt++ {
			interval := NextTime(clt.jobArrivalRate) //in min
			duration := time.Duration(interval*float64(time.Minute))
			//log.Printf("next time: %v", duration)
			time.Sleep(duration)
			clt.buildSendOneTask(records[recordCnt])
		}

		//wait until all the pretrain tasks are received
		for ;clt.EvaluationStatus == EvaluationStatus_Pretrain; {
			time.Sleep(5*time.Second)
		}

		log.Printf("Begin to produce evaluation tasks")
		for ;recordCnt < clt.PretrainNum + clt.EvalSamplesNum; recordCnt++ {
			interval := NextTime(clt.jobArrivalRate)
			duration := time.Duration(interval*float64(time.Second)*60)
			time.Sleep(duration)
			clt.buildSendOneTask(records[recordCnt])
		}
	}

	/*
	for ;clt.EvaluationStatus < EvaluationStatus_Finish ; {
		time.Sleep(3*time.Second)
	}*/
}

func (clt *Client) StartTaskProducer(wg *sync.WaitGroup) {
	time.Sleep(3*time.Second)
	clt.ProduceTasks()
	wg.Done()
}