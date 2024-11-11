# local-search-engine

## Usage

```sh
# install deps
go mod tidy

# run
go run ./cmd <dir to index> <search terms>

# example
go run ./cmd ./testdir 'workspace height'
# returns filepath with terms' row:col in that file
```
