package taskmgt


import (
	"container/list"
	"fmt"
	"os"
	"log"
	"encoding/json"
)

type taskEntity struct {
	TaskId int32 `json:"TaskId"`
	TaskName string `json:"TaskName"`
	TaskLocation string `json:"TaskLocation"`
	//其他属性
}

var taskList list.List
var taskTable = make(map[int32] (*list.Element) )

//创建任务并加入队列
func CreateTask(TaskId int32) *taskEntity {
	task := new(taskEntity)
	task.TaskId = TaskId
	e := taskList.PushBack(task)
	taskTable[TaskId] = e
	log.Printf("entity of task(Id:%d) is created", TaskId);
	return e.Value.(*taskEntity)
}

//从队列中删除任务
func DeleteTask(TaskId int32) {
	if e, ok := taskTable[TaskId]; ok {
		delete(taskTable, TaskId)
		//n := e.Value.(*taskEntity)
		//n = nil
		taskList.Remove(e)
		log.Printf("entity of task(IP:%d) is deleted", TaskId);
	}
}

//根据TaskId 查找任务
func FindTask(TaskId int32) *taskEntity {
	if e, ok := taskTable[TaskId]; ok {
		return e.Value.(*taskEntity)
	}
	return nil
}

//查看当前队列有多少任务
func GettaskNum() int {
	return taskList.Len()
}

//从json文件读入任务列表
func ReadTaskList(listpath string) {
	filePtr, err := os.Open(listpath)
  if err != nil {
    fmt.Println("Open file failed [Err:%s]", err.Error())
    return
  }
  defer filePtr.Close()
  log.Printf("opened %s", filePtr)

  var task []taskEntity
  // 创建json解码器
  decoder := json.NewDecoder(filePtr)
  err = decoder.Decode(&task)
  if err != nil {
    fmt.Println("Decoder failed", err.Error())
  } else {
    fmt.Println("Decoder success")
    fmt.Println(task)
  }
  fmt.Println("---taskloading end---")
}