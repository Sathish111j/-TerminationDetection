// main package
package main

import(
    "fmt"
    "time"
 "sync"
    "DS_case_study/graph"
    "DS_case_study/termination"
    "math/rand"
 )


 func simulateNodeProcessing(node *graph.Node, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Node %d starting processing.\n", node.ID)
    // Simulate random processing time
    time.Sleep(time.Duration(rand.Intn(5) + 1) * time.Second)
    node.Mutex.Lock()
    node.Completed = true
    node.Mutex.Unlock()
    fmt.Printf("Node %d finished processing.\n", node.ID)
}


func main() {
    g := graph.NewGraph()
    g.Initialize()
    g.BuildMST()

    var wg sync.WaitGroup
    for _, node := range g.Nodes {
        wg.Add(1)
        go simulateNodeProcessing(node, &wg)
    }

    wg.Wait()

    // Print the tree after building the MST
    fmt.Println("Minimum Spanning Tree:")
    termination.PrintTree(g.RootNode, "", true)
    fmt.Println()

    logger := &termination.OrderLogger{}
    detector := termination.NewDetector()
    if detector.CheckTermination(g.RootNode, logger) {
        fmt.Println("All nodes have completed their tasks.")
    } else {
        fmt.Println("Some nodes are still active.")
    }

    fmt.Println("Order of checks:", logger.Order)
}