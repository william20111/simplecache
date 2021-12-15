# simplecache

Simplecache is a key value in memory cache for applications that do not need a distributed cache.

## Installation

```shell
go get github/com/william20111/simplecache
```

## Usage

```go

import (
    "github.com/william20111/simplecache"
)

package main

func main() {
    cache := simplecache.New(1000)
    cache.Set("example", "test123", 1 * time.Hour)
}
```
