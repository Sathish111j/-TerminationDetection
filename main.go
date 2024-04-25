package main

import (
    "fmt"
    "DS_case_study/graph"
    "DS_case_study/node"
    "DS_case_study/termination"
    "time"
)

func main() {
    fmt.Println("Starting application...")
    
    // Initialize the graph
    g := graph.NewGraph()
    g.Initialize()

    // Create and add a few nodes
    n1 := node.NewNode(1)
    n2 := node.NewNode(2)
    n3 := node.NewNode(3)
    g.AddNode(n1)
    g.AddNode(n2)
    g.AddNode(n3)

    // Initialize the termination detector
    d := termination.NewDetector()

    // Simulate node activity and termination detection
    simulateNodeActivity(n1, d)
    simulateNodeActivity(n2, d)
    simulateNodeActivity(n3, d)

    // Wait for all node activities to finish
    time.Sleep(5 * time.Second)

    // Check system termination
    if d.CheckTermination() {
        fmt.Println("All nodes have completed their tasks.")
    } else {
        fmt.Println("Some nodes are still active.")
    }
}

// simulateNodeActivity simulates a node performing some work and then turning inactive
func simulateNodeActivity(n *node.Node, d *termination.Detector) {
    go func() {
        fmt.Printf("Node %d is now active.\n", n.ID)
        d.RegisterActive()

        // Simulating work by sleeping
        time.Sleep(time.Duration(n.ID) * time.Second)

        // Node completes its work
        fmt.Printf("Node %d has completed its work and is now inactive.\n", n.ID)
        d.RegisterInactive()
    }()
}