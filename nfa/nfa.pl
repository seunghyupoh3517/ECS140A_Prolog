reachable(Nfa, StartState, FinalState, Input) :-
    %% TODO: remove fail and add body/other cases for this predicate
    helper(Nfa, StartState, FinalState, Input).

%% Basic condition: when no more input symbol - whether at final state or not
helper(Nfa, StartState, FinalState, []) :-
    StartState = FinalState;
    member(FinalState, StartState).

%% Single next state 
helper(Nfa, StartState, FinalState, Input) :-
    Input = [IH|IT],
    transition(Nfa, StartState, IH, NextState), 
    helper(Nfa, NextState, FinalState, IT).

%% Multiple next states - Check the head and pass the tail 
helper(Nfa, StartState, FinalState, Input) :-
    StartState = [SH|ST],
    Input = [IH|IT],
    transition(Nfa, SH, IH, NextState),
    helper(Nfa, NextState, FinalState, IT).

helper(Nfa, StartState, FinalState, Input) :-
    StartState = [SH|ST],
    helper(Nfa, ST, FinalState, Input).


