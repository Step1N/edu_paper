package main

import (
	"flag"
	"log"
	"time"

	clnt "edu_paper/educlient"

	"google.golang.org/grpc"
)

const (
	username        = "admin1"
	password        = "secret"
	refreshDuration = 30 * time.Second
)

func testCreateQset(qsetClient *clnt.QSetClient) {

	println("I am on client side got server request")
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	flag.Parse()
	log.Printf("dial server %s, TLS = %t", *serverAddress, *enableTLS)
	transportOption := grpc.WithInsecure()

	cc1, err := grpc.Dial(*serverAddress, transportOption)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	qsetClient := clnt.NewQSetClient(cc1)
	testCreateQset(qsetClient)
}
