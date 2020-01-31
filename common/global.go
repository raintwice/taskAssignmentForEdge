package common

const (
	MasterIP = "127.0.0.1"
	MasterPort  = 50051  //master端口
	NodePort  = 50052   //node端口
	MasterPortForClient = 50053
	ClientPort = 50054 //client端口

	Timeout = 2000   //心跳间隔（毫秒）
	AssignInterval = 1000 //调度器工作间隔（毫秒）
	)

type NodeIdentity struct {
	IP  string
	Port  int
}

//indicate the nodes with the same hardware and software configuration
const (
	MachineType_Simualted = iota
	MachineType_RaspPi_3B
	MachineType_RaspPi_4B
	MachineType_Surface_M3
)

//indicate the nodes with same network environment and mobility
const (
	GroupIndex_Simulated = iota
)

const PreDispatch_RR_Cnt = 5*50