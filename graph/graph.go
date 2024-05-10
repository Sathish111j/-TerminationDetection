package graph

import (
    "fmt"
    "sort"
    "sync"
)

// Node struct enhanced with repeat signal handling and token tracking
type Node struct {
    ID        int
    Edges     []*Edge
    Children  []*Node
    Parent    *Node
    Completed bool
    Mutex     sync.Mutex
    Color     string
    Token     *Token
}

type Token struct {
    Color string
}

type Edge struct {
    From   *Node
    To     *Node
    Weight int
}

type Graph struct {
    Nodes    []*Node
    Edges    []*Edge
    RootNode *Node
    SetS     map[int]*Node // Map to track nodes with tokens
}

func NewGraph() *Graph {
    return &Graph{
        SetS: make(map[int]*Node),
    }
}

func NewNode(id int) *Node {
    return &Node{
        ID:    id,
        Color: "white",
    }
}

// Initialize graph with nodes and edges
func (g *Graph) Initialize() {
    for i := 1; i <= 10; i++ {
        node := NewNode(i)
        g.Nodes = append(g.Nodes, node)
        if i == 1 {
            g.RootNode = node
        }
    }

    g.AddEdge(g.Nodes[0], g.Nodes[1], 10)
    g.AddEdge(g.Nodes[0], g.Nodes[3], 20)
    g.AddEdge(g.Nodes[1], g.Nodes[2], 30)
    g.AddEdge(g.Nodes[2], g.Nodes[4], 40)
    g.AddEdge(g.Nodes[3], g.Nodes[4], 50)
    g.AddEdge(g.Nodes[1], g.Nodes[5], 15)
    g.AddEdge(g.Nodes[5], g.Nodes[6], 10)
    g.AddEdge(g.Nodes[6], g.Nodes[7], 10)
    g.AddEdge(g.Nodes[7], g.Nodes[8], 10)
    g.AddEdge(g.Nodes[8], g.Nodes[9], 10)
}


// Enhanced edge adding and token propagation logic
func (g *Graph) AddEdge(from, to *Node, weight int) error {
    if from == nil || to == nil {
        return fmt.Errorf("cannot add edge with nil from or to node")
    }
    edge := &Edge{From: from, To: to, Weight: weight}
    from.Edges = append(from.Edges, edge)
    to.Edges = append(to.Edges, edge)
    g.Edges = append(g.Edges, edge)
    return nil
}

func (g *Graph) BuildMST() {
    sort.Slice(g.Edges, func(i, j int) bool {
        return g.Edges[i].Weight < g.Edges[j].Weight
    })
    uf := NewUnionFind(len(g.Nodes))
    for _, edge := range g.Edges {
        if uf.union(edge.From.ID-1, edge.To.ID-1) {
            edge.From.Children = append(edge.From.Children, edge.To)
            edge.To.Parent = edge.From
        }
    }
}

// UnionFind struct for MST construction
type UnionFind struct {
    parent []int
    rank   []int
}

// NewUnionFind constructor
func NewUnionFind(size int) *UnionFind {
    parent := make([]int, size)
    rank := make([]int, size)
    for i := range parent {
        parent[i] = i
    }
    return &UnionFind{parent: parent, rank: rank}
}

// Find method for UnionFind
func (uf *UnionFind) find(n int) int {
    if uf.parent[n] != n {
        uf.parent[n] = uf.find(uf.parent[n])
    }
    return uf.parent[n]
}

// Union method for UnionFind
func (uf *UnionFind) union(x, y int) bool {
    rootX := uf.find(x)
    rootY := uf.find(y)
    if rootX != rootY {
        if uf.rank[rootX] > uf.rank[rootY] {
            uf.parent[rootY] = rootX
        } else if uf.rank[rootX] < uf.rank[rootY] {
            uf.parent[rootX] = rootY
        } else {
            uf.parent[rootY] = rootX
            uf.rank[rootX]++
        }
        return true
    }
    return false
}

func (g *Graph) DetectTermination() bool {
    for {
        fmt.Println("Current state of the tree:")
        PrintTree(g.RootNode, "", true)
        fmt.Println()

        rootReady := true
        for _, node := range g.Nodes {
            node.Mutex.Lock()
            if node.Token != nil {
                g.SetS[node.ID] = node // Track nodes with tokens
            }
            if len(node.Children) == 0 && node.Completed {
                node.propagateToken()
            } else if node.Completed && allChildrenCompleted(node) {
                node.propagateToken()
            }
            if node == g.RootNode && node.Token != nil && node.Token.Color == "black" {
                rootReady = false
            }
            node.Mutex.Unlock()
        }

        if rootReady && g.RootNode.Token != nil {
            if g.RootNode.Token.Color == "white" {
                return true
            }
        }

        if !rootReady {
            sendRepeatSignal(g) // Send repeat signal if black token is received
            for _, node := range g.Nodes {
                node.Mutex.Lock()
                node.Token = nil
                node.Color = "white"
                node.Mutex.Unlock()
            }
            g.SetS = make(map[int]*Node) // Reset the set S
        }
    }
}

func (n *Node) propagateToken() {
    if n.Parent == nil {
        return
    }

    tokenColor := "white"
    if n.Color == "black" {
        tokenColor = "black"
    }

    n.Parent.Mutex.Lock()
    if n.Parent.Token == nil {
        n.Parent.Token = &Token{Color: tokenColor}
    } else if tokenColor == "black" {
        n.Parent.Token.Color = "black"
    }
    n.Parent.Mutex.Unlock()

    n.Color = "white"
    n.Token = nil
}

func sendRepeatSignal(g *Graph) {
    fmt.Println("Sending repeat signal from root...")
    for _, node := range g.Nodes {
        node.Mutex.Lock()
        if node.Token == nil {
            node.Token = &Token{Color: "white"} // Initiate with white token
        }
        node.Mutex.Unlock()
    }
}

func allChildrenCompleted(node *Node) bool {
    for _, child := range node.Children {
        child.Mutex.Lock()
        defer child.Mutex.Unlock()
        if !child.Completed || child.Token != nil {
            return false
        }
    }
    return true
}

func PrintTree(node *Node, prefix string, isTail bool) {
    if node == nil {
        return
    }
    linePrefix := prefix
    if isTail {
        linePrefix += "└── "
    } else {
        linePrefix += "├── "
    }
    fmt.Println(linePrefix + fmt.Sprintf("Node %d (Color: %s, Token: %v)", node.ID, node.Color, node.Token))
    for i, child := range node.Children {
        if isTail {
            PrintTree(child, prefix+"    ", i == len(node.Children)-1)
        } else {
            PrintTree(child, prefix+"│   ", i == len(node.Children)-1)
        }
    }
}