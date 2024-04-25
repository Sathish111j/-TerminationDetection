package termination

import (
    "fmt"
    "sync"
)

type Detector struct {
    activeCount int
    lock        sync.Mutex
}

// NewDetector creates a new Detector instance
func NewDetector() *Detector {
    return &Detector{}
}

// RegisterActive is called when a node becomes active
func (d *Detector) RegisterActive() {
    d.lock.Lock()
    defer d.lock.Unlock()
    d.activeCount++
    fmt.Println("Node became active, total active:", d.activeCount)
}

// RegisterInactive is called when a node becomes inactive
func (d *Detector) RegisterInactive() {
    d.lock.Lock()
    defer d.lock.Unlock()
    if d.activeCount > 0 {
        d.activeCount--
    }
    fmt.Println("Node became inactive, total active:", d.activeCount)
}

// CheckTermination checks if all nodes are inactive
func (d *Detector) CheckTermination() bool {
    d.lock.Lock()
    defer d.lock.Unlock()
    if d.activeCount == 0 {
        fmt.Println("All nodes are inactive. System has terminated.")
        return true
    } else {
        fmt.Println("System not terminated. Active nodes remaining:", d.activeCount)
        return false
    }
}

// Reset allows resetting the detector's count (useful for tests or reinitializations)
func (d *Detector) Reset() {
    d.lock.Lock()
    defer d.lock.Unlock()
    d.activeCount = 0
    fmt.Println("Detector reset. Active count is now zero.")
}
