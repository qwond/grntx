// Dumb gRPC client
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	v1 "github.com/qwond/grntx/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	pair       = flag.String("pair", "usdtrub", "Currency pair to get rates for")
	interval   = flag.Duration("interval", 10*time.Second, "Polling interval")
)

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := v1.NewRatesServiceClient(conn)

	// Run health check first
	if err := checkHealth(client); err != nil {
		log.Fatalf("Health check failed: %v", err)
	}

	// Start polling for rates
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := getRates(client, *pair); err != nil {
				log.Printf("Error getting rates: %v", err)
			}
		}
	}
}

func checkHealth(client v1.RatesServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.HealthCheck(ctx, &v1.HealthCheckRequest{})
	if err != nil {
		return fmt.Errorf("could not check health: %v", err)
	}

	if resp.Status != "SERVING" {
		return fmt.Errorf("service is not healthy: %s", resp.Status)
	}

	log.Println("Service is healthy")
	return nil
}

func getRates(client v1.RatesServiceClient, pair string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetRates(ctx, &v1.GetRatesRequest{Pair: pair})
	if err != nil {
		return fmt.Errorf("could not get rates: %v", err)
	}

	log.Printf("Pair: %s, Ask: %s, Bid: %s, Timestamp: %d",
		resp.Pair, resp.Ask, resp.Bid, resp.Timestamp)

	return nil
}
