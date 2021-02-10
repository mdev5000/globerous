# File pattern matching for go lang

A package for running flexible file matching patterns to filter and match files.

{{ embed "examples_simple_test.go" }}

## Supported globbing patterns

A limited set of globbing patterns are supported specifically:

- prefix wildcard: `*.txt`
- suffix wildcard: `file.*`
- recursive wildcard: `**`

For more control of patterns use the **regex** or **hybrid** match compiler. See the
next section for more details.


## Running regex patterns

You can run regex patterns with either the `globerous.RegexPartCompiler` or via the
`globerous.HybridGlobRegexPartCompiler`. The latter allows use of both regex and glob
pattern matching.

Here's an example of the hybrid compiler:

{{ embed "examples_hybdridmatch_test.go" }}

## Virtual filesystems

Globerous has support for the **blang vfs** library (https://github.com/blang/vfs). This
allows you to virtualize the filesystem and for example use an in-memory fs.

{{ embed "examples_vfs_test.go" }}

## Custom pattern matching

You can extend or replace the pattern matching used by using your own part compiler.

{{ embed "examples/custompartcompiler/main.go" }}

You can also use the raw matchers instead of compiling.

{{ embed "examples/manualmatching/main.go" }}

# Testing paths

You can test paths either using real filesystem information or simple path information.

{{ embed "examples/testingpaths/main.go" }}
