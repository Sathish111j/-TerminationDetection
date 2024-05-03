// graph package
package graph

import (
    "sort"
    "sync"
)

// Node struct
type Node struct {
    ID            int
    Edges         []*Edge
    Children      []*Node
    Parent        *Node
    Active        bool
    TaskCompleted chan struct{}
    Completed     bool
    Mutex         sync.Mutex
    IsChannelClosed bool
}

// Edge struct
type Edge struct {
    From   *Node
    To     *Node
    Weight int
}

// Graph struct
type Graph struct {
    Nodes    []*Node
    Edges    []*Edge
    RootNode *Node
}

// NewGraph constructor
func NewGraph() *Graph {
    return &Graph{}
}

// NewNode constructor
func NewNode(id int) *Node {
    return &Node{
        ID:             id,
        TaskCompleted:  make(chan struct{}),
        Completed:      false,
        IsChannelClosed: false,
        Mutex:          sync.Mutex{},
    }
}

// AddEdge method
func (g *Graph) AddEdge(from, to *Node, weight int) {
    edge := &Edge{From: from, To: to, Weight: weight}
    from.Edges = append(from.Edges, edge)
    to.Edges = append(to.Edges, edge)
    g.Edges = append(g.Edges, edge)
}

// Initialize graph with nodes and edges
func (g *Graph) Initialize() {
    for i := 1; i <= 5; i++ {
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

// BuildMST builds a minimum spanning tree using Kruskal's algorithm
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