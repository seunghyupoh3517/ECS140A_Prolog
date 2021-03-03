% A list is a 1-D array of numbers.
% A matrix is a 2-D array of numbers, stored in row-major order.

% You may define helper functions here.

% are_adjacent(List, A, B) returns true iff A and B are neighbors in List.
are_adjacent(List, A, B) :-
    (match(A, List), matchNext(B, List),!);
    (match(B, List), matchNext(A, List)),!.
are_adjacent([_ | X], A, B) :- are_adjacent(X, A, B).

match(X, [X | _]).
matchNext(X, [_ | Y]) :- match(X, Y).

% matrix_transpose(Matrix, Answer) returns true iff Answer is the transpose of
% the 2D matrix Matrix.
matrix_transpose(Matrix, Answer) :-
    make_transpose(Matrix, Tmatrix),
    Answer = Tmatrix.

% learned the implementation to make a transpose matrix at
% https://stackoverflow.com/questions/4280986/
% how-to-transpose-a-matrix-in-prolog
make_transpose([], []).
make_transpose([F|Fs], Ts) :- make_transpose(F, [F|Fs], Ts).

make_transpose([], _, []).
make_transpose([_|Rs], Ms, [Ts|Tss]) :-
        lists_firsts_rests(Ms, Ts, Ms1),
        make_transpose(Rs, Ms1, Tss).

lists_firsts_rests([], [], []).
lists_firsts_rests([[F|Os]|Rest], [F|Fs], [Os|Oss]) :-
        lists_firsts_rests(Rest, Fs, Oss).

% are_neighbors(Matrix, A, B) returns true iff A and B are neighbors in the 2D
% matrix Matrix.
are_neighbors(Matrix, A, B) :-
    check_row(Matrix, A, B),!;
    (
        make_transpose(Matrix, Tmatrix),
        check_row(Tmatrix, A, B)
    ).

check_row([Row | _], A, B) :- are_adjacent(Row, A, B),!.
check_row([_ | NextRow], A, B) :- check_row(NextRow, A, B),!.
