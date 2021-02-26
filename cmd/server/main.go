package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal"
	"github.com/rsachdeva/illuminatingdeposits-grpc/interestcal/interestcalpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/kafkawriter"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbconn"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth"
	"github.com/rsachdeva/illuminatingdeposits-grpc/mongodbhealth/mongodbhealthpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/readenv"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn"
	"github.com/rsachdeva/illuminatingdeposits-grpc/userauthn/userauthnpb"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt"
	"github.com/rsachdeva/illuminatingdeposits-grpc/usermgmt/usermgmtpb"
	"google.golang.org/grpc"
)

const (
	// https://stackoverflow.com/questions/64093550/grpc-server-not-working-in-docker-compose
	address  = "0.0.0.0:50052"
	topic    = "depositcalculation-grpc"
	kafkaURL = "kafka:9092"
)

func main() {
	tls := readenv.TlsEnabled()
	fmt.Println("TLS is", tls)

	log.SetFlags(log.LstdFlags | log.Ltime | log.Lshortfile)
	log.Println("Starting ServiceServer...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mt := mongodbconn.Connect(ctx, 10)
	mdb := mt.Database("depositsmongodb")
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("could not listen %v", err)
	}

	zurl := os.Getenv("DEPOSITS_TRACE_URL")
	log.Println("zurl  empty ", zurl == "")
	if zurl == "" {
		zurl = "http://zipkin:9411/api/v2/spans"
	}
	log.Println("zipkin trace url is ", zurl)
	closer, err := RegisterTracer(
		"illuminatingdeposits-grpc",
		address,
		zurl,
		1,
	)
	if err != nil {
		log.Fatalf("could not regsiter tracer %v", err)
	}
	defer func() {
		err := closer()
		if err != nil {
			log.Println("could not close reporter", err)
		}
	}()
	log.Println("Tracer registered...")

	var opts []grpc.ServerOption
	if tls {
		opts = tlsOpts()
	}
	opts = append(opts, grpc.UnaryInterceptor(userauthn.EnsureValidToken))
	s := grpc.NewServer(opts...)
	log.Println("Registering gRPC proto MongoDBHealthService...")
	mongodbhealthpb.RegisterMongoDbHealthServiceServer(s, mongodbhealth.ServiceServer{
		Mct: mongodbconn.Connect(ctx, 2),
	})
	log.Println("Registering gRPC proto UserMgmtService...")
	usermgmtpb.RegisterUserMgmtServiceServer(s, usermgmt.ServiceServer{
		Mdb: mdb,
	})
	log.Println("Registering gRPC proto UserAuthenticationService...")
	userauthnpb.RegisterUserAuthnServiceServer(s, userauthn.ServiceServer{
		Mdb: mdb,
	})
	kafkawriter.CreateTopic(kafkaURL, topic)
	kafkaWriter := kafkawriter.Configure(kafkaURL, topic)
	log.Println("kafkaWriter is ", kafkaWriter)
	defer func() {
		err := kafkaWriter.Close()
		if err != nil {
			log.Fatalln("could not close depositcalculation connection", err)
		}
	}()
	log.Println("Kafka writer registered...")
	log.Println("Registering gRPC proto InterestCalService...")
	interestcalpb.RegisterInterestCalServiceServer(s, interestcal.ServiceServer{
		KafkaWriter: kafkaWriter,
	})

	// Trace struct {
	// 	URL         string  `conf:"default:http://zipkin:9411/api/v2/spans"`
	// 	Service     string  `conf:"default:illuminatingdeposits-rest"`
	// 	Probability float64 `conf:"default:1"`
	// }

	serveWithShutdown(ctx, s, lis, mt)
}
