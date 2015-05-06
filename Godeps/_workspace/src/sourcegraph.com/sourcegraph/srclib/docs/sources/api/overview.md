page_title: Overview

# Overview
<div class="alert alert-danger" role="alert">Note: The API is still in flux, and may change throughout the duration of this beta.</div>

The srclib API will be used through the invocation of subcommands of the `src` executable.

## API Commands

API commands return their responses as JSON, to facilitate the building of tools on top of srclib. Sourcegraph's [plugins](#TODO-plugins-overview) all make heavy use of the API commands.

<div class="alert alert-danger" role="alert">Note: The docs currently only show the Go representation of the output. See <a href="https://blog.golang.org/json-and-go">this blog post</a> for a primer on how Go types are marshaled into JSON.</div>

<!-- TODO: This should be generated from 'commands' in mkdocs.yml -->

### `src api describe`

[[.doc "src/api_cmds.go" "APIDescribeCmdDoc"]]

#### Usage

[[.run src api describe -h]]

#### Output

[[.doc "src/api_cmds.go" "APIDescribeCmdOutput"]]

### `src api list`
[[.doc "src/api_cmds.go" "APIListCmdDoc"]]

#### Usage
[[.run src api list -h]]

#### Output
[[.doc "src/api_cmds.go" "APIListCmdOutput"]]

### `src api deps`
[[.doc "src/api_cmds.go" "APIDepsCmdDoc"]]

#### Usage
[[.run src api list -h]]

#### Output
[[.doc "src/api_cmds.go" "APIDepsCmdOutput"]]

### `src api units`
[[.doc "src/api_cmds.go" "APIUnitsCmdDoc"]]

#### Usage
[[.run src api units -h]]

#### Output
[[.doc "src/api_cmds.go" "APIUnitsCmdOutput"]]

## Standalone Commands

Standalong commands are for the srclib power user: most people will use srclib through an editor plugin or Sourcegraph, but the following commands are useful for modifying the state of a repository's analysis data.

### `src config`
`src config` is used to detect what kinds of source units (npm/pip/Go/etc. packages) exist in a repository or directory tree.

### `src make`
`src make` is used to perform analysis on a given directory. See the [src make docs](make.md) for usage instructions.
