package graph

import (
  "sort"
)

type Node struct {
  ID           int
  Edges        []*Edge
  Children     []*Node
  Parent       *Node
  Active       bool // Indicates whether the node is currently processing a task
  TaskCompleted chan struct{} // Channel signaled when the node's task is complete
  Completed    bool           // Flag to indicate if processing is complete
}

type Edge struct {
  From *Node
  To   *Node
  Weight int
}

type Graph struct {
  Nodes  []*Node
  Edges  []*Edge
  RootNode *Node // Define the root node
}

func NewGraph() *Graph {
  return &Graph{}
}

func NewNode(id int) *Node {
  return &Node{
    ID:           id,
    TaskCompleted: make(chan struct{}),
    Completed:    false,
  }
}

func (g *Graph) AddEdge(from, to *Node, weight int) {
  edge := &Edge{From: from, To: to, Weight: weight}
  from.Edges = append(from.Edges, edge)
  to.Edges = append(to.Edges, edge)
  g.Edges = append(g.Edges, edge)
}

func (g *Graph) Initialize() {
  // Create nodes
  for i := 1; i <= 5; i++ {
    node := NewNode(i)
    g.Nodes = append(g.Nodes, node)
    // Set the first node as the root node
    if i == 1 {
      g.RootNode = node
    }
  }
  // Connect nodes with edges and weights (replace with your actual connections)
  g.AddEdge(g.Nodes[0], g.Nodes[1], 10)
  g.AddEdge(g.Nodes[0], g.Nodes[3], 20)
  g.AddEdge(g.Nodes[1], g.Nodes[2], 30)
  g.AddEdge(g.Nodes[2], g.Nodes[4], 40)
  g.AddEdge(g.Nodes[3], g.Nodes[4], 50)
  g.AddEdge(g.Nodes[3], g.Nodes[1], 60) // Added edge for clarity
}
// Union-Find Helper for Kruskal's MST
type UnionFind struct {
    parent []int
    rank   []int
}

func NewUnionFind(size int) *UnionFind {
    parent := make([]int, size)
    rank := make([]int, size)
    for i := range parent {
        parent[i] = i
    }
    return &UnionFind{parent: parent, rank: rank}
}

func (uf *UnionFind) find(n int) int {
    if uf.parent[n] != n {
        uf.parent[n] = uf.find(uf.parent[n])
    }
    return uf.parent[n]
}

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

// Kruskal's algorithm to build the MST (corrected)
func (g *Graph) BuildMST() {
    sort.Slice(g.Edges, func(i, j int) bool {
      return g.Edges[i].Weight < g.Edges[j].Weight
    })
  
    uf := NewUnionFind(len(g.Nodes))
    mst := make([]*Edge, 0, len(g.Nodes)-1)
  
    for _, edge := range g.Edges {
      rootX := uf.find(edge.From.ID - 1)
      rootY := uf.find(edge.To.ID - 1)
      if rootX != rootY {
        mst = append(mst, edge)
  
        // Establish parent-child relationships based on the lower ranked node
        if uf.rank[rootX] > uf.rank[rootY] {
          edge.To.Parent = edge.From
        } else if uf.rank[rootX] < uf.rank[rootY] {
          edge.From.Parent = edge.To
        } else {
          edge.To.Parent = edge.From
          uf.rank[rootX]++
        }
        edge.From.Children = append(edge.From.Children, edge.To)
      }
    }
  }