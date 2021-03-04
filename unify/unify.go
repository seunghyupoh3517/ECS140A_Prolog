package unify

import (
	"errors"
 	"hw4/disjointset"
	"hw4/term"
)

// ErrUnifier is the error value returned by the Parser if the string is not a
// valid term.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrUnifier = errors.New("unifier error")

// UnifyResult is the result of unification. For example, for a variable term
// `s`, `UnifyResult[s]` is the term which `s` is unified with.
type UnifyResult map[*term.Term]*term.Term

// Unifier is the interface for the term unifier.
// Do not change the definition of this interface
type Unifier interface {
	Unify(*term.Term, *term.Term) (UnifyResult, error)
}

type GeneralUnifier struct{
	disjointsets map[*term.Term]disjointset.DisjointSet		// disjoint sets for each representitive
	size map[*term.Term]int									// size for each disjoint set
	schema map[*term.Term]*term.Term						// map from disjoint set to a Term
	vars map[*term.Term][]*term.Term						// representitive -> list of TermVariable
}

/******************** Global Variable Definition ******* ***************/

// global variable to store the return unify result
var unifyMap = UnifyResult{}

// map for node <-> int in order to apply disjointset
var mapToInt = map[*term.Term]int{}
var mapToTerm = map[int]*term.Term{}
var nodeCounter = 0

// visited flag for all node
var visited = map[*term.Term]bool{}

// acyclic flag for all nodes
var acyclic = map[*term.Term]bool{}

/****************************** End *************************************/


// NewUnifier creates a struct of a type that satisfies the Unifier interface.
func NewUnifier() Unifier {
	var unifyObj Unifier = GeneralUnifier{ make(map[*term.Term]disjointset.DisjointSet), 
										   make(map[*term.Term]int),
										   make(map[*term.Term]*term.Term), 
									 	   make(map[*term.Term][]*term.Term) }
	return unifyObj 
}

// implements the Parse method with UnionSet struct
func (unif GeneralUnifier) Unify(t1 *term.Term, t2 *term.Term) (UnifyResult, error) {
	// initialization

	// map *Term to int 
	mapToInt[t1] = nodeCounter
	mapToTerm[nodeCounter] = t1
	nodeCounter++
	mapToInt[t2] = nodeCounter
	mapToTerm[nodeCounter] = t2
	nodeCounter++

	// create the disjointset class for two terms
	set1 := disjointset.NewDisjointSet()
	set2 := disjointset.NewDisjointSet()

	unif.disjointsets[t1] = set1
	unif.disjointsets[t2] = set2

	// intialize the schema for two classes
	unif.schema[t1] = t1
	unif.schema[t2] = t2

	// initialize the size for two classes
	unif.size[t1] = 1
	unif.size[t2] = 1

	// initialize the vars for two terms
	if t1.Typ == term.TermVariable {
		unif.vars[t1] = append(unif.vars[t1], t1)
	} else if t2.Typ == term.TermVariable {
		unif.vars[t2] = append(unif.vars[t2], t2)
	}

	
	unif.UnifClousure(t1, t2)
	unif.FindSolution(t1)

	return unifyMap, nil
}

func (unif GeneralUnifier) UnifClousure(t1 *term.Term, t2 *term.Term) error {
	num1 := unif.disjointsets[t1].FindSet(mapToInt[t1])		// num1 is int
	num2 := unif.disjointsets[t2].FindSet(mapToInt[t2])		// num2 is int return from FindSet()
	s := mapToTerm[num1]
	t := mapToTerm[num2]

	if s != t {
		schema_s := unif.schema[s]
		schema_t := unif.schema[t]
		if schema_s.Typ == term.TermVariable || schema_t.Typ == term.TermVariable {
			// one of their schema is variable 
			unif.Union(s, t)
		} else {
			// both are non-variable term
			if schema_s.Typ == term.TermCompound && schema_t.Typ == term.TermCompound {
				// both are compound term, compare the functor
				if schema_s.Functor == schema_t.Functor && len(schema_s.Args) == len(schema_t.Args){
					unif.Union(s, t)
					// loop through the args in two compound terms
					for i := range schema_s.Args {
						// try matching the parameters of compounds
						unif.UnifClousure(schema_s.Args[i], schema_t.Args[i])
					}
				} else {
					// not match with functor
					return ErrUnifier
				}
			}
		}

	}
	return nil
}

func (unif GeneralUnifier) Union(t1 *term.Term, t2 *term.Term) {
	s := unif.schema[t1]
	t := unif.schema[t2]
	
	if _, ok := mapToInt[s]; !ok {
		// not in the map
		mapToInt[s] = nodeCounter
		mapToTerm[nodeCounter] = s
		nodeCounter++
	}

	if _, ok := mapToInt[t]; !ok {
		mapToInt[t] = nodeCounter
		mapToTerm[nodeCounter] = t
		nodeCounter++
	}

	// TODO: Double check with using size to do comparison
	// Because the Union func in the disjoint class to comparing the rank
	// instead of size
	if unif.size[s] >= unif.size[t] {
		unif.size[s] += unif.size[t]
		unif.vars[s] = append(unif.vars[s], unif.vars[t]...)	// append the vars(t) to vars(s)
		// re assign the representitive of s to t's
		if unif.schema[s].Typ == term.TermVariable {
			unif.schema[s] = unif.schema[t]
		}
		unif.disjointsets[s].UnionSet(mapToInt[s], mapToInt[t])
	} else {
		// insert s's disjoint set to t's
		unif.size[t] += unif.size[s]
		unif.vars[t] = append(unif.vars[t], unif.vars[s]...)
		// re assign the representitive of t's to s's
		if unif.schema[t].Typ == term.TermVariable {
			unif.schema[t] = unif.schema[s]
		}
		unif.disjointsets[t].UnionSet(mapToInt[s], mapToInt[t])
	}
}

func (unif GeneralUnifier) FindSolution(t *term.Term) error {
	num := unif.disjointsets[t].FindSet(mapToInt[t])
	s := mapToTerm[num]
	s = unif.schema[s]

	if val, ok := acyclic[s]; ok {
		// TODO: Double check what need to return here??
		if val == true {

		}
	}

	if val, ok := visited[s]; ok {
		if val == true {
			return ErrUnifier
		}
	}

	if s.Typ == term.TermCompound {
		visited[s] = true
		for i := range s.Args {
			unif.FindSolution(s.Args[i])
		}
		visited[s] = false
	}

	acyclic[s] = true


	num2 := unif.disjointsets[t].FindSet(mapToInt[s])
	s = mapToTerm[num2]

	varsList := unif.vars[s]
	if len(varsList) > 0 {

		for _, x := range varsList {
			if x != s {
				unifyMap[x] = s
			}
		}
	}
	return nil
}
