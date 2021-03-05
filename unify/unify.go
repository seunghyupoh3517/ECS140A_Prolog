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

/****************************** End *************************************/


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

	if _, ok := mapToInt[t2]; !ok {
		mapToInt[t2] = nodeCounter
		mapToTerm[nodeCounter] = t2
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

	if _, ok := unif.disjointsets[t2]; !ok {
		set2 := disjointset.NewDisjointSet()
		unif.disjointsets[t2] = set2
		unif.size[t2] = 1
		unif.schema[t2] = t2
		if t2.Typ == term.TermVariable {
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
		return nil, ErrUnifier
	}
	// fmt.Println(" *** Debug info:  from line 107")
	err_findSol := unif.FindSolution(t1)
	if err_findSol != nil {
		return nil, ErrUnifier
	}
	// fmt.Println(" *** Debug info:  from line 112")
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
			// fmt.Println(" *** Debug info: from line 127")
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
		// fmt.Println("The unif.vars[t] in line 193", unif.vars[t])
		if unif.schema[s].Typ == term.TermVariable {
			unif.schema[s] = unif.schema[t]
		}
		// fmt.Println(unif.disjointsets[s].UnionSet(mapToInt[s], mapToInt[t]))
	} else {
		// insert s's disjoint set to t's
		unif.size[t] += unif.size[s]
		unif.vars[t] = append(unif.vars[t], unif.vars[s]...)
		fmt.Println("The unif.vars[t] in line 201", unif.vars[t])
		// re assign the representitive of t's to s's
		if unif.schema[t].Typ == term.TermVariable {
			unif.schema[t] = unif.schema[s]
		}
		unif.disjointsets[t].UnionSet(mapToInt[s], mapToInt[t])
	}
}


// 这个asyclic很迷,我感觉一开始并不能等于true(默认是true好像),因为这样的话大部分Unify都会报错
// if s.Typ == term.TermCompound 在这个function里面calling FindSolution出了问题
// 只要最后一项不是多项compound term就可以pass， 比如最后是f(1),f(x)就pass不了
// 按理来说最后的compoundterm应该是empty set， 但是咱们的program判定如果是compound term就会还有后续，继续往里找就出错了
// 单项中比如{"f(X)", "f(B)"}不会进行recursive，所以就不会出错
// 出错在num := unif.disjointsets[t].FindSet(mapToInt[t])这一行

// I let the line 241 val == false, original is val == true
// if keep the ErrorSymbolClash open, the "TestUnifySuccess" will fail somehow????? Global ?


func (unif GeneralUnifier) FindSolution(t *term.Term) error {
	// fmt.Println(" *** Debug info: ", t, "line 210")
	// fmt.Println(" *** Debug info: unif.disjointSet =", unif.disjointsets[t], "line 195")	
	// fmt.Println(" *** Debug info: map =", unif.disjointsets)
	// fmt.Println("Beginning of the FindSol	")
	// fmt.Println(mapToInt[t])
	if mapToInt[t] == 0{
		return nil
	}
	num := unif.disjointsets[t].FindSet(mapToInt[t])
	s := mapToTerm[num]
	s = unif.schema[s]
	
	// acyclic[s] = false
	// fmt.Println(" *** Debug info: ", s, "line 219")
	// fmt.Println(" *** Debug info: schema[] =", unif.schema, "line 218")
	if val, ok := acyclic[s]; ok {
		// TODO: Double check what need to return here??
		// fmt.Println(" *** Debug info: from line 221 val == ", val)
		if val == false{
			// fmt.Println(" *** Debug info: from line 225")
			return ErrUnifier
		}
		
	}
	// fmt.Println(" *** Debug info line 236")
	 if val, ok := visited[s]; ok {
		if val == true {
			// fmt.Println(" *** Debug info: from line 229")
			return ErrUnifier		// exits a cycle
		}
	}
	// fmt.Println(" *** Debug info line 243, s=", s)
	if s.Typ == term.TermCompound {
		// fmt.Println(" *** Debug info: from line 216")
		visited[s] = true
		for i := range s.Args {
			// fmt.Println(" *** Debug info: i=", s.Args[i], "from line 248")
			err := unif.FindSolution(s.Args[i])
			// fmt.Println("11111111111111")
			if err != nil {
				// fmt.Println(" *** Debug info: from line 222")
				return ErrUnifier
			}		
		}
		visited[s] = false
	}

	acyclic[s] = true

	// fmt.Println(" *** Debug info line 260")
	num2 := unif.disjointsets[t].FindSet(mapToInt[s])
	s = mapToTerm[num2]
	// fmt.Println(" *** Debug info line 263")
	// fmt.Println(" *** Debug info: from line 253")
	// fmt.Println(" *** Debug info: s =",s ,"- from line 254")
	varsList := unif.vars[s]
	// fmt.Println(" *** Debug info: varslist =",varsList ,"- from line 268")

	if len(varsList) > 0 {
		for _, x := range varsList {
			// fmt.Println("x= ", x,"s= ",s)
			if x != s {
				// fmt.Println(" *** Debug info: from line 259")
				unifyMap[x] = s
				// fmt.Println("The UnifyMap is ", unifyMap[x])
			}
		}
	}

	return nil
}
