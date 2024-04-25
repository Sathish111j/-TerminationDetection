package main

import (
    "fmt"
    "DS_case_study/graph"
    "DS_case_study/termination"
    "time"
)

func main() {
    fmt.Println("Starting application...")
    
    g := graph.NewGraph()
    g.Initialize() // Populate graph with nodes and edges

    n1 := graph.NewNode(1)
    n2 := graph.NewNode(2)
    n3 := graph.NewNode(3)

    g.AddNode(n1)
    g.AddNode(n2)
    g.AddNode(n3)

    d := termination.NewDetector()

    simulateNodeActivity(n1, d)
    simulateNodeActivity(n2, d)
    simulateNodeActivity(n3, d)

    time.Sleep(5 * time.Second)

    if d.CheckTermination() {
        fmt.Println("All nodes have completed their tasks.")
    } else {
        fmt.Println("Some nodes are still active.")
    }
}

func simulateNodeActivity(n *graph.Node, d *termination.Detector) {
    go func() {
        fmt.Printf("Node %d is now active.\n", n.ID)
        d.RegisterActive()

        time.Sleep(time.Duration(n.ID) * time.Second)

        fmt.Printf("Node %d has completed its work and is now inactive.\n", n.ID)
        d.RegisterInactive()
    }()
}
