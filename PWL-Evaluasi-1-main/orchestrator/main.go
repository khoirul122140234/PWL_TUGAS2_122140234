
package main

import (
    "context"
    "log"
    "time"

    orderpb "create_order_saga/proto/order"
    paymentpb "create_order_saga/proto/payment"
    shippingpb "create_order_saga/proto/shipping"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // Connect to all services
    orderConn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    paymentConn, _ := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
    shippingConn, _ := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))

    orderClient := orderpb.NewOrderServiceClient(orderConn)
    paymentClient := paymentpb.NewPaymentServiceClient(paymentConn)
    shippingClient := shippingpb.NewShippingServiceClient(shippingConn)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Step 1: Create Order
    orderRes, err := orderClient.CreateOrder(ctx, &orderpb.CreateOrderRequest{UserId: 1})
    if err != nil {
        log.Fatalf("Failed to create order: %v", err)
    }
    log.Printf("Order created: %v", orderRes.OrderId)

    // Step 2: Process Payment
    paymentRes, err := paymentClient.ProcessPayment(ctx, &paymentpb.ProcessPaymentRequest{
        OrderId: orderRes.OrderId,
        Amount:  100.0,
    })
    if err != nil || paymentRes.Status != "SUCCESS" {
        log.Printf("Payment failed or error: %v", err)
        log.Println(">>> Compensating: Cancel Order")
        orderClient.CancelOrder(ctx, &orderpb.CancelOrderRequest{OrderId: orderRes.OrderId})
        return
    }
    log.Printf("Payment success: %v", paymentRes.PaymentId)

    // Step 3: Start Shipping
    shippingRes, err := shippingClient.StartShipping(ctx, &shippingpb.StartShippingRequest{
        OrderId: orderRes.OrderId,
    })
    if err != nil || shippingRes.Status != "SHIPPED" {
        log.Printf("Shipping failed or error: %v", err)
        log.Println(">>> Compensating: Refund Payment and Cancel Order")
        paymentClient.RefundPayment(ctx, &paymentpb.RefundPaymentRequest{PaymentId: paymentRes.PaymentId})
        orderClient.CancelOrder(ctx, &orderpb.CancelOrderRequest{OrderId: orderRes.OrderId})
        return
    }
    log.Printf("Shipping success: %v", shippingRes.ShipmentId)
    log.Println(">>> All steps completed successfully")
}
