package common

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"
)

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func RemoveTaskFile(taskFileName string) {
	time.Sleep(time.Second)
	//remove the task file
	taskFileDir := Task_File_Dir + "/" + taskFileName
	cmdPara := fmt.Sprintf("rm %s", taskFileDir)
	//log.Printf("%s", cmdPara)
	cmd := exec.Command("/bin/bash", "-c", cmdPara)
	err := cmd.Run()
	if err != nil {
		log.Printf("Error, cannot rm %s due to %v", taskFileName, err)
	}
}

func RemoveAllTaskFile() {
	time.Sleep(time.Second)
	//remove the task file
	taskFileDir := Task_File_Dir + "/*"
	cmdPara := fmt.Sprintf("rm %s", taskFileDir)
	//log.Printf("%s", cmdPara)
	cmd := exec.Command("/bin/bash", "-c", cmdPara)
	err := cmd.Run()
	if err != nil {
		log.Printf("Error, cannot rm these task files due to %v", err)
	}
}