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
  // Replace with your actual node processing logic
  time.Sleep(time.Duration(node.ID) * time.Second) // Simulate varying processing times
  fmt.Printf("Node %d finished processing.\n", node.ID)
  node.Completed = true  // Set the completion flag
  close(node.TaskCompleted) // Still close the channel for potential future use
}

func main() {
  g := graph.NewGraph()
  g.Initialize()
  
  
  g.AddEdge(g.Nodes[4], g.Nodes[1], 10)  

  g.BuildMST()

  detector := termination.NewDetector()

  var wg sync.WaitGroup
  for _, node := range g.Nodes {
    wg.Add(1)
    // Create a copy of the node and pass it to the goroutine
    nodeCopy := *node
    go func(n graph.Node) {
      defer wg.Done()
      n.Active = true // Set node to active before processing
      simulateNodeProcessing(&n)
      n.Active = false // Set node to inactive after processing
    }(nodeCopy)
  }

  // Wait for all processing goroutines to finish before checking termination
  wg.Wait()

  if detector.CheckTermination(g.RootNode) {
    fmt.Println("All nodes have completed their tasks.")
  } else {
    fmt.Println("Some nodes are still active.")
  }
}
