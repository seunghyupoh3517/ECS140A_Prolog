/* fooTransitions */
% 0 -a-> 1
% 0 -a-> 2
% 1 -b-> 3
% 2 -c-> 3
transition(fooTransitions, 0, a, [1,2]).
transition(fooTransitions, 0, b, []).
transition(fooTransitions, 0, c, []).
transition(fooTransitions, 1, a, []).
transition(fooTransitions, 1, b, [3]).
transition(fooTransitions, 1, c, []).
transition(fooTransitions, 2, a, []).
transition(fooTransitions, 2, b, []).
transition(fooTransitions, 2, c, [3]).

/* expTransitions */
% 0 -a-> 1
% 0 -a-> 2
% 0 -b-> 2
% 1 -b-> 0
transition(expTransitions, 0, a, [1,2]).
transition(expTransitions, 0, b, [2]).
transition(expTransitions, 1, a, []).
transition(expTransitions, 1, b, [0]).
transition(expTransitions, 2, a, []).
transition(expTransitions, 2, b, []).

/* langTransitions */
% 0 -a-> 0
% 0 -b-> 1
% 1 -a-> 1
% 1 -b-> 0
transition(langTransitions, 0, a, [0]).
transition(langTransitions, 0, b, [1]).
transition(langTransitions, 1, a, [1]).
transition(langTransitions, 1, b, [0]).


/* manualTransitions */
% 0 -a-> 1
% 0 -a-> 2
% 0 -a-> 3
% 0 -a-> 4
% 1 -a-> 0
% 2 -b-> 0
% 2 -b-> 3
% 3 -c-> 0
% 3 -c-> 4
% 4 -d-> 0
% 4 -d-> 1
transition(manualTransitions, 0, a, [1, 2, 3, 4]).
transition(manualTransitions, 1, a, [0]).
transition(manualTransitions, 2, b, [0, 3]).
transition(manualTransitions, 3, c, [0, 4]).
transition(manualTransitions, 4, d, [0, 1]).