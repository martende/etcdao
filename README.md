# etcdao

etcdao is the Go library for etcd, that implements deserealisation of data from etcd for 
some simple types.

## Currently supported types:

* Type primitives: int,string,bool
* time.Time
* struct with supported types
* slice 
* map with string as key

## Extensions:

Fields in struct can be configured via tags.
By default , module finds a variable with tha same name, but it can be changed via 'name' tag ( see usage example )

For time.Time 'format' configures time format in usual time.Parse format. Default formatting for time is "2006-01-02"

* Maps are encoded in their native format:

```
map[string] int {"A":1,"B":1}
```

awaits in etcd:

```
/objectPath/A -> 1
/objectPath/B -> 1
```

* Slices are using the following notation

```

[]int{1,2,3}

/slicePath/0 -> 1
/slicePath/1 -> 2
/slicePath/2 -> 3
```


## Instalation

```
go get github.com/martende/etcdao
```

Tests use etcd service on a localhost machine, if there is no etcd they fail, to skip integration tests use 

```
go test -short github.com/martende/etcdao
```


## Usage

```go


package main

import (
	"log"
	"time"
	"github.com/martende/etcdao"
	"golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

type Config2 struct {
	A map[string] int
	B []int
	C string	`name:"attributesFile"`
	D time.Time	`name:"updateTime" format:"2006-01-02 15:04:05"`
} 

type Config struct {
	A map[string] int
	B []int
	C Config2
}

func main() {
	cfg := client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := client.NewKeysAPI(c)

	configVal := A{}

	err = etcdao.ReadObject(kapi, context.Background(), "/config", &configVal)

}


```