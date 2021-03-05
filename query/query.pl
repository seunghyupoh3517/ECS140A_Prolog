/* All novels published either during the year 1953 or during the year 1996*/
year_1953_1996_novels(Book) :-
    novel(Book, 1953);
    novel(Book, 1996).

/* List of all novels published during the period 1800 to 1900 (not inclusive)*/
period_1800_1900_novels(Book) :-
    novel(Book, Year),
    Year > 1800,
    Year < 1900.

/* Characters who are fans of LOTR */
lotr_fans(Fan) :-
    fan(Fan, Books),
    find_book(the_lord_of_the_rings, Books).

/* Find if X contained in Books */
find_book(X, [X | _]).
find_book(X, [_ | Books]) :- find_book(X, Books).

/* Authors of the novels that heckles is fan of. */
heckles_idols(Author) :-
    author(Author, Author_books),
    fan(heckles, Fan_books),
    common_books(Author_books, Fan_books).

common_books([X |_], Books) :- find_book(X, Books).
common_books([_| Y], Books) :- common_books(Y, Books).


/* Characters who are fans of any of Robert Heinlein's novels */
heinlein_fans(Fan) :-
    author(robert_heinlein, Author_books),
    fan(Fan, Fan_books),
    common_books(Author_books, Fan_books).

/* Novels common between either of Phoebe, Ross, and Monica */
mutual_novels(Book) :-
    fan(phoebe, Phoebe_books),
    fan(ross, Ross_books),
    fan(monica, Monica_books),

    (between_fans(Book, Phoebe_books, Ross_books);
    between_fans(Book, Phoebe_books, Monica_books);
    between_fans(Book, Ross_books, Monica_books)).

between_fans(Book, FanBooks1, FanBooks2) :-
    find_book(Book, FanBooks1),
    find_book(Book, FanBooks2).

