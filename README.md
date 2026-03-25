# v0tov1_expressions

A Go package for converting Harness CI v0 pipeline expressions to v1 format using trie-based pattern matching.

## Installation

```bash
go get github.com/peroxidemonke7/v0tov1_expressions
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    ce "github.com/peroxidemonke7/v0tov1_expressions/convert_expressions"
)

func main() {
    // Simple conversion without context
    v0Expr := "<+pipeline.stages.build.spec.execution.steps.myStep.output.outputVariables.foo>"
    v1Expr := ce.ConvertExpressionWithTrie(v0Expr, nil, false)
    
    fmt.Println(v1Expr)
    // Output: <+pipeline.stages.build.steps.myStep.output.outputVariables.foo>
}
```

### With Conversion Context

For context-aware conversions (e.g., step-specific transformations):

```go
package main

import (
    "fmt"
    ce "github.com/peroxidemonke7/v0tov1_expressions/convert_expressions"
)

func main() {
    // Create conversion context
    ctx := &ce.ConversionContext{
        StepType:        "Run",
        StepID:          "myStep",
        StepTypeMap:     map[string]string{"myStep": "Run"},
        CurrentStepType: "Run",
        UseFQN:          false,
    }

    v0Expr := "<+step.spec.envVariables.MY_VAR>"
    v1Expr := ce.ConvertExpressionWithTrie(v0Expr, ctx, false)
    
    fmt.Println(v1Expr)
}
```

### Using FQN (Fully Qualified Names)

When you need to convert expressions with full pipeline paths:

```go
ctx := &ce.ConversionContext{
    UseFQN:            true,
    CurrentStepV1Path: "pipeline.stages.build.steps.restoreCache",
    StepV1PathMap: map[string]string{
        "restoreCache": "pipeline.stages.build.steps.restoreCache",
        "buildImage":   "pipeline.stages.build.steps.buildImage",
    },
}

v0Expr := "<+steps.restoreCache.spec.key>"
v1Expr := ce.ConvertExpressionWithTrie(v0Expr, ctx, false)
```

## API Reference

### Main Function

#### `ConvertExpressionWithTrie(expr string, context *ConversionContext, inner bool) string`

Converts v0 expressions to v1 format using trie-based pattern matching.

**Parameters:**
- `expr`: The v0 expression string to convert
- `context`: Optional conversion context for step-aware transformations (can be `nil`)
- `inner`: Set to `false` for top-level calls (used internally for nested expressions)

**Returns:** Converted v1 expression string

### ConversionContext

```go
type ConversionContext struct {
    StepType          string            // Current step type (e.g., "Run", "GCSUpload")
    StepID            string            // Current step ID
    StepTypeMap       map[string]string // Map of step ID to step type
    CurrentStepType   string            // Fallback step type for "step.spec.*" expressions
    UseFQN            bool              // Enable fully qualified name mode
    CurrentStepV1Path string            // v1 FQN base path for current step
    StepV1PathMap     map[string]string // Map of step ID to v1 FQN paths
}
```

### Utility Functions

#### `FindHarnessExprs(s string) [][2]int`

Finds all top-level `<+...>` expressions in a string, handling nested expressions.

**Returns:** List of `[start, end]` index pairs for each expression

#### `GetPipelineTrie() *PipelineTrie`

Returns the cached singleton trie instance used for pattern matching.

## Local Development

If the module isn't published to a remote repository yet, use a `replace` directive in your project's `go.mod`:

```go
module your-project

go 1.24

require github.com/peroxidemonke7/v0tov1_expressions v0.0.0

replace github.com/peroxidemonke7/v0tov1_expressions => /path/to/local/v0tov1_expressions
```

## Features

- **Trie-based matching**: Efficient pattern matching for expression conversion
- **Nested expression support**: Handles arbitrarily nested `<+...>` expressions
- **Context-aware conversion**: Step-type specific transformations
- **FQN support**: Fully qualified name resolution for pipeline paths
- **Singleton trie caching**: Performance optimization for repeated conversions

## License

See LICENSE file for details.
