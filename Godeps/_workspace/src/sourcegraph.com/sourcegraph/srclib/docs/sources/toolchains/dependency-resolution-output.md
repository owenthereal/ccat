page_title: Dependency Resolution Output

# Dependency Resolution Output

The output of the dependency resolution tool should be printed to standard
output.

## Output Schema

The schema of the dependency resolution output should be an array of
`Resolution` objects, with the structure of each object being as follows.

[[.code "dep/resolve.go" "Resolution"]]

The `Raw` field is language specific, but the Target field follows the following format.

[[.code "dep/resolve.go" "ResolvedTarget"]]

If an error occurred during resolution, a detailed description should be placed in the `Error` field.

## Example: Depresolve on [gorilla/mux](https://github.com/gorilla/mux)

```json
[
    {
        "Raw": "bytes",
        "Target": {
            "ToRepoCloneURL": "",
            "ToUnit": "bytes",
            "ToUnitType": "GoPackage",
            "ToVersionString": "",
            "ToRevSpec": ""
        }
    },
    {
        "Raw": "errors",
        "Target": {
            "ToRepoCloneURL": "",
            "ToUnit": "errors",
            "ToUnitType": "GoPackage",
            "ToVersionString": "",
            "ToRevSpec": ""
        }
    },
    {
        "Raw": "fmt",
        "Target": {
            "ToRepoCloneURL": "",
            "ToUnit": "fmt",
            "ToUnitType": "GoPackage",
            "ToVersionString": "",
            "ToRevSpec": ""
        }
    },
    {
        "Raw": "github.com/gorilla/context",
        "Target": {
            "ToRepoCloneURL": "",
            "ToUnit": "github.com/gorilla/context",
            "ToUnitType": "GoPackage",
            "ToVersionString": "",
            "ToRevSpec": ""
        }
    },
    {
        "Raw": "net/http",
        "Target": {
            "ToRepoCloneURL": "",
            "ToUnit": "net/http",
            "ToUnitType": "GoPackage",
            "ToVersionString": "",
            "ToRevSpec": ""
        }
    },
    {
        "Raw": "net/url",
        "Target": {
            "ToRepoCloneURL": "",
            "ToUnit": "net/url",
            "ToUnitType": "GoPackage",
            "ToVersionString": "",
            "ToRevSpec": ""
        }
    },
    {
        "Raw": "path",
        "Target": {
            "ToRepoCloneURL": "",
            "ToUnit": "path",
            "ToUnitType": "GoPackage",
            "ToVersionString": "",
            "ToRevSpec": ""
        }
    },
    {
        "Raw": "regexp",
        "Target": {
            "ToRepoCloneURL": "",
            "ToUnit": "regexp",
            "ToUnitType": "GoPackage",
            "ToVersionString": "",
            "ToRevSpec": ""
        }
    },
    {
        "Raw": "strings",
        "Target": {
            "ToRepoCloneURL": "",
            "ToUnit": "strings",
            "ToUnitType": "GoPackage",
            "ToVersionString": "",
            "ToRevSpec": ""
        }
    },
    {
        "Raw": "testing",
        "Target": {
            "ToRepoCloneURL": "",
            "ToUnit": "testing",
            "ToUnitType": "GoPackage",
            "ToVersionString": "",
            "ToRevSpec": ""
        }
    }
]
```
