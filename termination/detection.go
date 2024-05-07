// termination package
package termination

import (
	"fmt"
	"strings"
	"sync"

	"DS_case_study/graph"
)

type Detector struct {
	wg sync.WaitGroup // Use a single WaitGroup for all nodes
}

func NewDetector() *Detector {
	return &Detector{}
}

type OrderLogger struct {
	Order []int
}

func PrintTree(node *graph.Node, prefix string, isTail bool) {
	if node == nil {
		return
	}
	linePrefix := prefix
	if isTail {
		linePrefix += "└── "
	} else {
		linePrefix += "├── "
	}
	fmt.Println(linePrefix + fmt.Sprintf("Node %d", node.ID))
	for i, child := range node.Children {
		if isTail {
			PrintTree(child, prefix+"    ", i == len(node.Children)-1)
		} else {
			PrintTree(child, prefix+"│   ", i == len(node.Children)-1)
		}

	}
}

func calculatePositions(node *graph.Node, depth int, levels *[][]graph.Node, nodes map[*graph.Node]int) {
	if len(*levels) == depth {
		*levels = append(*levels, []graph.Node{})
	}
	nodes[node] = len((*levels)[depth])
	(*levels)[depth] = append((*levels)[depth], *node)
	for _, child := range node.Children {
		calculatePositions(child, depth+1, levels, nodes)
	}
}

func drawNodes(lines *[]string, level []graph.Node, nodes map[*graph.Node]int) {
	line := ""
	lastPos := 0
	for _, node := range level {
		pos := nodes[&node] * 12
		if pos > lastPos {
			line += strings.Repeat(" ", pos-lastPos)
		}
		line += fmt.Sprintf("Node %d", node.ID)
		lastPos = pos + len(fmt.Sprintf("Node %d", node.ID))
	}
	*lines = append(*lines, line)
}

func drawLines(lines *[]string, level []graph.Node, nodes map[*graph.Node]int) {
	line := ""
	lastPos := 0
	for _, node := range level {
		if node.Parent != nil {
			parentPos := nodes[node.Parent]*12 + 6
			childPos := nodes[&node]*12 + 6
			if childPos > lastPos {
				line += strings.Repeat(" ", childPos-lastPos)
			}
			if parentPos < childPos {
				line += "/"
			} else {
				line += "\\"
			}
			lastPos = childPos + 1
		}
	}
	*lines = append(*lines, line)
}

func (d *Detector) CheckTermination(root *graph.Node, logger *OrderLogger) bool {
	if root == nil {
		return true
	}

	root.Mutex.Lock()
	completed := root.Completed
	root.Mutex.Unlock()

	logger.Order = append(logger.Order, root.ID) // Log the order of checking
	fmt.Printf("Checking Node %d: %t\n", root.ID, completed)

	for _, child := range root.Children {
		if !d.CheckTermination(child, logger) {
			return false
		}
	}
	return completed
}
func checkCompletionRecursive(node *graph.Node, completed *bool, completedMutex *sync.Mutex) {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	if !node.Completed {
		completedMutex.Lock()
		*completed = false
		completedMutex.Unlock()
		return
	}

	for _, child := range node.Children {
		checkCompletionRecursive(child, completed, completedMutex)
		if !*completed {
			return
		}
	}
}
