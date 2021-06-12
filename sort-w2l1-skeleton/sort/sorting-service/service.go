package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/tbancheva/golang-ocado/sort/gen"
)

func newSortingService() gen.SortingRobotServer {
	return &sortingService{}
}

type sortingService struct {
	LoadedItems  []*gen.Item
	SelectedItem *gen.Item
	MovedItems   map[string]*gen.Item
}

//LoadItems loads an input array of items in the service. E.g. ["tomatoes", "cucumber", "potato", "cheese"]
func (s *sortingService) LoadItems(ctx context.Context, r *gen.LoadItemsRequest) (*gen.LoadItemsResponse, error) {
	s.LoadedItems = append(s.LoadedItems, r.Items...)
	fmt.Printf("Loaded items %v\n", s.LoadedItems)

	return &gen.LoadItemsResponse{}, nil
}

//SelectItem chooses an item at random from the remaining ones in the array. E.g. choose "tomatoes" at random && remove item from existing array
func (s *sortingService) SelectItem(context.Context, *gen.SelectItemRequest) (res *gen.SelectItemResponse, err error) {
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	randomIndex := r.Intn(len(s.LoadedItems))
	fmt.Printf("Randomly selected index: %v\n", randomIndex)

	res = &gen.SelectItemResponse{Item: s.LoadedItems[randomIndex]}
	fmt.Printf("Selected item is: %v\n", s.LoadedItems[randomIndex])

	s.LoadedItems = append(s.LoadedItems[:randomIndex], s.LoadedItems[randomIndex+1:]...)
	fmt.Printf("Loaded items after selection %v\n", s.LoadedItems)

	return res, nil
}

//MoveItem moves the selected item in the input cubby. Simply return "Success" here.
func (s *sortingService) MoveItem(context.Context, *gen.MoveItemRequest) (*gen.MoveItemResponse, error) {

	return nil, errors.New("not implemented")
}
