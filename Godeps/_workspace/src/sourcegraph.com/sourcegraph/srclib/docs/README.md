# srclib Documentation Site

The srclib documentation files are written in Markdown and converted to a
browsable HTML static site, hosted at [srclib.org](http://srclib.org/).
**[MkDocs](http://www.mkdocs.org/)**, installable through pip, is used to
generate the HTML.

Run `pip install mkdocs` to install MkDocs.

You will also need SASS, in order to compile the style sheets. This can be
obtained with `gem install sass`. In order to use the testing script `test.sh`,
you will need to run `sudo apt-get install inotify-tools`.

## WARNING

**Doc tags (like "START SomeTag") are NOT preserved when protobuf
files are generated. You will need to add them by hand.**


## Testing Documentation

To start testing, run `./test.sh`. This will watch for changes, building into
`site/` whenever a file is modified. Additionally, `python -m SimpleHTTPServer`
is used to start a local server on port 8000.

## Editing

Markdown files are located in the `sources` directory, with subdirectories
intended to correspond with the drop-down menus on the site. The actual
categories, however, are specified through configuration in `mkdocs.yml`.

### Embedding Code Segments
In the docs, you can embed code taken directly from the repository source.
The format is as follows:
```
[[.code "grapher/grapher.go" 23 30]]
```
Note that the file name is specified as relative to the repository root.
You can also specify two strings instead of two integers.
```
[[.code "grapher/grapher.go" "START REGION" "END REGION"]]
```
This will look for the first lines that contain those strings (case insensitive).

Like in the go-present tool, lines ending with the word `OMIT` are omitted from the
embedded code. If you have a segment as follows:
```go
// START RegionName OMIT
Go code here...
// END RegionName OMIT
```
You can use the following short-hand:
```
[[.code "grapher/grapher.go" "RegionName"]]
```

## Deploying

The site can be deployed to srclib.org simply by running `./deploy.sh`. This
builds the SASS stylesheet, builds the docs, and builds the other static pages
of the site. It then uses `aws s3 sync` to store these in the correct bucket.
