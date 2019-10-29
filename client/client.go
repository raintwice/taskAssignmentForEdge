package main

import (
	"sync"
	"taskAssignmentForEdge/client/connect"
)

func main() {
	client := connect.NewClient()
	var wg sync.WaitGroup
	client.StartRecvResultServer(&wg)
	wg.Wait()
}