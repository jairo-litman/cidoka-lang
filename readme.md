# The Boludo Programming Language

This repo is an implementation of the Monkey Programming Language. Built by following Thorsten Ball's two books on the topics of [interpreters](https://interpreterbook.com/) and [compilers](https://compilerbook.com/).

The implementation features both a tree-walking interpreter as well as a bytecode compiler with virtual machine. 

The user is free to choose which backend to use. In either case the language will eventually be executed in Go.

## Getting Started

### Requirements

* Go >= 1.14

### Running the REPL

`go run main.go`

The command above will run the REPL using the compiler+virtual machine, to run it using the interpreter add the `-engine=eval` flag or copy the following command.

`go run main.go -engine=eval`

Look in `benchmark/benchmark.go` for an example of how to run pre written code.

## Supported Types

**Booleans**

Booleans are backed by go's native bool type.


```
true
false
```

**Strings**

Strings are backed by go's native string type. Printing is supported via the built-in print() function. String concatenation is supported with the `+` operator. Strings in BoludoLang take the form of characters delimited by a pair of double quotes.

Examples:


```
"Graciela"
"Daniela"
"Hugo"

print("Boludo")
"Boludo " + "Lang"
```

**Integers**

Integers are backed by go's native int type. BoludoLang supports basic arithmetic operations on integers.

Examples:


```
1
10000
9122873

1 / 2
1 + 2
1 - 2
1 * 2
```

**Arrays**

Arrays are backed by go's native slice type. Arrays are not scoped to a particular type in BoludoLang so you can mix and match to your hearts content. BoludoLang arrays take the form:

`[<expression>, <expression>, ...];`

You can index into an Array with an index expression. Array index expressions take the form:

 `<array>[<expression>];`

Examples:


```
["Ralph", "Abigail", "Bret", "Alejandro"]
[1,2,3,4]
[true, 1, false, "hello"]

["Ralph", "Abigail", "Bret", "Alejandro"][0] -> "Ralph"

let people = ["Ralph", "Abigail", "Bret", "Alejandro"]
people[3] -> "Alejandro"
```


**HashMaps/Dicts/Hashes**

BoludoLang's kv data type is the Hash and it is backed by a go map. Like Arrays, they are not typed. Hashes take the form:

`{<expression>:<expression, <expression>:<expression, ....};`

You can index into a Hash with an index expression. Hash index expressions takes the form:

 `<hash>[<expression>];`

Examples:

```
{1:2, 3:4, 5:6}"
let animals = {"Rodrigo":"parrot", "William":"giraffe", "Matt":"octopus"}

animals["Rodrigo"] -> "parrot"
animals["Rod" + "rigo"] -> "parrot"

```

**Functions**

Functions are first class in BoludoLang. Additionally, closures are supported. If you don't have an explicit return in your BoludoLang function, it will implicitly return the last expression.

Functions in BoludoLang take the form:

```
fn(<optional comma-delimited identifiers>) {
    <optional statements>
    <optional return statement>
}
```

Example self-referential recursive function:

```
let fibonacci = fn(x) {
  if (x == 0) {
    return 0;
  } else {
    if (x == 1) {
      return 1;
    } else {
      return fibonacci(x - 1) + fibonacci(x - 2);
    }
  }
};

fibonacci(10);
```

Example Closure

```
let newClosure = fn(a,b) {
    let one = fn() {a;};
    let two = fn() {b;};
    return fn() {one() + two();};
};

let closure = newClosure(9,90);
closure(); -> 99
```

## Statements and Expressions

**Statements**

Programs in BoludoLang are a series of statements.

Statements don't produce values. There are three types of statements in BoludoLang.

1. let statements
    - Bind expressions to an identifier
2. return statements
    - return the value produced by an expression from a function
3. expression statements
    - wrap expressions, these values are not reused


**Expressions**

Expressions produce values. These values can be reused in other expressions and combined with the statements listed in the previous section in order to bind an expression to a variable or return an expression...

BoludoLang supports both infix and prefix expressions.

**Let Statements**

Let statements allow you to bind expressions to names in the environment. Let statements scope to where you define them. If you use a let statement in the global scope it will be available to all functions. If you use it within a function, it will be bound to the lexical scope of that function.

`let <name> = <expression>;`

```
let result = 5 + 5 * 2 / 9
let concat = "fizz" + "buzz"
```


**If Expressions**

BoludoLang supports conditional logic / flow control. This takes the form of:

`if (<expression>) { <statements> } else { <statements> };`

A nice feature of BoludoLang is it doesn't use if statements but rather if expressions. This allows you to assign to a variable based on conditional logic.

```
let comparison = 5 > 3;

let val = if (comparison) { "Greater than" } else { "less than or equal"};

val -> "Greater than"
```

## License

This project is licensed under the GNU GENERAL PUBLIC LICENSE V3 License - see the LICENSE file for details

## Acknowledgments

Inspiration, code snippets, etc.
* [Writing An Interpreter In Go - Thorsten Ball](https://interpreterbook.com/)
* [Writing A Compiler In Go - Thorsten Ball](https://compilerbook.com/)
* [EvanErcolano/monkey](https://github.com/EvanErcolano/monkey)