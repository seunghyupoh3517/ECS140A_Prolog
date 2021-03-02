package disjointset

// DisjointSet is the interface for the disjoint-set (or union-find) data
// structure.
// Do not change the definition of this interface.
// import (
// 	"fmt"
// )
type DisjointSet interface {
	// UnionSet(s, t) merges (unions) the sets containing s and t,
	// and returns the representative of the resulting merged set.
	UnionSet(int, int) int
	// FindSet(s) returns representative of the class that s belongs to.
	FindSet(int) int
}



type Collection struct{
	parent map[int]int 		// map from node -> parent
	rank map[int]int 		// map from node -> rank
}




// TODO: implement a type that satisfies the DisjointSet interface.

// NewDisjointSet creates a struct of a type that satisfies the DisjointSet interface.
func NewDisjointSet() DisjointSet {
	// create two maps and pass to collection
	var collection DisjointSet = Collection{make(map[int]int), make(map[int]int)}
	return collection
}

func (sets Collection) FindSet(num int) int {
	// trace back the root of given num
	for true {
		if val, ok := sets.parent[num]; ok {
			num = val
		} else {
			break
		}
	}

	// reach the root and return it
	return num
}



func (sets Collection) UnionSet(x, y int) int {
	// find root of two numbers
	rootA := sets.FindSet(x)
	rootB := sets.FindSet(y)
	// within the same set
	if rootA == rootB {
		return rootA
	}

	// get the rank for two trees
	rankA := sets.rank[rootA]
	rankB := sets.rank[rootB]

	// append the shorter tree's root to the child of higher tree
	if rankA > rankB {
		sets.parent[rootB] = rootA
		return rootA
	} else if rankA < rankB {
		sets.parent[rootA] = rootB
		return rootB
	} else {
		sets.parent[rootA] = rootB
		sets.rank[rootB] = rankB + 1
		return rootB
	}

}