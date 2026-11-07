# go-mongo

Thin, testable MongoDB client built on the official [mongo-driver v2](https://pkg.go.dev/go.mongodb.org/mongo-driver/v2).

## Install

```bash
go get github.com/andywang2011/go-mongo@v0.1.0
```

## Quick start

```go
import (
	"context"
	"fmt"

	"github.com/andywang2011/go-mongo"
)

cfg := mongo.ConfigDefault()
cfg.Database = "app"
cfg.Address = "127.0.0.1:65530"

client, err := mongo.New(context.Background(), cfg)
if err != nil {
	fmt.Println("unavailable")
	return
}
defer client.Close(context.Background())
fmt.Println("ok")
```

## Examples

```bash
export MONGO_URI='mongodb://127.0.0.1:27017'
export MONGO_DB=app
cd go-mongo/examples/basic && go run .
cd go-mongo/examples/replica && go run .
```

## License

MIT — see [LICENSE](LICENSE).
