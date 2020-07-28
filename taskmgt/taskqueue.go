package taskmgt

import (
	"container/list"
	"log"
	"sort"
	"sync"
)

/*
* TaskQueue
*/

type TaskQueue struct {
	Name string
	TaskList list.List
	TaskTable map[int32](*list.Element)

	Rwlock sync.RWMutex
}

func NewTaskQueue(name string)  *TaskQueue {
	return &TaskQueue{
		Name : name,
		TaskList:  list.List{},
		TaskTable: make(map[int32](*list.Element)),
	}
}

func (tq *TaskQueue) EnqueueTask(task *TaskEntity) {
	if task == nil {
		log.Printf("Error: empty task when enqueueing in queue(%s)", tq.Name)
	}
	tq.Rwlock.Lock()
	e := tq.TaskList.PushBack(task)
	tq.TaskTable[task.TaskId] = e
	tq.Rwlock.Unlock()
}

func (tq *TaskQueue) DequeueTask(TaskId int32) (task *TaskEntity) {
	tq.Rwlock.Lock()
	if e, ok := tq.TaskTable[TaskId]; ok {
		task = e.Value.(*TaskEntity)
		delete(tq.TaskTable, TaskId)
		tq.TaskList.Remove(e)
	} else {
		task = nil
	}
	tq.Rwlock.Unlock()
	return
}

func (tq *TaskQueue) DequeueFirstTask( ) (task *TaskEntity) {
	if tq.TaskList.Len() == 0 {
		return nil
	}
	tq.Rwlock.Lock()
	e := tq.TaskList.Front()
	task = e.Value.(*TaskEntity)
	delete(tq.TaskTable, task.TaskId)
	tq.TaskList.Remove(e)

	tq.Rwlock.Unlock()
	return
}

func (tq *TaskQueue) GetAllTasks( ) (tasks []*TaskEntity) {
	if tq.TaskList.Len() == 0 {
		return nil
	}
	tq.Rwlock.RLock()
	for e := tq.TaskList.Front(); e != nil; e = e.Next(){
		task := e.Value.(*TaskEntity)
		tasks = append(tasks, task)
	}
	tq.Rwlock.RUnlock()
	return
}

func (tq *TaskQueue) GetTasksByStatus(statusDemand int32) (tasks []*TaskEntity) {
	if tq.TaskList.Len() == 0 {
		return nil
	}
	tq.Rwlock.RLock()
	for e := tq.TaskList.Front(); e != nil; e = e.Next(){
		task := e.Value.(*TaskEntity)
		if task.Status == statusDemand {
			tasks = append(tasks, task)
		}
	}
	tq.Rwlock.RUnlock()
	return
}

func (tq *TaskQueue) DequeueAllTasks( ) (tasks []*TaskEntity) {
	if tq.TaskList.Len() == 0 {
		return nil
	}
	tq.Rwlock.Lock()
	for e := tq.TaskList.Front(); e != nil; {
		next_e := e.Next()

		task := e.Value.(*TaskEntity)
		tasks = append(tasks, task)

		tq.TaskList.Remove(e)
		delete(tq.TaskTable, task.TaskId)

		e = next_e
	}

	tq.Rwlock.Unlock()
	return
}

func (tq *TaskQueue) FindTask(TaskId int32)  (task *TaskEntity) {
	tq.Rwlock.RLock()
	if e, ok := tq.TaskTable[TaskId]; ok {
		task = e.Value.(*TaskEntity)
	} else {
		task = nil
	}
	tq.Rwlock.RUnlock()
	return
}

func (tq *TaskQueue) GettaskNum() (len int) {
	tq.Rwlock.RLock()
	len = tq.TaskList.Len()
	tq.Rwlock.RUnlock()
	return len
}

//clean mtq and add all its tasks to tq
//不涉及task里面修改
func (tq *TaskQueue) MergeTasks(mtq *TaskQueue) {
	tasks := mtq.DequeueAllTasks()
	if len(tasks) == 0 {
		return
	}

	tq.Rwlock.Lock()
	for _, task := range tasks {
		e := tq.TaskList.PushBack(task)
		tq.TaskTable[task.TaskId] = e
	}
	tq.Rwlock.Unlock()
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


type TaskSimpleQueue struct {
	tq []*TaskEntity
	rwlock sync.RWMutex
}

func NewTaskSimpleQueue()  *TaskSimpleQueue {
	return &TaskSimpleQueue{
		tq: make([]*TaskEntity, 0),
	}
}

func (tsq *TaskSimpleQueue) EnqueueTask(task *TaskEntity) {
	tsq.rwlock.Lock()
	tsq.tq = append(tsq.tq, task)
	tsq.rwlock.Unlock()
}

type SchedulerFunc func(tsq *TaskSimpleQueue)

func FindScheduler(schedulerID int)  (f SchedulerFunc) {
	switch schedulerID {
	case Scheduler_Default:
		f = DefaultScheduler
	case Scheduler_EDF:
		f = DeadlineFirstScheduler
	default:
		log.Fatalf("Wrong Scheduler ID!")
	}
	return  f
}

func DefaultScheduler(tsq *TaskSimpleQueue) {
	//First come first serve; do nothing
}

func DeadlineFirstScheduler(tsq *TaskSimpleQueue) {
	tsq.rwlock.Lock()
	sort.Slice(tsq.tq, func(i,j int) bool {
		return tsq.tq[i].Deadline < tsq.tq[j].Deadline
	})
	tsq.rwlock.Unlock()
}

func (tsq *TaskSimpleQueue) PopFirstTask(f SchedulerFunc) (task *TaskEntity) {
	if f != nil {
		f(tsq)
	}
	if len(tsq.tq) == 0 {
		return nil
	}
	tsq.rwlock.Lock()
	task = tsq.tq[0]
	tsq.tq = tsq.tq[1:]
	tsq.rwlock.Unlock()
	log.Printf("schedule task(%d) with earliest deadline(%d)", task.TaskId, task.Deadline)
	return task
}

func (tsq *TaskSimpleQueue) Clear() {
	if len(tsq.tq) > 0 {
		tsq.tq = tsq.tq[0:0]
	}
}

func (tsq *TaskSimpleQueue) GetTaskNum() (taskNum int){
	tsq.rwlock.RLock()
	taskNum = len(tsq.tq)
	tsq.rwlock.RUnlock()
	return taskNum
}