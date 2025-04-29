package field

import (
	"container/heap"
	"math"
)

// PathNode represents a node in the A* pathfinding algorithm
type PathNode struct {
	cell      *Cell
	parent    *PathNode
	gCost     float64 // Cost from start to this node
	hCost     float64 // Heuristic cost from this node to end
	fCost     float64 // Total cost (gCost + hCost)
	heapIndex int     // Index in the priority queue
}

// PathNodeHeap is a min-heap of PathNodes
type PathNodeHeap []*PathNode

// Len returns the length of the heap
func (h PathNodeHeap) Len() int {
	return len(h)
}

// Less compares two nodes by their fCost
func (h PathNodeHeap) Less(i, j int) bool {
	if h[i].fCost == h[j].fCost {
		// If fCosts are equal, prefer the one with lower hCost
		return h[i].hCost < h[j].hCost
	}
	return h[i].fCost < h[j].fCost
}

// Swap swaps two nodes in the heap
func (h PathNodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].heapIndex = i
	h[j].heapIndex = j
}

// Push adds a node to the heap
func (h *PathNodeHeap) Push(x interface{}) {
	node := x.(*PathNode)
	node.heapIndex = len(*h)
	*h = append(*h, node)
}

// Pop removes and returns the minimum node from the heap
func (h *PathNodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	node.heapIndex = -1
	*h = old[0 : n-1]
	return node
}

// Path represents a path between two points
type Path struct {
	cells []*Cell
	field *Field
}

// NewPath creates a new Path
func NewPath(cells []*Cell, field *Field) *Path {
	return &Path{
		cells: cells,
		field: field,
	}
}

// Cells returns the cells in the path
func (p *Path) Cells() []*Cell {
	return p.cells
}

// Length returns the number of cells in the path
func (p *Path) Length() int {
	return len(p.cells)
}

// Start returns the starting cell of the path
func (p *Path) Start() *Cell {
	if len(p.cells) == 0 {
		return nil
	}
	return p.cells[0]
}

// End returns the ending cell of the path
func (p *Path) End() *Cell {
	if len(p.cells) == 0 {
		return nil
	}
	return p.cells[len(p.cells)-1]
}

// Positions returns the positions in the path
func (p *Path) Positions() []Position {
	positions := make([]Position, len(p.cells))
	for i, cell := range p.cells {
		positions[i] = cell.Position()
	}
	return positions
}

// IsValid checks if the path is valid (all cells are walkable)
func (p *Path) IsValid() bool {
	for _, cell := range p.cells {
		if !cell.IsWalkable() {
			return false
		}
	}
	return true
}

// PathFinder handles pathfinding operations
type PathFinder struct {
	grid *CellGrid
}

// NewPathFinder creates a new PathFinder
func NewPathFinder(grid *CellGrid) *PathFinder {
	return &PathFinder{
		grid: grid,
	}
}

// FindPath finds a path between two positions using the A* algorithm
func (pf *PathFinder) FindPath(start, end Position) *Path {
	startCell := pf.grid.GetCellAtPosition(start)
	endCell := pf.grid.GetCellAtPosition(end)

	if startCell == nil || endCell == nil {
		return nil
	}

	// If start and end are the same, return a path with just that cell
	if startCell == endCell {
		return NewPath([]*Cell{startCell}, pf.grid.Field())
	}

	// If the end cell is not walkable, find the closest walkable cell
	if !endCell.IsWalkable() {
		closestPos := pf.grid.Field().ClosestWalkableSpot(end, 5)
		if closestPos == nil {
			return nil
		}
		endCell = pf.grid.GetCellAtPosition(*closestPos)
	}

	// Initialize the open and closed sets
	openSet := make(map[*Cell]*PathNode)
	closedSet := make(map[*Cell]bool)

	// Create the start node
	startNode := &PathNode{
		cell:   startCell,
		parent: nil,
		gCost:  0,
		hCost:  calculateHeuristic(startCell, endCell),
	}
	startNode.fCost = startNode.gCost + startNode.hCost

	// Initialize the priority queue
	openHeap := &PathNodeHeap{startNode}
	heap.Init(openHeap)
	openSet[startCell] = startNode

	// Main A* loop
	for openHeap.Len() > 0 {
		// Get the node with the lowest fCost
		currentNode := heap.Pop(openHeap).(*PathNode)
		currentCell := currentNode.cell

		// Remove from open set
		delete(openSet, currentCell)

		// Add to closed set
		closedSet[currentCell] = true

		// If we've reached the end, reconstruct the path
		if currentCell == endCell {
			return reconstructPath(currentNode, pf.grid.Field())
		}

		// Check all neighbors
		for _, neighbor := range currentCell.WalkableNeighbors() {
			// Skip if already processed
			if closedSet[neighbor] {
				continue
			}

			// Calculate the cost to reach this neighbor
			moveCost := 1.0

			// Diagonal movement costs more
			if neighbor.X() != currentCell.X() && neighbor.Y() != currentCell.Y() {
				moveCost = 1.414 // sqrt(2)
			}

			// Add weight factor
			moveCost *= float64(neighbor.Weight() + 1) // +1 to avoid zero weight

			// Calculate the new gCost
			newGCost := currentNode.gCost + moveCost

			// Check if this is a better path
			neighborNode, inOpenSet := openSet[neighbor]
			if !inOpenSet || newGCost < neighborNode.gCost {
				// Create or update the neighbor node
				if !inOpenSet {
					neighborNode = &PathNode{
						cell:  neighbor,
						hCost: calculateHeuristic(neighbor, endCell),
					}
					openSet[neighbor] = neighborNode
					heap.Push(openHeap, neighborNode)
				}

				// Update the node
				neighborNode.parent = currentNode
				neighborNode.gCost = newGCost
				neighborNode.fCost = neighborNode.gCost + neighborNode.hCost

				// Update the node in the heap
				if inOpenSet {
					heap.Fix(openHeap, neighborNode.heapIndex)
				}
			}
		}
	}

	// No path found
	return nil
}

// calculateHeuristic calculates the heuristic cost between two cells
func calculateHeuristic(from, to *Cell) float64 {
	// Manhattan distance
	dx := float64(abs(from.X() - to.X()))
	dy := float64(abs(from.Y() - to.Y()))

	// Use diagonal distance
	return math.Max(dx, dy) + (math.Sqrt(2)-1)*math.Min(dx, dy)
}

// reconstructPath reconstructs the path from the end node
func reconstructPath(endNode *PathNode, field *Field) *Path {
	// Count the path length
	length := 0
	for node := endNode; node != nil; node = node.parent {
		length++
	}

	// Create the path in reverse order
	cells := make([]*Cell, length)
	node := endNode
	for i := length - 1; i >= 0; i-- {
		cells[i] = node.cell
		node = node.parent
	}

	return NewPath(cells, field)
}

// FindClientPath finds a path using the client's pathfinding algorithm
// This is a simplified version that mimics the client's behavior
func (pf *PathFinder) FindClientPath(start, end Position) *Path {
	// Check if we can move directly
	if pf.grid.Field().CheckLOS(start, end, false) {
		// Direct path is possible
		startCell := pf.grid.GetCellAtPosition(start)
		endCell := pf.grid.GetCellAtPosition(end)

		if startCell != nil && endCell != nil {
			return NewPath([]*Cell{startCell, endCell}, pf.grid.Field())
		}
	}

	// Fall back to A* pathfinding
	return pf.FindPath(start, end)
}

// CheckPathFree checks if there are no obstacles in a given path
func (pf *PathFinder) CheckPathFree(start, end Position) bool {
	return pf.grid.Field().CheckLOS(start, end, false)
}

// CanAttack checks if an attack is possible between two positions
func (pf *PathFinder) CanAttack(start, end Position, canSnipe bool, attackRange int, clientSight bool) int {
	// Calculate distance
	dx := abs(start.X - end.X)
	dy := abs(start.Y - end.Y)
	distance := int(math.Max(float64(dx), float64(dy)))

	// Check range
	if distance > attackRange {
		return 0 // Out of range
	}

	// Check line of sight
	if !pf.grid.Field().CheckLOS(start, end, canSnipe) {
		return -1 // No line of sight
	}

	// If clientSight is true, we need to check if the client would consider this in range
	if clientSight {
		// Client uses a different distance calculation for some skills
		clientDistance := dx + dy
		if clientDistance > attackRange*2 {
			return 0 // Out of range according to client
		}
	}

	return 1 // Success
}
