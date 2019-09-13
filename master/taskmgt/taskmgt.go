package taskmgt


import (
	"container/list"
	//"fmt"
	"log"
)

type taskEntity struct {
	taskId int32
	taskName string
	taskLocation string 
	//其他属性
}

var taskList list.List
var taskTable = make(map[int32] (*list.Element) )

//创建任务并加入队列
func CreateTask(taskId int32) *taskEntity {
	task := new(taskEntity)
	task.taskId = taskId
	e := taskList.PushBack(task)
	taskTable[taskId] = e
	log.Printf("entity of task(Id:%d) is created", taskId);
	return e.Value.(*taskEntity)
}

//从队列中删除任务
func DeleteTask(taskId int32) {
	if e, ok := taskTable[taskId]; ok {
		delete(taskTable, taskId)
		//n := e.Value.(*taskEntity)
		//n = nil
		taskList.Remove(e)
		log.Printf("entity of task(IP:%d) is deleted", taskId);
	}
}

//根据taskId 查找任务
func FindTask(taskId int32) *taskEntity {
	if e, ok := taskTable[taskId]; ok {
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

}
