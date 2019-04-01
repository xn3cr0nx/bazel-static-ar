# bazelar
Bazelar is a tool written in Go that uses ar to merge all archive files produced by bazel in order to easily link the library using cmake.

## Usage
> go get https://github.com/xn3cr0nx/bazelar

You can specify the directory (-d/--dir) where bazelar should look for bazel-bin folder.
You can specify the output pathname (-o/--out) of the file .a bazelar will produce.

### Example
*bazelar --dir ./libs/project/ -o ./libs/include/out.a*
