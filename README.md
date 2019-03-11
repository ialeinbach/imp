# imp

*imp* is a toy, recursive assembly language. It is a work in progress. Currently, it can only be lexed and parsed.

#### Usage

Run `make test` to lex and parse a test source file with maximum verbosity enabled.

The programming models consists of 16 global, general purpose, 64-bit registers referred to as `@0`, `@1`, ..., `@15`. Integer constants are prefixed with a `#` (e.g. `#2`, `#-5`, `#0032`). There are only two types of statements: *subroutine definitions* and *subroutine calls*. Each must be terminated with a new line. Additionally, there must be at least one new line after the open brace of a definition.

#### Subroutine Definition
```
:subroutine_name @reg_alias, #int_alias, ... {
    subroutine_body
}
```

#### Subroutine Call
```
subroutine_name @reg, #int, ...
```

Register and integer constant arguments in a subroutine call are bound to register and integer constant aliases by position in the subroutine definition for the execution of the subroutine body.  Integer constants cannot be bound to register aliases and vice versa. There are two special subroutine calls: `rec` and `ret`.

`ret` given no arguments returns to the caller and continues execution at the next subroutine call. If `ret` is called with two arguments, it gets a question mark (i.e. must be called as `ret?`) and the arguments are compared for equality. On failure, the `ret` is skipped.

`rec` given no arguments jumps to the beginning of the currently executing subroutine body. If `rec` is called with two arguments, it gets a question mark (i.e. must be called as `rec?`), and the arguments are compared for equality. On failure, the `rec` is skipped.

Subroutine names can consist of lowercase ASCII letters and ASCII underscores (i.e. [a-z\_]). Register and integer aliases can be single lowercase or uppercase ASCII letters (i.e. [A-Za-z]). The same register cannot be passed to a subroutine call in more than one position.

Files are executed line-by-line starting at the first. Subroutine definitions can be nested but are only callable within the top level of scope within which they're defined (not including the use of `rec`). Subroutine definitions are also hoisted in the sense that they are callable before execution encounters their definitions.
