package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type Item struct {
	Title string
	Body  string
}

type API struct{}

var wg = sync.WaitGroup{}
var database []Item

func (a *API) GetDB(empty string, reply *[]Item) error {
	*reply = database
	return nil
}

func (a *API) GetByName(title string, reply *Item) error {
	var getItem Item
	for _, val := range database {
		if val.Title == title {
			getItem = val
			*reply = getItem
			return nil;
		}
	}
	return fmt.Errorf("Not Found");
}

func (a *API) AddItem(item Item, reply *Item) error {
	database = append(database, item)
	*reply = item
	return nil
}

func (a *API) EditItem(item Item, reply *Item) error {
	var changed Item
	found := false;
	for idx, val := range database {
		if val.Title == item.Title {
			database[idx] = Item{item.Title, item.Body}
			changed = database[idx]
			found = true;
		}
	}

	*reply = changed
	if !found{
		return fmt.Errorf("Not Found");
	}
	return nil;
}

func (a *API) DeleteItem(item Item, reply *Item) error {
	var del Item
	found := false;

	for idx, val := range database {
		if val.Title == item.Title && val.Body == item.Body {
			database = append(database[:idx], database[idx+1:]...)
			del = item
			found = true;
			break
		}
	}

	*reply = del
	if !found{
		return fmt.Errorf("Not Found");
	}
	return nil
}

func main() {
	api := &API{}
	err := rpc.Register(api)
	if err != nil {
		log.Fatal("Rrror registering API", err)
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":4040")

	if err != nil {
		log.Fatal("Listener error", err)
	}
	log.Printf("serving rpc on port %d", 4040)
	wg.Add(1);
	go http.Serve(listener, nil)
	wg.Wait();
}