package connect

import (
	"errors"
	"io"
	"log"
	"os"
	"sync"
	"taskAssignmentForEdge/common"
	"taskAssignmentForEdge/master/nodemgt"
	"taskAssignmentForEdge/taskmgt"
	"time"
	pb "taskAssignmentForEdge/proto"
	"golang.org/x/net/context"
)

//分配

func (ms *Master) AssignOneTask(task *taskmgt.TaskEntity, node *nodemgt.NodeEntity) (* pb.SendStatus, error) {
	if task == nil || node == nil {
		return nil, errors.New("The task to assign is empty!")
	}

	var (
		writing   = true
		buf       []byte
		n         int
		file      *os.File
		status    *pb.SendStatus
		err       error
		chunkSize  = 1<<20
		isSendInfo = false
	)

	c := pb.NewMaster2NodeConnClient(node.Conn)
	file, err = os.Open(task.TaskName)
	if err != nil {
		log.Printf("Cannot open task file %s ", task.TaskName)
		return nil, err
	}
	defer file.Close()

	stream, err := c.AssignTask(context.Background())
	if err != nil {
		log.Printf("Failed to create AssignTask stream for task file %s", task.TaskName)
		return nil, err
	}
	defer stream.CloseSend()

	//var stats Stats
	//stats.StartedAt = time.Now()

	buf = make([]byte, chunkSize)
	for writing {
		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}
			log.Printf("Failed to copy from task file %s to buffer", task.TaskName)
			return nil, err
		}

		if isSendInfo {
			err = stream.Send(&pb.TaskChunk{
				Content: buf[:n],
			})
		} else {
			err = stream.Send(&pb.TaskChunk{
				Info:&pb.TaskInfo{
					TaskName: task.TaskName,
					TaskId: task.TaskId,
				},
				Content: buf[:n],
			})
			isSendInfo = true
		}

		if err != nil {
			log.Printf("Failed to send chunk via stream for task file %s", task.TaskName)
			return nil, err
		}
	}

	//stats.FinishedAt = time.Now()

	status, err = stream.CloseAndRecv()
	if err != nil {
		log.Printf("Failed to receive response from node(IP:%s)", node.IpAddr)
		return nil, err
	}

	if status.Code != pb.SendStatusCode_Ok {
		log.Printf("Failed to assign task file %s", task.TaskName)
		return status, errors.New(
			"Failed to assign task file")
	}

	//将发送成功的任务 移到 已发送任务队列
	node.TqPrepare.Rwlock.RUnlock()
	node.TqPrepare.RemoveTask(task.TaskId)
	node.TqPrepare.Rwlock.RLock()
	node.TqAssign.AddTask(task)

	return status, nil
}

func (ms *Master) AssignTaskForNode(node *nodemgt.NodeEntity) {
	if node == nil || node.TqPrepare.TaskNum == 0 {
		return
	}
	node.TqPrepare.Rwlock.RLock()
	for e := node.TqPrepare.TaskList.Front(); e != nil; e = e.Next() {
		task := e.Value.(* taskmgt.TaskEntity)
		ms.AssignOneTask(task, node)
		//TBD, 分配失败先不管了
	}
	node.TqPrepare.Rwlock.RUnlock()
}

func (ms *Master) StartDispatcher(wg *sync.WaitGroup) {
	for range time.Tick(time.Millisecond*common.AssgnTimeout) {
		ms.Nq.Rwlock.RLock()
		for e:= ms.Nq.NodeList.Front(); e != nil; e = e.Next() {
			node := e.Value.(*nodemgt.NodeEntity)
			go ms.AssignTaskForNode(node)
		}
		ms.Nq.Rwlock.RUnlock()
	}

	if wg != nil {
		wg.Done()
	}
}