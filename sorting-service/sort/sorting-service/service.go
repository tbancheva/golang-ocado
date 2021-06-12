package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/tbancheva/golang-ocado/sort/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const CubbyDefaultCapacity = 2

func newSortingService() gen.SortingRobotServer {
	return &sortingService{}
}

type sortingService struct {
	LoadedItems  []*gen.Item
	SelectedItem *gen.Item
	Cubbys       map[string]cubby
}

type cubby struct {
	Capacity    int
	SortedItems []*gen.Item
}

func (c *cubby) addItem(i *gen.Item) error {
	if len(c.SortedItems) >= c.Capacity {
		return errors.New("can't add item because cubby is full")
	}
	c.SortedItems = append(c.SortedItems, i)
	return nil
}

//LoadItems loads an input array of items in the service. E.g. ["tomatoes", "cucumber", "potato", "cheese"]
func (s *sortingService) LoadItems(ctx context.Context, r *gen.LoadItemsRequest) (*gen.LoadItemsResponse, error) {
	if r.Items == nil || len(r.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "loading items called with empty cargo")
	}

	s.LoadedItems = append(s.LoadedItems, r.Items...)
	log.Printf("Successfully loaded items: %v\n", s.LoadedItems)

	return &gen.LoadItemsResponse{}, nil
}

//SelectItem chooses an item at random from the remaining ones in the array. E.g. choose "tomatoes" at random && remove item from existing array
func (s *sortingService) SelectItem(ctx context.Context, req *gen.SelectItemRequest) (res *gen.SelectItemResponse, err error) {
	if s.SelectedItem != nil {
		return nil, errors.New("an item has already been selected")
	}

	if s.LoadedItems == nil || len(s.LoadedItems) == 0 {
		return nil, errors.New("can't select an item because cargo is empty")
	}

	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	randomIndex := r.Intn(len(s.LoadedItems))
	log.Printf("Randomly selected index: %v\n", randomIndex)

	s.SelectedItem = s.LoadedItems[randomIndex]
	res = &gen.SelectItemResponse{Item: s.SelectedItem}
	log.Printf("Selected item is: %v\n", s.SelectedItem)

	s.LoadedItems = append(s.LoadedItems[:randomIndex], s.LoadedItems[randomIndex+1:]...)
	log.Printf("Loaded items after selection %v\n", s.LoadedItems)

	return res, nil
}

//MoveItem moves the selected item in the input cubby. Simply return "Success" here.
func (s *sortingService) MoveItem(ctx context.Context, r *gen.MoveItemRequest) (*gen.MoveItemResponse, error) {
	if s.SelectedItem == nil {
		return nil, errors.New("no item has been selected")
	}

	if s.Cubbys == nil {
		s.Cubbys = make(map[string]cubby)
	}

	if cubbyEl, ok := s.Cubbys[r.Cubby.Id]; ok {
		err := cubbyEl.addItem(s.SelectedItem)
		if err != nil {
			return nil, err
		}
		s.Cubbys[r.Cubby.Id] = cubbyEl
	} else {
		s.Cubbys[r.Cubby.Id] = cubby{Capacity: CubbyDefaultCapacity, SortedItems: []*gen.Item{s.SelectedItem}}
	}

	log.Printf("Moved item:%v to cubby: %v\n", s.SelectedItem.Code, r.Cubby.Id)
	s.SelectedItem = nil

	log.Printf("Total sorted items: %v\n", s.Cubbys)

	return &gen.MoveItemResponse{}, nil
}
