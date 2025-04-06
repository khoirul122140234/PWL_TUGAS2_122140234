
package main

import (
    "context"
    "log"
    "math/rand"
    "net"
    "time"

    pb "create_order_saga/proto/shipping"
    "google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedShippingServiceServer
}

func (s *server) StartShipping(ctx context.Context, req *pb.StartShippingRequest) (*pb.StartShippingResponse, error) {
    log.Printf("Starting shipping for Order ID: %s", req.OrderId)

    // Simulasi sukses/gagal pengiriman (acak)
    status := "SHIPPED"
    if rand.Intn(100) < 20 { // 20% kemungkinan gagal
        status = "CANCELLED"
    }

    return &pb.StartShippingResponse{
        ShipmentId: generateShipmentID(),
        Status:     status,
    }, nil
}

func (s *server) CancelShipping(ctx context.Context, req *pb.CancelShippingRequest) (*pb.CancelShippingResponse, error) {
    log.Printf("Cancelling shipment ID: %s", req.ShipmentId)
    return &pb.CancelShippingResponse{Status: "CANCELLED"}, nil
}

func generateShipmentID() string {
    rand.Seed(time.Now().UnixNano())
    return "SHP" + time.Now().Format("150405") + string(rand.Intn(1000))
}

func main() {
    lis, err := net.Listen("tcp", ":50053")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    grpcServer := grpc.NewServer()
    pb.RegisterShippingServiceServer(grpcServer, &server{})
    log.Println("Shipping Service running on port 50053...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
