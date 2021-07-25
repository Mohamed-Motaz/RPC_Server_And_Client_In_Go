package main

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"
)

//Ideally would be from a class library serving both the
//client and the server rather than rewriting the same
//definition in both
type Item struct {
	Title string
	Body  string
}

var wg = sync.WaitGroup{};

func createItems(client *rpc.Client, ch chan struct{}){
	fmt.Println("Currently sending the items to be created over rpc");

	localWg := sync.WaitGroup{}; 
	for i := 0; i < 8; i++{
		localWg.Add(1);
		go func(i int){
			item := Item{};
			client.Call("API.AddItem", 
			Item{Title: fmt.Sprintf("Title %d", i), Body : fmt.Sprintf("Body %d", i)}, 
			&item);
			fmt.Println("This is the created item ", item);
			localWg.Done();
		}(i);
	}
	localWg.Wait();
	fmt.Println("Done sending all creation requests")
	ch <- struct{}{};
}

func getItems(client *rpc.Client, ch chan struct{}){
	fmt.Println("Currently getting the items from the db");
	wg.Add(1);
	var db []Item;
	go func(){
		client.Call("API.GetDB", "", &db);
		fmt.Println(db);
		wg.Done();
	}();
}

func main(){
	client, err := rpc.DialHTTP("tcp", "localhost:4040")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	//spawn multiple goroutines to communicate with the server
	doneCreatingItems := make(chan struct{});
	go createItems(client, doneCreatingItems);
	select{
	case <- doneCreatingItems:
		getItems(client, doneCreatingItems);
	} 

	wg.Wait();
	


}