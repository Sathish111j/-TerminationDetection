package termination

import (
  "fmt"
  "sync"
  "DS_case_study/graph"
  "time"
)

type Detector struct {
  wg sync.WaitGroup
  done chan bool
}

func NewDetector() *Detector {
  return &Detector{done: make(chan bool)}
}

func (d *Detector) PropagateCompletion(n *graph.Node) {
  fmt.Printf("Node %d checking children completion.\n", n.ID)
  for _, child := range n.Children {
    if !child.Completed { // Check the completion flag before channel operation
      d.wg.Add(1)
      go func(c *graph.Node) {
        defer d.wg.Done()
        d.PropagateCompletion(c) // Recursive call to check child completion
      }(child)
    }
  }
  // Wait for all goroutines spawned for child completion checks to finish
  d.wg.Wait()
}

func (d *Detector) CheckTermination(root *graph.Node) bool {
  d.wg.Add(1)
  go func() {
    defer d.wg.Done()
    d.PropagateCompletion(root)
    close(d.done) // Signal completion of all checks
  }()
  select {
  case <-d.done:
    return true
  case <-time.After(time.Minute * 5): // Increased timeout for debugging
    fmt.Println("Termination check timed out.")
    return false
  }
}
