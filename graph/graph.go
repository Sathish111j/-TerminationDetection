package graph

import "sort"

type Node struct {
    ID    int
    Edges []*Edge
}

type Edge struct {
    From   *Node
    To     *Node
    Weight int
}

type Graph struct {
    Nodes []*Node
    Edges []*Edge
}

func NewGraph() *Graph {
    return &Graph{}
}

// Method to add an edge
func (g *Graph) AddEdge(from, to *Node, weight int) {
    edge := &Edge{From: from, To: to, Weight: weight}
    from.Edges = append(from.Edges, edge)
    to.Edges = append(to.Edges, edge)
    g.Edges = append(g.Edges, edge)
}

// Initialize a simple graph for demonstration
func (g *Graph) Initialize() {
    // Create nodes
    for i := 1; i <= 5; i++ {
        g.Nodes = append(g.Nodes, &Node{ID: i})
    }
    // Connect nodes with edges and weights
    g.AddEdge(g.Nodes[0], g.Nodes[1], 10)
    g.AddEdge(g.Nodes[0], g.Nodes[3], 20)
    g.AddEdge(g.Nodes[1], g.Nodes[2], 30)
    g.AddEdge(g.Nodes[2], g.Nodes[4], 40)
    g.AddEdge(g.Nodes[3], g.Nodes[4], 50)
    g.AddEdge(g.Nodes[3], g.Nodes[1], 60)
}

// Union-Find structure to help with Kruskal's MST
type UnionFind struct {
    parent, rank []int
}

func NewUnionFind(size int) *UnionFind {
    uf := &UnionFind{
        parent: make([]int, size),
        rank:   make([]int, size),
    }
    for i := range uf.parent {
        uf.parent[i] = i
    }
    return uf
}

func (uf *UnionFind) find(n int) int {
    if uf.parent[n] != n {
        uf.parent[n] = uf.find(uf.parent[n]) // path compression
    }
    return uf.parent[n]
}
// AddNode adds a new node to the graph.
func (g *Graph) AddNode(node *Node) {
    g.Nodes = append(g.Nodes, node)
}

func (uf *UnionFind) union(x, y int) {
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
    }
}

// Kruskal's algorithm to find the Minimum Spanning Tree (MST)
func (g *Graph) KruskalsMST() []*Edge {
    var mst []*Edge
    uf := NewUnionFind(len(g.Nodes))

    // Sort edges by weight
    sort.Slice(g.Edges, func(i, j int) bool {
        return g.Edges[i].Weight < g.Edges[j].Weight
    })

    for _, edge := range g.Edges {
        if uf.find(edge.From.ID-1) != uf.find(edge.To.ID-1) {
            mst = append(mst, edge)
            uf.union(edge.From.ID-1, edge.To.ID-1)
        }
    }
    return mst
}
