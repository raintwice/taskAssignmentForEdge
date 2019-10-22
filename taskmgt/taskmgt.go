package taskmgt

import (
	"container/list"
)

/*
* Task
*/

type TaskEntity struct {
	TaskId int32 `json:"TaskId"`
	TaskName string `json:"TaskName"`
	TaskLocation string `json:"TaskLocation"`
	//其他属性
}

//创建任务
func CreateTask(TaskId int32) *TaskEntity {
	task := new(TaskEntity)
	task.TaskId = TaskId
	return task
}


/*
* TaskQueue
*/

type TaskQueue struct {
	TaskList list.List
	TaskTable map[int32](*list.Element) 
	TaskNum int 
}

func NewTaskQueue()  *TaskQueue {
	return &TaskQueue{
		TaskList:  list.List{},
		TaskTable: make(map[int32](*list.Element)),
		TaskNum:   0,
	}
}

func (tq *TaskQueue) AddTask(task *TaskEntity) {
	e := tq.TaskList.PushBack(task)
	tq.TaskTable[task.TaskId] = e
	tq.TaskNum++
}

func (tq *TaskQueue) RemoveTask(TaskId int32) {
	if e, ok := tq.TaskTable[TaskId]; ok {
		delete(tq.TaskTable, TaskId)
		tq.TaskList.Remove(e)
		tq.TaskNum--
	}
}

func (tq *TaskQueue) FindTask(TaskId int32) *TaskEntity {
	if e, ok := tq.TaskTable[TaskId]; ok {
		return e.Value.(*TaskEntity)
	}
	return nil
}

func (tq *TaskQueue) GettaskNum() int {
	return tq.TaskNum
}


//从json文件读入任务列表
/*func ReadTaskList(listpath string) {
	filePtr, err := os.Open(listpath)
  if err != nil {
    fmt.Println("Open file failed [Err:%s]", err.Error())
    return
  }
  defer filePtr.Close()
  log.Printf("opened %s", filePtr)

  var task []TaskEntity
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
}*/