package eduserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	pb "edu_paper/edupb"

	"github.com/jinzhu/copier"
)

// ErrAlreadyExists is returned when a record with the same ID already exists in the store
var ErrAlreadyExists = errors.New("record already exists")

// QSetStore is an interface to store qset
type QSetStore interface {
	//QueSet save
	Save(qset *pb.QueSet) error
	// Search searches for qset
	Search(ctx context.Context, filter *pb.Filter, found func(qset *pb.QueSet) error) error
}

// InMemoryQSetStore stores qset in memory
type InMemoryQSetStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.QueSet
}

// NewInMemoryQSetStore returns a new InMemoryQSetStore
func NewInMemoryQSetStore() *InMemoryQSetStore {
	{
		return &InMemoryQSetStore{
			data: make(map[string]*pb.QueSet),
		}
	}
}

// Save saves the qset to the store
func (store *InMemoryQSetStore) Save(qset *pb.QueSet) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[qset.PaperId] != nil {
		return ErrAlreadyExists
	}

	other, err := deepCopy(qset)
	if err != nil {
		return err
	}

	store.data[other.PaperId] = other
	return nil
}

// Search searches for question set with filter, returns one by one via the found function
func (store *InMemoryQSetStore) Search(ctx context.Context, filter *pb.Filter, found func(qset *pb.QueSet) error,
) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, qset := range store.data {
		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Print("context is cancelled")
			return nil
		}

		other, err := deepCopy(qset)
		if err != nil {
			return err
		}

		err = found(other)
		if err != nil {
			return err
		}
	}

	return nil
}
func deepCopy(qset *pb.QueSet) (*pb.QueSet, error) {
	other := &pb.QueSet{}

	err := copier.Copy(other, qset)
	if err != nil {
		return nil, fmt.Errorf("cannot copy qset data: %w", err)
	}

	return other, nil
}
