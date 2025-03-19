# Let's Byte

Annie and Kevin's compilers project for MPCS 51300.

## Instructions

If running on the UChicago CS Linux servers, first run `module load golang/1.21.7` to load the proper version of Go. The golitec executable can be created by running `source build.sh` in the root directory, and the executable can be run with `./golitec <flags...> <PATH>`, where PATH is a .golite file or path to a .golite file.

Example:

```
module load golang/1.21.7
source build.sh
./golitec benchmarks/simple.golite
```

## Compiler Flags

An explanation of the various compiler flags and their meanings.

`-l`: The compiler prints out the tokens of the given program and the execution stops at the end of lexing.

`-ast`: The compiler prints out the AST of the given program before continuing.

`-llvm-ir=[stack|reg]`: The LLVM translation is written to a `.ll` file with the same name and location as the input file. If the input file were located at `/example/test.golite` then the LLVM output would be written to `/example/test.ll`. If `-llvm-ir=stack` is specified, then the stack-based IR is generated, and if `-llvm-ir=reg` is specified then the register-based IR is generated.

`-target=STRING`: Changes the LLVM target triple value to STRING. The default is `"x86_64-linux-gnu"`.

`-S`: The assembly translation is written to a `.s` file with the same name and location as the input file. If the input file were located at `/example/test.golite` then the assembly output would be written to `/example/test.s`.

Note that the execution of the compiler will terminate if it fails lexing, parsing, or semantic analysis.
