# imp

*imp* is a toy, recursive assembly language. It is a work in progress. Currently, the compiler (imp) can generate psuedo-instructions that will eventually be translated to architecture-specific machine code. There is an interpreter (twerp) with basic debugging features.

#### Usage

Download (without installing) by running `go get -u -d github.com/ialeinbach/imp`.

To build the compiler, run `make imp`.
To build the interpreter, run `make twerp`.
To build both, run `make`.

There are two types of statements: procedure calls (calls) and procedure declarations (decls). Newlines must be placed at the end of a call, end of a decl, and after the open brace of a decl. Decls cannot be nested (yet...?).

The calls in a decl body can only reference the aliases in that decl's parameter list (i.e. no globals). Parameter lists can contain integer and/or register aliases. Register parameters must be passed register arguments, but integer parameters can be passed either integer or register arguments. Typechecking is performed on calls to enforce these rules.

The programming model will eventually be dynamic with respect to compilation flags and target architecture limitations. Currently (and arbitrarily), there are 8 registers and procedures can have at most 6 arguments.

Control flow is implemented in a recursive style. There are two special builtins `ret` and `rec`. When passed 0 arguments, `ret` simply returns from the procedure and `rec` recurses (i.e. jumps to the beginning of the procedure). When passed 2 arguments, only when the arguments are equal do they return or recurse.

#### Todo

* Decide how to and implement plug-and-play target architectures.
* Optimize reg X passed as arg X to produce no psuedo-instructions (see examples/test3.imp).
* Read unicode point by unicode point rather than byte by byte.
