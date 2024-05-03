// main package
package main

import (
    "fmt"
    "sync"
    "time"
    "DS_case_study/graph"
    "DS_case_study/termination"
)

func simulateNodeProcessing(node *graph.Node) {
    fmt.Printf("Node %d starting processing.\n", node.ID)
    time.Sleep(time.Duration(node.ID) * time.Second)
    fmt.Printf("Node %d finished processing.\n", node.ID)
    node.Mutex.Lock()
    node.Completed = true
    node.Mutex.Unlock()
    close(node.TaskCompleted)
}

func main() {
    g := graph.NewGraph()
    g.Initialize()
    g.BuildMST()

    var wg sync.WaitGroup
    for _, node := range g.Nodes {
        wg.Add(1)
        go func(n *graph.Node) {
            defer wg.Done()
            simulateNodeProcessing(n)
        }(node)
    }

    wg.Wait()

    detector := termination.NewDetector()
    if detector.CheckTermination(g.RootNode) {
        fmt.Println("All nodes have completed their tasks.")
    } else {
        fmt.Println("Some nodes are still active.")
    }
}
