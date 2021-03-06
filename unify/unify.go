package unify

import (
	"errors"
 	"hw4/disjointset"
	"hw4/term"
	"fmt"
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


// TODO: need to reset the global value before return the result
/****************************** End *************************************/


// reset the global variable
func resetGlobal() {
	// unifyMap = UnifyResult{}
	mapToInt = map[*term.Term]int{}
	mapToTerm = map[int]*term.Term{}
    nodeCounter = 0
	visited = map[*term.Term]bool{}
	acyclic = map[*term.Term]bool{}
}

// NewUnifier creates a struct of a type that satisfies the Unifier interface.
func NewUnifier() Unifier {
	var unifyObj Unifier = GeneralUnifier{ make(map[*term.Term]disjointset.DisjointSet), 
										   make(map[*term.Term]int),
										   make(map[*term.Term]*term.Term), 
									 	   make(map[*term.Term][]*term.Term) }
	return unifyObj 
}

func (unif GeneralUnifier) Initializer(t1 *term.Term, t2 *term.Term) {
	// maps  *Term <--->  int 
	if _, ok := mapToInt[t1]; !ok {
		mapToInt[t1] = nodeCounter
		mapToTerm[nodeCounter] = t1
		nodeCounter++
	}

	// create the disjointset class for two terms
	if _, ok := unif.disjointsets[t1]; !ok {
		set1 := disjointset.NewDisjointSet()
		unif.disjointsets[t1] = set1
		unif.size[t1] = 1		// initialize the size for the class
		unif.schema[t1] = t1	// intialize the schema for the class
		// initialize the vars for the term
		if t1.Typ == term.TermVariable {
			unif.vars[t1] = append(unif.vars[t1], t1)
		}
	}

	if t2 == nil {
		return
	}

	if _, ok := mapToInt[t2]; !ok {
		mapToInt[t2] = nodeCounter
		mapToTerm[nodeCounter] = t2
		nodeCounter++
	}

	if _, ok := unif.disjointsets[t2]; !ok {
		set2 := disjointset.NewDisjointSet()
		unif.disjointsets[t2] = set2
		unif.size[t2] = 1
		unif.schema[t2] = t2

		if t2 != nil && t2.Typ == term.TermVariable {
			unif.vars[t2] = append(unif.vars[t2], t2)
		}
	}
}

// implements the Unify method with UnionSet struct
func (unif GeneralUnifier) Unify(t1 *term.Term, t2 *term.Term) (UnifyResult, error) {
	// initialization
	unif.Initializer(t1, t2)

	err_clousure := unif.UnifClousure(t1, t2)
	if err_clousure != nil {
		resetGlobal()
		return nil, ErrUnifier
	}
	// fmt.Println(" *** Debug info:  from line 107")
	err_findSol := unif.FindSolution(t1)
	if err_findSol != nil {
		resetGlobal()
		return nil, ErrUnifier
	}
	fmt.Println(" *** Debug info:  from line 112")
	resetGlobal()
	return unifyMap, nil
}

func (unif GeneralUnifier) UnifClousure(t1 *term.Term, t2 *term.Term) error {
	num1 := unif.disjointsets[t1].FindSet(mapToInt[t1])		// num1 is int
	num2 := unif.disjointsets[t2].FindSet(mapToInt[t2])		// num2 is int return from FindSet()
	s := mapToTerm[num1]
	t := mapToTerm[num2]
	// fmt.Println(" *********** Debug info:  from line 121 **********")
	if s != t {
		schema_s := unif.schema[s]
		schema_t := unif.schema[t]
		if schema_s.Typ == term.TermVariable || schema_t.Typ == term.TermVariable {
			// one of their schema is variable 
			fmt.Println(" *** Debug info: s = ",s , ", t =", t ,"- line 127")
			fmt.Println(" *** Debug info: schema_s = ",schema_s , ", schema_t =",schema_t ,"- line 128")
			unif.Union(s, t)
		} else {
			// fmt.Println(" *** Debug info: from line 123")
			// both are non-variable term
			if schema_s.Typ == term.TermCompound && schema_t.Typ == term.TermCompound {
				// both are compound term, compare the functor
				if schema_s.Functor == schema_t.Functor && len(schema_s.Args) == len(schema_t.Args){
					// fmt.Println(" *** Debug info:  from line 128")
					unif.Union(s, t)
					// fmt.Println(" *** Debug info:  disjointsets =", unif.disjointsets)
					// loop through the args in two compound terms
					for i := range schema_s.Args {
						// fmt.Println(" *** Debug info: s =",schema_s.Args[i] , " from line 133")
						// fmt.Println(" *** Debug info: t =",schema_t.Args[i] , " from line 134")						
						
						// create the disjointSet for parameters
						unif.Initializer(schema_s.Args[i], schema_t.Args[i])
						
						// try matching the parameters of compounds
						err := unif.UnifClousure(schema_s.Args[i], schema_t.Args[i])
						// fmt.Println(" *** Debug info: from line 136")
						if err != nil {
							return ErrUnifier
						}
					}
				} else {
					// not match with functor
					return ErrUnifier
				}
			} else {
				// unmatched term type or single term
				// TODO: Double check if   f, f --> return true
				// fmt.Println(" *** Debug info: from line 144")
				return ErrUnifier
			}
		}

	}
	return nil
}

func (unif GeneralUnifier) Union(t1 *term.Term, t2 *term.Term) {
	// TODO: Double check if this one need to call schema
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
		fmt.Println(" *** Debug info: check s=",s ," - line 193")
		fmt.Println(" *** Debug info: check vars(s)=",unif.vars[s] ," - line 194")
		fmt.Println(" ***************************************************")
		
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
	// fmt.Println(" *** Debug info: ", t, "line 216")
	// fmt.Println(" *** Debug info: unif.disjointSet =", unif.disjointsets[t], "line 195")	
	// fmt.Println(" *** Debug info: map =", unif.disjointsets)
	num := unif.disjointsets[t].FindSet(mapToInt[t])

	s := mapToTerm[num]
	// fmt.Println(" *** Debug info: ", s, "line 222")
	s = unif.schema[s]
	// fmt.Println(" *** Debug info: ", s, "line 224")


	// fmt.Println(" *** Debug info: ", s, "line 217")
	// fmt.Println(" *** Debug info: schema[] =", unif.schema, "line 218")
	if val, ok := acyclic[s]; ok {
		// TODO: Double check what need to return here??
		fmt.Println(" *** Debug info: from line 233")
		if val == true {
			// return ErrUnifier
		}
	}

	if val, ok := visited[s]; ok {
		if val == true {
			// fmt.Println(" *** Debug info: from line 210")
			return ErrUnifier		// exits a cycle
		}
	}

	if s.Typ == term.TermCompound {
		// fmt.Println(" *** Debug info: from line 216")
		visited[s] = true
		for i := range s.Args {
			// fmt.Println(" *** Debug info: i=", s.Args[i], "from line 248")
			// fmt.Println("********************************")
			unif.Initializer(s.Args[i], nil)
			err := unif.FindSolution(s.Args[i])
			if err != nil {
				// fmt.Println(" *** Debug info: from line 252")
				return ErrUnifier
			}		
		}
		visited[s] = false
	}

	acyclic[s] = true


	num2 := unif.disjointsets[t].FindSet(mapToInt[s])
	s = mapToTerm[num2]
	// fmt.Println(" *** Debug info: from line 253")
	// fmt.Println(" *** Debug info: s =",s ,"- from line 254")
	varsList := unif.vars[s]
	// fmt.Println(" *** Debug info: varslist =",varsList ,"- from line 256")
	if len(varsList) > 0 {
		for _, x := range varsList {
			if x != s {
				// fmt.Println(" *** Debug info: from line 259")
				unifyMap[x] = s
			}
		}
	}
	return nil
}
