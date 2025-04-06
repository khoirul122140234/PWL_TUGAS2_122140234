package main

import (
    "context"
    "log"
    "math/rand"
    "net"
    "time"

    pb "create_order_saga/proto/order"
    "google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedOrderServiceServer
}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
    log.Printf("Creating order for user ID: %d", req.UserId)
    return &pb.CreateOrderResponse{
        OrderId: generateOrderID(),
        Status:  "PENDING",
    }, nil
}

func (s *server) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
    log.Printf("Cancelling order ID: %s", req.OrderId)
    return &pb.CancelOrderResponse{Status: "CANCELLED"}, nil
}

func generateOrderID() string {
    rand.Seed(time.Now().UnixNano())
    return "ORD" + time.Now().Format("150405") + string(rand.Intn(1000))
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    grpcServer := grpc.NewServer()
    pb.RegisterOrderServiceServer(grpcServer, &server{})
    log.Println("Order Service running on port 50051...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
