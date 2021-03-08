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

	// TODO: need to delete
	// ToString() string
}

type DisjointSetImpl struct{
	parent map[int]int 		// map from node -> parent
	size map[int]int 		// map from node -> rank
}

//func (sets DisjointSetImpl) ToString() string {
//	return fmt.Sprint(" ", sets.parent)
//}

func (ds *DisjointSetImpl) FindSet(s int) int {
	if p, ok := ds.parent[s]; !ok {
		ds.parent[s] = s
		ds.size[s] = 1
		return s 
	} else if s != p {
		r := ds.FindSet(p)
		ds.parent[s] = r
		return r
	} else {
		return s
	}
}	

func (ds *DisjointSetImpl) UnionSet(s, t int) int {
	s, t = ds.FindSet(s), ds.FindSet(t)
	if s == t {
		return s
	}
	sizeS, sizeT := ds.size[s], ds.size[t]
	if sizeS < sizeT {
		s, t = t, s
	}
	ds.parent[t] = s
	ds.size[s] = sizeS + sizeT
	return s
}

func NewDisjointSet() DisjointSet {
	return &DisjointSetImpl{make(map[int]int), make(map[int]int)}
}

func MergeToUnionSets(set1 DisjointSet, set2 DisjointSet) DisjointSet {
	set_1 := set1.(*DisjointSetImpl)
	set_2 := set2.(*DisjointSetImpl)
	for key, val := range set_2.parent {
		if _, ok := set_1.parent[key]; !ok {
            set_1.parent[key] = val
		}
	}
	return set_1
}

func GetParent(setDis DisjointSet) map[int]int {
	set := setDis.(*DisjointSetImpl)
	return set.parent
}

// TODO: Need to delete, only for debug mode
// func (sets DisjointSetImpl) ToString() string {
// 	return fmt.Sprint(" ", sets.parent)
// }



/*

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
		sets.parent[rootB] = rootA
		sets.rank[rootA] = rankA + 1
		return rootA
	}

}

*/