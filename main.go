package main

import (
  "fmt"
  "sync"
  
  "DS_case_study/graph"
  "DS_case_study/termination"
)

func simulateNodeProcessing(node *graph.Node) {
    fmt.Printf("Node %d starting processing.\n", node.ID)
    // Mock processing logic (doesn't take any time)
    for i := 0; i < 1000; i++ { // Simulate some work, but doesn't take actual time
      _ = i * i
    }
    fmt.Printf("Node %d finished processing.\n", node.ID)
    node.Completed = true
    close(node.TaskCompleted)
  }
  

func main() {
  g := graph.NewGraph()
  g.Initialize()
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
