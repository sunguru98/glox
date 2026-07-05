# glox

A tree-walk interpreter for [Lox](https://craftinginterpreters.com/the-lox-language.html) — the teaching language from Bob Nystrom's book *Crafting Interpreters* — written in Go.

Lox is a small, dynamically-typed scripting language with C-like syntax, closures, and classes. `glox` reads Lox source code and executes it directly (no compilation to bytecode or machine code — it walks the syntax tree and evaluates it as it goes).

## Features

- Arithmetic, comparison, and logical expressions (`+ - * / == != < > and or !`)
- Variables and assignment (`var`)
- Control flow: `if` / `else`, `while`, `for`
- Functions with closures (`fun`, `return`)
- Classes with inheritance, `this`, and `super`
- A native `clock()` function
- A REPL (interactive prompt) and file execution mode

## Requirements

- Go 1.26+

## Installation

```bash
git clone https://github.com/sunguru98/glox.git
cd glox
go build -o glox .
```

## Usage

Run a `.lox` script:

```bash
./glox lox-programs/functions.lox
```

Or start the interactive REPL (leave the input blank to exit):

```bash
./glox
> print "Hello, world!";
Hello, world!
> var x = 40 + 2;
> print x;
42
>
```

Example scripts to try live in [`lox-programs/`](./lox-programs), including recursive functions, closures, and scoping demos.

## Language quick reference

```js
// Variables
var name = "Lox";
var count = 0;

// Functions and closures
fun makeCounter() {
  var i = 0;
  fun count() {
    i = i + 1;
    print i;
  }
  return count;
}

var counter = makeCounter();
counter(); // 1
counter(); // 2

// Control flow
for (var i = 0; i < 5; i = i + 1) {
  if (i == 3) print "three!";
  else print i;
}

// Classes
class Animal {
  init(name) {
    this.name = name;
  }
  speak() {
    print this.name + " makes a sound.";
  }
}

class Dog < Animal {
  speak() {
    print this.name + " barks.";
  }
}

var d = Dog("Rex");
d.speak(); // Rex barks.
```

## How it works

Source code passes through four stages, wired together in `main.go`:

```
source text
    │
    ▼
Scanner   (scanner/)     → splits raw text into tokens (numbers, identifiers, operators, keywords)
    │
    ▼
Parser    (parser/)      → turns tokens into an AST based on Lox's grammar
    │
    ▼
Resolver  (resolver/)    → walks the AST once ahead of time to resolve variable scope,
    │                       so closures capture the correct variable at the correct depth
    ▼
Interpreter (parser/interpret.go) → walks the AST and evaluates it
```

### Package layout

| Path | Responsibility |
|---|---|
| `scanner/` | Lexer: converts source text into a stream of `Token`s. |
| `parser/parser.go` | Recursive-descent parser: builds the AST from tokens. |
| `parser/expression.go` | AST node types for expressions (`Binary`, `Unary`, `Literal`, `Call`, `Get`, `Set`, `This`, `Super`, ...). |
| `parser/statement.go` | AST node types for statements (`If`, `While`, `Function`, `Class`, `Return`, ...). |
| `parser/environment.go` | Variable scope storage — a chain of maps, one per block/function. |
| `parser/callable.go` | Runtime representations of functions, classes, and instances. |
| `parser/interpret.go` | The tree-walking evaluator that executes the AST. |
| `resolver/` | Static analysis pass for resolving variable scope before execution. |
| `lib/` | Shared error-reporting state (`HadError`, `HadRuntimeError`) used across all stages. |
| `lox-programs/` | Example `.lox` scripts. |

## Exit codes

- `0` — success
- `64` — usage error (more than one script path passed)
- `65` — a parse/static error was reported
- `70` — a runtime error occurred

## Reference

- [Crafting Interpreters](https://craftinginterpreters.com/) — the book this implementation follows (the `jlox` reference implementation is in Java; this is a Go port with the same grammar and semantics).