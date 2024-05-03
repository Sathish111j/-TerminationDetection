package termination

import (
    "fmt"
    "sync"
    "DS_case_study/graph"
)

type Detector struct {
    wg sync.WaitGroup
}

func NewDetector() *Detector {
    return &Detector{}
}

func safeClose(ch chan struct{}) {
    defer func() {
        if recover() != nil {
            fmt.Println("Recovered from closing closed channel.")
        }
    }()
    close(ch)
}

func (d *Detector) PropagateCompletion(n *graph.Node) {
    fmt.Printf("Node %d checking children completion.\n", n.ID)
    for _, child := range n.Children {
        <-child.TaskCompleted // Wait for the task completed signal from the child
        d.PropagateCompletion(child) // Recursive call to check child completion
    }
    n.Mutex.Lock()
    if !n.Completed {
        n.Completed = true // Mark as completed
        safeClose(n.TaskCompleted) // Safely close the channel
    }
    n.Mutex.Unlock()
}

func (d *Detector) CheckTermination(root *graph.Node) bool {
    d.wg.Add(1)
    go func() {
        defer d.wg.Done()
        d.PropagateCompletion(root)
    }()
    d.wg.Wait() // Wait for all checks to complete
    return true
}
