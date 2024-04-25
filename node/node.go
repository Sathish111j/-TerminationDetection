package node

import "fmt"

// Status type defines the possible states of a Node
type Status int

// These constants represent the status of the Node
const (
    Idle Status = iota
    Active
)

// Node represents a node in the network with an ID and a State
type Node struct {
    ID    int
    State Status
    // other properties can be added here such as links to other nodes, data, etc.
}

// NewNode creates a new node with the given ID, defaulting to Idle state
func NewNode(id int) *Node {
    return &Node{
        ID:    id,
        State: Idle,
    }
}

// SetState changes the state of the node
func (n *Node) SetState(state Status) {
    n.State = state
    fmt.Printf("Node %d state changed to %v\n", n.ID, state)
}

// Activate sets the node's state to Active
func (n *Node) Activate() {
    n.SetState(Active)
}

// Deactivate sets the node's state to Idle
func (n *Node) Deactivate() {
    n.SetState(Idle)
}

// Example of additional functionality:
// SendMessage simulates sending a message from this node to another
func (n *Node) SendMessage(to *Node, message string) {
    fmt.Printf("Node %d sending message to Node %d: %s\n", n.ID, to.ID, message)
    to.ReceiveMessage(n, message)
}

// ReceiveMessage simulates this node receiving a message from another node
func (n *Node) ReceiveMessage(from *Node, message string) {
    fmt.Printf("Node %d received message from Node %d: %s\n", n.ID, from.ID, message)
    // Handle message based on its content or simply acknowledge it
}

// Additional methods that define behavior of the node can be added below
