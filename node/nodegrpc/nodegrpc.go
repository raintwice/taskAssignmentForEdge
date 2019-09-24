package nodegrpc

import (
    "log"
    //"fmt"
    "net"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    pb "taskAssignmentForEdge/proto"
)

const (
    masterAddress  = "localhost:50051"
    //port  = "50051"
)

type server struct {}

func Join(nodeaddr string) {
    conn, err := grpc.Dial(masterAddress , grpc.WithInsecure())
    if err != nil {
        log.Fatal("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewConnectionClient(conn)

    r, err := c.JoinGroup(context.Background(), &pb.JoinRequest{IpAddr: nodeaddr})
    if err != nil {
        log.Fatal("could not register: %v", err)
    }
    if(r.Reply) {
        log.Printf("Successed to join")
    } else {
        log.Printf("Failed to join")
    }
}

func (s *server) SendFile(ctx context.Context, in *pb.FileContent) (*pb.FileReply, error) {
    log.Printf("Received from master: %v", ctx, in)
    return &pb.FileReply{Status: 0}, nil
}

func BootNodeGrpc(nodeaddr string) {
    log.Printf("Node ready to listen %s", nodeaddr)
    lis, err := net.Listen("tcp", nodeaddr)
    if err != nil {
        log.Fatal("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterNodeConnectionServer(s, &server{})
    s.Serve(lis)
}

