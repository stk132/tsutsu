# tsutsu
fireworq api client

## usage

``` go
package main

import (
    "github/stk132/tsutsu"
    "log"
    "fmt"
)

func main() {
    baseURL := "http://localhost:8080" //fireworq url
    queueName := "sample_queue"
    client := tsutsu.NewTsutsu(baseURL)
    
    //create queue
    if _, err := client.CreateQueue(queueName, 20, 1); err != nil {
        log.Fatal(err)
    }
    
    //show queue
    q, err := client.Queue(queueName)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(q.QueueName) //sample_queue
    
    //show queue list
    queues, err := client.Queues()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, queue := range queues {
        fmt.Println(queue.QueueName) //default, sample_queue
    }
    
    //delete queue
    if _, err := client.DeleteQueue(queueName); err != nil {
        log.Fatal(err)
    }
    
    categoryName := "sample_category
    
    //create routing
    if _, err := client.CreateRouting(categoryName, "default") {
        log.Fatal(err)
    }
    
    //show routing
    r, err := client.Routing(categoryName)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(r.QueueName) //default
    
    //show routing list
    routings, err := client.Routings()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, routing := range routings {
        fmt.Println(routing.QueueName) //default
    }
    
    //delete routing
    if _, err := client.DeleteRouting(categoryName); err != nil {
        log.Fatal(err)
    }
}
```