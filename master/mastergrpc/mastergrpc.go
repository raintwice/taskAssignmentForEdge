
package mastergrpc

import (
    "log"
    "net"
    //"strconv"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    pb "taskAssignmentForEdge/proto"
    "taskAssignmentForEdge/master/nodemgt"
)

const (
    masterAddr = "localhost:50051"
)

type server struct {}

func callback(ipAddr string) {
    log.Printf("Now calling back to ")
    conn, err := grpc.Dial(ipAddr , grpc.WithInsecure())
    if err != nil {
        log.Fatal("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewNodeConnectionClient(conn)
    //pb.NewConnectionClient(conn)

    r, err := c.SendFile(context.Background(), &pb.FileContent{Filename: "file111"})
    if err != nil {
        log.Fatal("could not send file: %v", err)
    }
    if(r.Status == 0) {
        log.Printf("Succeed to send file")
    } else {
        log.Printf("Failed send file")
    }
}

func (s *server) JoinGroup(ctx context.Context, in *pb.JoinRequest) (*pb.JoinReply, error) {
    nodemgt.CreateNode(in.IpAddr)
    log.Printf("node %s has joined in the group", in.IpAddr)
    callback(in.IpAddr)
    return &pb.JoinReply{Reply: true}, nil
}

func (s *server) ExitGroup(ctx context.Context, in *pb.ExitRequest) (*pb.ExitReply, error) {
    nodemgt.DeleteNode(in.IpAddr);
    log.Printf("node %s has left from the group", in.IpAddr)
    return &pb.ExitReply{Reply: true}, nil
}

func (s *server) HeartBeat(ctx context.Context, in *pb.HeartBeatRequest) (*pb.HeartBeatReply, error) {
    return &pb.HeartBeatReply{Time: 10001}, nil
}

func BootupGrpcServer() {
    lis, err := net.Listen("tcp", masterAddr)
    if err != nil {
        log.Fatal("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterConnectionServer(s, &server{})
    s.Serve(lis)
}
