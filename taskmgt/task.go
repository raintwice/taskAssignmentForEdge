package taskmgt

import (
	"errors"
	"os/exec"
	"reflect"
)

/*
* Task
 */

const (
	TaskStatusCode_Created = iota
	TaskStatusCode_WaitForAssign
	TaskStatusCode_Assigned
	TaskStatusCode_WaitForExec
	TaskStatusCode_Running
	TaskStatusCode_Success
	TaskStatusCode_Failed
)

type TaskEntity struct {
	TaskId int32 `json:"TaskId"`
	TaskName string `json:"TaskName"`          //任务文件名
	TaskLocation string `json:"TaskLocation"` //任务路径

	//执行状态
	Status int
	Err error

	//节点
	NodeIP  string

	//任务执行后的回调函数
	callback interface{}
	arg []interface{}

	//其他属性
}

//创建任务
func CreateTask(TaskId int32) *TaskEntity {
	task := new(TaskEntity)
	task.TaskId = TaskId
	task.Status = TaskStatusCode_Created
	return task
}

func (t *TaskEntity) Execute() error {
	if t.TaskName == "" {
		t.Err = errors.New("Empty task file name")
		t.Status = TaskStatusCode_Failed
		return t.Err
	}
	cmd := exec.Command(t.TaskLocation+"/"+t.TaskName)
	err := cmd.Run()
	if err != nil {
		t.Err = err
		t.Status = TaskStatusCode_Failed
	}
	if t.callback != nil {
		RunCallBack(t.callback, t.arg)
	}
	return err
}

func (t *TaskEntity)SetTaskCallback(f interface{}, args ...interface{}) {
	t.callback = f
	t.arg = args
}

func RunCallBack(callback interface{}, args []interface{}) {
	v := reflect.ValueOf(callback)
	if v.Kind() != reflect.Func {
		panic("Parameter callback is not a function")
	}
	vargs := make([]reflect.Value, len(args))
	for i, arg := range args {
		vargs[i] = reflect.ValueOf(arg)
	}

	v.Call(vargs)
}