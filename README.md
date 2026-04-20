# dupefinder

A concurrent duplicate file finder CLI written in Go. Walks a directory tree, hashes files
using a configurable worker pool, and reports duplicates grouped by content hash.

----

```
go run . find ./testdata --min-size 1B
```