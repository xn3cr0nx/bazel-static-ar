# gAR
gAR is a tool written in Go that uses ar to merge recursevily all object files produced in the specified directory in order to easily link the static library.

## Usage
> go get https://github.com/xn3cr0nx/gAR

You can specify the directory (-d/--dir) where gAR should look for object files.
You can specify the output pathname (-o/--out) of the file .a gAR will produce.

### Example
*gar --dir ./libs/project/ -o ./libs/include/out.a*
