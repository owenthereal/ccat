page_title: Scanner Output

# Scanner Output

The scanner should descend through the directory on which it was invoked,
searching for source units and printing its results to stdout.

## Output Schema

The data in the output file should consist of a JSON array of objects, each of which conforms
to the data in this struct.

[[.code "unit/source_unit.go" "SourceUnit"]]

## Example: Scanner output on [gorilla/mux](https://github.com/gorilla/mux)

```json
[
    {
        "Name": "github.com/gorilla/mux",
        "Type": "GoPackage",
        "Repo": "",
        "Globs": null,
        "Files": [
            "doc.go",
            "mux.go",
            "regexp.go",
            "route.go",
            "bench_test.go",
            "mux_test.go",
            "old_test.go"
        ],
        "Dependencies": [
            "bytes",
            "errors",
            "fmt",
            "github.com/gorilla/context",
            "net/http",
            "net/url",
            "path",
            "regexp",
            "strings",
            "testing"
        ],
        "Data": {
            "Dir": ".",
            "Name": "mux",
            "Doc": "Package gorilla/mux implements a request router and dispatcher.",
            "ImportPath": "github.com/gorilla/mux",
            "Root": "",
            "SrcRoot": "",
            "PkgRoot": "",
            "BinDir": "",
            "Goroot": false,
            "PkgObj": "",
            "AllTags": null,
            "ConflictDir": "",
            "GoFiles": [
                "doc.go",
                "mux.go",
                "regexp.go",
                "route.go"
            ],
            "CgoFiles": null,
            "IgnoredGoFiles": null,
            "CFiles": null,
            "CXXFiles": null,
            "MFiles": null,
            "HFiles": null,
            "SFiles": null,
            "SwigFiles": null,
            "SwigCXXFiles": null,
            "SysoFiles": null,
            "CgoCFLAGS": null,
            "CgoCPPFLAGS": null,
            "CgoCXXFLAGS": null,
            "CgoLDFLAGS": null,
            "CgoPkgConfig": null,
            "Imports": [
                "bytes",
                "errors",
                "fmt",
                "github.com/gorilla/context",
                "net/http",
                "net/url",
                "path",
                "regexp",
                "strings"
            ],
            "ImportPos": null,
            "TestGoFiles": [
                "bench_test.go",
                "mux_test.go",
                "old_test.go"
            ],
            "TestImports": [
                "bytes",
                "fmt",
                "github.com/gorilla/context",
                "net/http",
                "testing"
            ],
            "TestImportPos": null,
            "XTestGoFiles": null,
            "XTestImports": null,
            "XTestImportPos": null
        },
        "Ops": {
            "depresolve": null,
            "graph": null
        }
    }
]
```
