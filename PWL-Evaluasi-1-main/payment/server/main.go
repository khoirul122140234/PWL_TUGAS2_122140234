
package main

import (
    "context"
    "log"
    "math/rand"
    "net"
    "time"

    pb "create_order_saga/proto/payment"
    "google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedPaymentServiceServer
}

func (s *server) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
    log.Printf("Processing payment for Order ID: %s, Amount: %.2f", req.OrderId, req.Amount)

    // Simulasi sukses/gagal pembayaran (acak)
    status := "SUCCESS"
    if rand.Intn(100) < 30 { // 30% kemungkinan gagal
        status = "FAILED"
    }

    return &pb.ProcessPaymentResponse{
        PaymentId: generatePaymentID(),
        Status:    status,
    }, nil
}

func (s *server) RefundPayment(ctx context.Context, req *pb.RefundPaymentRequest) (*pb.RefundPaymentResponse, error) {
    log.Printf("Refunding payment ID: %s", req.PaymentId)
    return &pb.RefundPaymentResponse{Status: "REFUNDED"}, nil
}

func generatePaymentID() string {
    rand.Seed(time.Now().UnixNano())
    return "PAY" + time.Now().Format("150405") + string(rand.Intn(1000))
}

func main() {
    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    grpcServer := grpc.NewServer()
    pb.RegisterPaymentServiceServer(grpcServer, &server{})
    log.Println("Payment Service running on port 50052...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
