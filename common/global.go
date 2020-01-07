package common

const (
	MasterIP = "127.0.0.1"
	MasterPort  = 50051  //master端口
	NodePort  = 50052   //node端口
	MasterPortForClient = 50053
	ClientPort = 50054 //client端口

	Timeout = 2000   //心跳间隔（毫秒）
	AssgnTimeout = 1000 //调度器工作间隔（毫秒）
	)

type NodeIdentity struct {
	IP  string
	Port  int
}