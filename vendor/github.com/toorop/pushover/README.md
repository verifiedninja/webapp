# pushover

pushover implements a go package client interface to the [pushover API](https://pushover.net/api).

## Documentation
See [godoc](http://godoc.org/github.com/thorduri/pushover)

## Example

Assuming pushover.go as:
```Go
package main

import (
	"fmt"
	"log"

	"github.com/thorduri/pushover"
)

const exampleToken = "KzGDORePK8gMaC0QOYAMyEEuzJnyUi"
const exampleUser = "uQiRzpo4DXghDmr9QzzfQu27cmVRsG"

func main() {
	po, err := pushover.NewPushover(exampleToken, exampleUser)
	if err != nil {
		log.Fatal(err)
	}

	err = po.Message("Hello Pushover!")
	if err != nil {
		log.Fatal(err)
	}
}
```
Then:
```bash
$ go run pushover.go
```
Will send a the message "Hello Pushover!" over the pushover API using default values for the provided token and user.
