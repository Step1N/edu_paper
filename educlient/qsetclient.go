package educlient

import (
	"context"
	pb "edu_paper/edupb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//QSetClient client for qest
type QSetClient struct {
	service pb.PaperServiceClient
}

// NewQSetClient returns a new QSet client
func NewQSetClient(cc *grpc.ClientConn) *QSetClient {
	service := pb.NewPaperServiceClient(cc)
	return &QSetClient{service}
}

// CreateQSet calls create qset RPC
func (qsetClient *QSetClient) CreateQSet(qset *pb.QueSet) {
	req := &pb.CreateQueSetRequest{
		QueSet: qset,
	}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := qsetClient.service.CreateQueSet(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			// not a big deal
			log.Print("question set already exists")
		} else {
			log.Fatal("cannot create question set: ", err)
		}
		return
	}

	log.Printf("created question set with id: %s", res.Id)
}

// SearchQueSet calls search question set RPC
func (qsetClient *QSetClient) SearchQueSet(filter *pb.Filter) {
	log.Print("search filter: ", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchQueSetRequest{Filter: filter}
	stream, err := qsetClient.service.SearchQueSet(ctx, req)
	if err != nil {
		log.Fatal("cannot search qset: ", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive response: ", err)
		}

		qset := res.GetQueSet()
		log.Print("- found: ", qset.GetPaperId())
		log.Print("  + paper name: ", qset.GetPaperName())
		log.Print("  + paper duration : ", qset.GetPaperDuration())
		log.Print("  + total number of questions ", len(qset.GetQuestions()))
		log.Print("  + paper type: ", qset.GetPaperType())
		log.Print("  + updated at: ", qset.GetUpdatedAt())
	}
}
