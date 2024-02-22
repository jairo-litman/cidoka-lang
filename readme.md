# The Cidoka Programming Language

This repo is an implementation of the Monkey Programming Language. Built by following Thorsten Ball's two books on the topics of [interpreters](https://interpreterbook.com/) and [compilers](https://compilerbook.com/) and then expanded upon with new features.

The implementation features both a tree-walking interpreter as well as a bytecode compiler with virtual machine. 

The user is free to choose which backend to use. In either case the code will eventually be interpreted in Go.

## Getting Started

### Requirements

* Go >= 1.14

### Running the REPL

The REPL runs using [peterh/liner](https://github.com/peterh/liner), visit the repo for more information and documentation on the library and commands.

`go run main.go`

The command above will run the REPL using the compiler+virtual machine, to run it using the interpreter add the `-engine=eval` flag or copy the following command.

`go run main.go -engine=eval`

The REPL features a history file written to the `tmp` directory. This file is used to store the history of commands entered into the REPL. The file is read from and written to when the REPL is started and stopped respectively.

The REPL also has an auto-completion feature. When you start typing a command, pressing the `tab` key will cycle through suggestions. The suggestions are the built-in functions and the language's keywords.

**Keyboard Shortcuts**

keystroke | action
--- | ---
`Ctrl + C` | Kill the REPL
`Ctrl + D` | Exit the REPL
`Ctrl + L` | Clear the screen
`Up Arrow` | Previous match from history
`Down Arrow` | Next match from history
`Tab` | Next completion
`Shift + Tab` | Previous completion

See the [liner](https://github.com/peterh/liner) repo for a full list of keyboard shortcuts.

### Running pre-written code

Look in `example/` for examples of Cidoka code. To run the code, use the `-input` flag followed by the absolute path to the file. Use the `-engine` flag to specify the backend to use, the default is the compiler+virtual machine.

`go run main.go -input=<absolute path to file>`

The following command will run the code in `example/helloworld.cidoka` using the vm.

`go run main.go -input=./example/helloworld.cidoka`

### Running the Benchmark

There's currenlty two benchmarks, both calculate the fibonacci sequence up to the 35th number, however one does so recursively and the other iteratively. Both benchmarks can be run using the compiler+virtual machine or the interpreter. 

At the end of the benchmark, the time taken to run the code will be printed to the console.

To run the benchmark recursively using the vm, run the following command:

`go run benchmark/benchmark.go`

The command accepts two flags, `-engine` and `-recursive`. The `-engine` flag accepts two values, `eval` and `vm`, the `-recursive` flag accepts two values, `true` and `false`. The default values are `vm` and `true` respectively.

## Supported Types

**Booleans**

Booleans are backed by go's native bool type.

```
true
false

true == false    -> false
true != false    -> true
true && false    -> false
true || false    -> true
!true            -> false
!false           -> true
```

**Strings**

Strings are backed by go's native string type. Printing is supported via the built-in print() function. String concatenation is supported with the `+` operator. Strings in Cidoka take the form of characters delimited by a pair of double quotes.

Examples:

```
"Graciela"
"Daniela"
"Hugo"

print("Cidoka")
"Cidoka " + "Lang"  -> "Cidoka Lang"
```

**Integers**

Integers are backed by go's native int type. Cidoka supports basic arithmetic operations on integers.

Examples:

```
1
10000
9122873

1 / 2   -> 0
1 + 2   -> 3
1 - 2   -> -1
1 * 2   -> 2
1 % 2   -> 1
```

**Floats**

Floats are backed by go's native float64 type. Cidoka supports basic arithmetic operations on floats.

Examples:

```
1.0
123.456
3.14159

1.0 / 2.0   -> 0.5
1.0 + 2.0   -> 3.0
1.0 - 2.0   -> -1.0
1.0 * 2.0   -> 2.0
```


**Arrays**

Arrays are backed by go's native slice type. Arrays are not scoped to a particular type in Cidoka so you can mix and match to your hearts content. Cidoka arrays take the form:

`[<expression>, <expression>, ...];`

You can index into an Array with an index expression. Array index expressions take the form:

 `<array>[<expression>];`

Examples:

```
["Ralph", "Abigail", "Bret", "Alejandro"]
[1,2,3,4]
[true, 1, false, "hello"]

["Ralph", "Abigail", "Bret", "Alejandro"][0]    -> "Ralph"

let people = ["Ralph", "Abigail", "Bret", "Alejandro"]
people[3]   -> "Alejandro"
```


**HashMaps/Dicts/Hashes**

Cidoka's kv data type is the Hash and it is backed by a go map. Like Arrays, they are not typed. Hashes take the form:

`{<expression>:<expression, <expression>:<expression, ....};`

It's worth noting that the keys of a hash can be any type that is hashable. These include integers, strings, and booleans. The values can be any type.

You can index into a Hash with an index expression. Hash index expressions takes the form:

 `<hash>[<expression>];`

Examples:

```
{1:2, 3:4, 5:6}

{"one":1, "two":2, "three":3}["one"]    -> 1
{1:"one", 2:"two", 3:"three"}[2]        -> "two"


let animals = {"Rodrigo":"parrot", "William":"giraffe", "Matt":"octopus"}

animals["Rodrigo"]          -> "parrot"
animals["Rod" + "rigo"]     -> "parrot"

```

**Functions**

Functions are first class in Cidoka. Additionally, closures are supported. 

If you don't have an explicit return in your Cidoka function, it will implicitly return the last expression.

Functions in Cidoka take the form:

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
closure();  -> 99
```

## Statements and Expressions

### Statements 

Programs in Cidoka are a series of statements.

Statements don't produce values. There are 7 types of statements in Cidoka.

**Expression Statements**

Expression statements are used to wrap expressions. These values are not reused.

`<expression>;`

```
5 + 5 * 2 / 9
```

**Let Statements**

Let statements allow you to declare names in the environment and bind expressions to them. Let statements scope to where you define them. If you use a let statement in the global scope it will be available to all functions. If you use it within a function, it will be bound to the lexical scope of that function.

`let <name> = <expression>;`

```
let result = 5 + 5 * 2 / 9
let concat = "fizz" + "buzz"
```

Once a name has been declared it cannot be redeclared in the same scope.

```
let hello = "world"
let hello = "foo"       -> Error: name hello already declared

let hello = "world"
fn() {
    let hello = "foo"   -> This is fine
}
```

**Return Statements**

Return statements are used to return a value from a function. If you don't have an explicit return in your Cidoka function, it will implicitly return the last expression.

`return <expression>;`

```
fn() {
    return 5 + 5 * 2 / 9
}
```

**Block Statements**

Block statements are used to group statements together. They are used by if, functions and loops. 

`{ <statements> }`

Note that the following example is not valid Cidoka code as block statements cannot be used by themselves.

```
{
    let result = 5 + 5 * 2 / 9
    let concat = "fizz" + "buzz"
}
```

**Loop Statements**

Loop statements are used to create for and while loops. Loops have their own scope, so any variables declared within a loop are not available to outer scopes.

`for (<statement>; <expression>; <statement>) { <statements> }`

`while (<expression>) { <statements> }`

Note that in most Cidoka code delimiting semicolons are optional. However, when defining a loop, the semicolons are required.

```
for (let i = 0; i < 10; i = i + 1) {
    print(i)
}

let i = 0;
while (i < 10) {
    print(i)
    i += 1
}
```

Under the hood, 'while' loops are just syntactic sugar for 'for' loops. The 'while' loop above is equivalent to the 'for' loop below.

```
let i = 0;
for (;i < 10;) {
    print(i)
    i += 1
}
```

The following loops are also equivalent and valid Cidoka code.

```
for (;;) {
    print("infinite loop")
}

while () {
    print("infinite loop")
}
```

**Break Statements**

Break statements are used to break out of a loop.

`break;`

```
for (let i = 0; i < 10; i = i + 1) {
    if (i == 5) {
        break;
    }
    print(i)
}
```

Will print 0, 1, 2, 3, 4

**Continue Statements**

Continue statements are used to skip the rest of the loop and start the next iteration.

`continue;`

```
for (let i = 0; i < 10; i = i + 1) {
    if (i == 5) {
        continue;
    }
    print(i)
}
```

Will print 0, 1, 2, 3, 4, 6, 7, 8, 9

### Expressions

Expressions produce values. These values can be reused in other expressions and combined with the statements listed in the previous section in order to bind an expression to a variable, return an expression, etc.

Most things in Cidoka are expressions, such as all the types listed previously.

**Identifier Expressions**

Identifier expressions are used to refer to a name in the environment. They evaluate to the value bound to the name.

`<name>`

```
let result = 1 + 2
result  -> 3
```

**Assignment Expressions**

Assignment expressions are used to rebind a declared name to a new value. They evaluate to the value that was assigned.

`<name> = <expression>`

```
let result = 1 + 2
result = 5  -> 5
```

They can also be used to modify an array or hash.

```
let arr = [1,2,3,4]
arr[2] = 5  -> 5
arr         -> [1,2,5,4]

let hash = {"one":1, "two":2, "three":3}
hash["two"] = 5     -> 5
hash                -> {"one":1, "two":5, "three":3}
```

They can also be used to rebind a name in an outer scope.

```
let result = 1 + 2
while () {
    result = 5  -> 5
    break;
}
result  -> 5

let arr = [1,2,3,4]
for (let i = 0; i < len(arr); i += 1) {
    arr[i] = 5 -> 5
}
arr     -> [5,5,5,5]
```

Assignment expressions support the following operators:

* `a += b` equivalent to `a = a + b`
* `a -= b` equivalent to `a = a - b`
* `a *= b` equivalent to `a = a * b`
* `a /= b` equivalent to `a = a / b`
* `a %= b` equivalent to `a = a % b`

**Prefix Expressions**

Prefix expressions are used to apply a prefix operator to an expression. They evaluate to the result of the operation.

`<prefix operator><expression>`

```
-5      -> -5
!true   -> false
!5      -> false
```

**Infix Expressions**

Infix expressions are used to apply an infix operator to two expressions. They evaluate to the result of the operation.

`<expression> <infix operator> <expression>`

```
5 + 5   -> 10
5 * 2   -> 15
```

The following is made up of a series of infix expressions.

```
5 + 5 * 2 / 9   -> 6
```

**If Expressions**

Cidoka supports conditional logic / flow control. This takes the form of:

`if (<expression>) { <statements> } else { <statements> };`

A nice feature of Cidoka is it doesn't use if statements but rather if expressions. This allows you to assign to a variable based on conditional logic.

```
let comparison = 5 > 3;

let val = if (comparison) { "Greater than" } else { "less than or equal"};

val     -> "Greater than"
```

It's worth noting that the else block is optional.

```
let comparison = 5 > 3;

let val = if (comparison) { "Greater than" };

val     -> "Greater than"
```

It's also possible to chain if expressions.

```
if (foo) {
    print("foo")
} else if (bar) {
    print("bar")
} else if (fizz) {
    print("fizz")
} else {
    print("bazz")
}
```

Note that most expressions in Cidoka are truthy.

```
if (5) {
    print("truthy")
}
```

**Function Call Expressions**

Function call expressions are used to call a function. They evaluate to the result of the function.

`<expression>(<comma-delimited expressions>)`

```
let add = fn(a,b) {a + b};
add(5,5)    -> 10
```

**Index Expressions**

Index expressions are used to index into an array or hash. They evaluate to the value at the given index.

`<expression>[<expression>]`

```
let arr = [1,2,3,4];
arr[2]  -> 3

let hash = {1:2, 3:4, 5:6};
hash[3] -> 4
```

## Built-in Functions

Cidoka comes with a few built-in functions which are run in Go. These functions are:

* `len(<array | string>)`
    - returns the length of an array or string
* `print(<string>)`
    - prints the given string to the console
* `first(<array>)`
    - returns the first element of an array
* `last(<array>)`
    - returns the last element of an array
* `tail(<array>)`
    - returns all elements of an array except the first
* `push(<array>, <element>)`
    - adds an element to the end of an array

## Missing Features and Possible Improvements

* Cli tool for generating binary executables
* More built-in functions
* Typing system
* Better error handling
* Better REPL
* More documentation
* More examples
* More comments

## License

This project is licensed under the GNU AFFERO GENERAL PUBLIC LICENSE Version 3 - see the LICENSE file for details

## Acknowledgments

Inspiration, code snippets, etc.
* [Writing An Interpreter In Go - Thorsten Ball](https://interpreterbook.com/)
* [Writing A Compiler In Go - Thorsten Ball](https://compilerbook.com/)
* [EvanErcolano/monkey](https://github.com/EvanErcolano/monkey)
* [peterh/liner](https://github.com/peterh/liner)
