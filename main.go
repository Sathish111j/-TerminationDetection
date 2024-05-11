package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
    "DS_case_study/graph" 
)

func simulateNodeProcessing(node *graph.Node, wg *sync.WaitGroup, g *graph.Graph) {
    defer wg.Done()
    fmt.Printf("Node %d starting processing.\n", node.ID)

    // Simulate random processing time
    time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

    node.Mutex.Lock()
    node.Completed = true
    node.Mutex.Unlock()

    // Simulate message passing
    if rand.Intn(2) == 0 {
        node.Mutex.Lock()
        node.Color = "black"
        fmt.Printf("Node %d turns black due to its own processing logic.\n", node.ID)
        node.Mutex.Unlock()

        // Randomly select a neighbor to send a message
        neighbors := make([]*graph.Node, 0)
        for _, edge := range node.Edges {
            if edge.To != node && edge.To != node.Parent { // Ensure we don't send a message to self or parent
                neighbors = append(neighbors, edge.To)
            }
        }
        if len(neighbors) > 0 {
            target := neighbors[rand.Intn(len(neighbors))]
            fmt.Printf("Node %d sending message to Node %d.\n", node.ID, target.ID)
            target.Mutex.Lock()
            if target.Color == "white" {
                target.Color = "black" // Change target's color to black on receiving a message
                fmt.Printf("Node %d turns black upon receiving a message from Node %d.\n", target.ID, node.ID)
            }
            target.Mutex.Unlock()
        }
    }

    fmt.Printf("Node %d finished processing.\n", node.ID)
}

func main() {
    g := graph.NewGraph()
    g.Initialize()
    g.BuildMST()

    var wg sync.WaitGroup
    for _, node := range g.Nodes {
        wg.Add(1)
        go simulateNodeProcessing(node, &wg, g)
    }

    wg.Wait() // Wait for all nodes to finish processing

    // Print the list of black nodes after all message passing is completed
    g.PrintBlackNodes()

    fmt.Println("Simulation complete. Final state of the graph:")
    if g.DetectTermination() {
        fmt.Println("Termination detected.")
    } else {
        fmt.Println("Termination detection failed.")
    }
}
