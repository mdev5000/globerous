// A package for filtering and matching file paths based on  flexible file matching patterns.
//
// Note the package does not support true glob patterns but instead supports an offset. The package also supports regex
// and a hybrid between this simplified glob and regex.
//
// Supported glob patterns are:
//
// - Prefix wildcard:  *.txt
//
// - Postfix wildcard:  text*
//
// - Recursive wildcards: **
//
// See the NewCompiler example for examples of the different compiling schemes.
//
package globerous
