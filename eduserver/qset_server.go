package eduserver

import (
	"context"
	"log"

	pb "edu_paper/edupb"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//QSetServer server
type QSetServer struct {
	qsetStore QSetStore
}

// NewQSetServer returns a new qsetServer
func NewQSetServer(qsetStore QSetStore) *QSetServer {
	return &QSetServer{qsetStore}
}

//CreateQueSet to create paper
func (server *QSetServer) CreateQueSet(
	ctx context.Context,
	req *pb.CreateQueSetRequest,
) (*pb.CreateQueSetResponse, error) {
	qset := req.GetQueSet()
	log.Printf("receive a create-paper request with id: %s", qset.PaperId)
	if len(qset.PaperId) > 0 {
		_, err := uuid.Parse(qset.PaperId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Question ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new Questin ID: %v", err)
		}
		qset.PaperId = id.String()
	}
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	err := server.qsetStore.Save(qset)
	if err != nil {
		code := codes.Internal
		code = codes.AlreadyExists
		return nil, status.Errorf(code, "cannot save qset to the store: %v", err)
	}

	log.Printf("saved qset with id: %s", qset.PaperId)

	res := &pb.CreateQueSetResponse{
		Id: qset.PaperId,
	}
	return res, nil

}

// SearchQueSet is a server-streaming RPC to search for qsets
func (server *QSetServer) SearchQueSet(
	req *pb.SearchQueSetRequest,
	stream pb.PaperService_SearchQueSetServer,
) error {
	filter := req.GetFilter()
	log.Printf("receive a search-question paper request with filter: %v", filter)

	err := server.qsetStore.Search(
		stream.Context(),
		filter,
		func(qset *pb.QueSet) error {
			res := &pb.SearchQueSetResponse{QueSet: qset}
			err := stream.Send(res)
			if err != nil {
				return err
			}

			log.Printf("sent qset with id: %s", qset.GetPaperId())
			return nil
		},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
