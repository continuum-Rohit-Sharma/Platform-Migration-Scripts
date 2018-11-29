## jsonMerge Utility
Merge 2 JSON Trees into a single structure.

## Build
(from JSON Merge Folder)

$ go build .

## Usage
$ jsonMerge < source-json-file> < delta-json-file> [< destination-json-file>]

If the destination file is not specified, the source json file is overwritten.
