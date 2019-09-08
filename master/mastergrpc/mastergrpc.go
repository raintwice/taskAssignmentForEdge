
package mastergrpc

import (
    "log"
    "net"

    "golang.org/x/net/context"
    "google.golang.org/grpc"
    pb "taskAssignment/proto"
    "taskAssignment/master/nodemgt"
)

const (
    port = ":50051"
)

type server struct {}

func (s *server) JoinGroup(ctx context.Context, in *pb.JoinRequest) (*pb.JoinReply, error) {
    nodemgt.CreateNode(in.IpAddr);
    log.Printf("node %d has joined in the group", in.IpAddr)
    return &pb.JoinReply{Reply: true}, nil
}

func (s *server) ExitGroup(ctx context.Context, in *pb.ExitRequest) (*pb.ExitReply, error) {
    nodemgt.DeleteNode(in.IpAddr);
    log.Printf("node %d has left from the group", in.IpAddr)
    return &pb.ExitReply{Reply: true}, nil
}

func (s *server) HeartBeat(ctx context.Context, in *pb.HeartBeatRequest) (*pb.HeartBeatReply, error) {
    return &pb.HeartBeatReply{Time: 10001}, nil
}

func BootupGrpcServer() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatal("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterConnectionServer(s, &server{})
    s.Serve(lis)
}
