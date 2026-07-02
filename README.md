# GNU Flag

A lightweight, GNU-style command-line flag parser for Go.

## Overview

GNU Flag is a simple yet powerful command-line argument parser that follows GNU-style conventions. It provides an easy way to define and parse command-line flags with support for short and long options, help generation, and error handling.

## Features

- **GNU-style flags**: Support for both short (`-f`) and long (`--flag`) options
- **Automatic help generation**: Built-in `--help` and `-h` flags with formatted output
- **Multiple flag types**: String, integer, and boolean flags
- **Error handling**: Clear error messages for invalid options and missing arguments
- **Positional arguments**: Support for non-flag arguments
- **Simple API**: Easy to use with minimal boilerplate

## Installation

```bash
go get github.com/robogg133/gnuflag
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "os"
    "github.com/robogg133/gnuflag"
)

func main() {
    // Create a new parser
    parser := gnuflag.NewParser(os.Args)
    
    // Set your program's usage information
    parser.Usage = "[options] [files...]"
    parser.Slogan = "A simple program that does something useful"
    parser.Description = "This program demonstrates the GNU Flag parser with various options and flags."
    
    // Define flags
    var name string
    parser.SetFlagString("name", "your name", &name, "n")
    
    var age int
    parser.SetFlagInt("age", "your age", &age, "a")
    
    var verbose bool
    parser.SetFlagBool("verbose", "enable verbose output", &verbose, "v")
    
    // Parse the arguments
    if err := parser.Parse(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    
    // Access flags
    if name != "" {
        fmt.Printf("Hello, %s!\n", name)
    }
    
    if age > 0 {
        fmt.Printf("You are %d years old.\n", age)
    }
    
    if verbose {
        fmt.Println("Verbose mode enabled")
    }
    
    // Access positional arguments
    for i := 0; i < parser.NArgs(); i++ {
        fmt.Printf("Argument %d: %s\n", i, parser.Arg(i))
    }
}
```

### Example Usage

```bash
# With flags
$ ./program --name "John" --age 30 --verbose file1.txt file2.txt
Hello, John!
You are 30 years old.
Verbose mode enabled
Argument 0: file1.txt
Argument 1: file2.txt

# With short flags
$ ./program -n John -a 30 -v

# Get help
$ ./program --help
Usage: program [options] [files...]
A simple program that does something useful

  -n, --name      your name
  -a, --age       your age
  -v, --verbose   enable verbose output
  -h, --help      display this help and exit

This program demonstrates the GNU Flag parser with various options and flags.
```

## API Reference

### Parser

#### `func NewParser(args []string) *Parser`

Creates a new parser instance. Typically called with `os.Args`.

#### `func (p *Parser) SetFlagString(full, help string, value *string, shorts ...string)`

Defines a string flag.

```go
var config string
parser.SetFlagString("config", "configuration file path", &config, "c")
```

#### `func (p *Parser) SetFlagInt(full, help string, value *int, shorts ...string)`

Defines an integer flag.

```go
var port int
parser.SetFlagInt("port", "port number", &port, "p")
```

#### `func (p *Parser) SetFlagBool(full, help string, value *bool, shorts ...string)`

Defines a boolean flag.

```go
var debug bool
parser.SetFlagBool("debug", "enable debug mode", &debug, "d")
```

#### `func (p *Parser) Parse() error`

Parses the command-line arguments. Returns an error if parsing fails.

#### `func (p *Parser) NArgs() int`

Returns the number of positional arguments.

#### `func (p *Parser) Arg(n int) string`

Returns the positional argument at index `n`.

### Properties

| Property | Description |
|----------|-------------|
| `CommandName` | Name of the command/program |
| `Usage` | Usage string (shown in help) |
| `Slogan` | Short description (shown in help) |
| `Description` | Long description (shown in help) |

## Error Handling

The parser provides specific error types for common cases:

- `ErrorInvalidOption`: Invalid short option
- `ErrorNotRecognized`: Unrecognized long option
- `ErrorRequiresAnArg`: Option requires an argument
- `ErrorRequired`: Missing required option

All errors include a helpful message with a suggestion to use `--help`.

## Help System

The help system is automatically generated and includes:

- Usage information
- Program slogan
- List of all defined flags with descriptions
- Program description

Help is triggered by either `--help` or `-h`.
