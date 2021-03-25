<img src="https://github.com/raphaelreyna/go-comptop/raw/master/logo/logo.png" height="256px">

# go-comptop [![](https://godoc.org/github.com/raphaelreyna/go-comptop?status.svg)](http://godoc.org/github.com/raphaelreyna/go-comptop) [![Go Report Card](https://goreportcard.com/badge/github.com/raphaelreyna/go-comptop)](https://goreportcard.com/report/github.com/raphaelreyna/go-comptop)

A computational topology package for gophers.

## Features
Simplices, simplicial complexes, simplicial chains, chain groups, simplicial sets, methods for computing boundaries, methods for computing Euler characteristics, Euler integration, and more (with even more to come)!

## Contributing
Contributions are welcome!

## Examples:
### Counting requests to load-balanced services with poor observability:
#### Scenario:
Suppose we have a backend with the network topology shown below.

- Services 0, 1, 2 are load balanced and accessed by 3.
- Services 6 and 7 are load balanced and accessed by 5.
- Services 8 and 9 are load balanced and accessed by 10. 

Let's assume we are using a load balancing scheme where requests are multicast to some random number of randomly selected instances.

Services 4 and 10 are the only public facing services.
Services forward client API keys with requests to internal services.

The services in red (0, 1, 2, 5, 6, 7, 9) are the only ones logging API keys, but without timestamps.
Logs are rotated daily.

We can query the services in red for the number of times it logged a clients API key today.

<img src="https://github.com/raphaelreyna/go-comptop/raw/master/examples/microservices//images/labels.png" height="256px">

#### Problem:
Suppose that client Alice made 3 requests today:
- A request to 4 which sent a request to 3, which load balanced to 0 and 1.
- A request to 10 which load balanced to 9 then 10.
- A request to 4 which load balanced to 3, which load balanced to 0, 1 and 2.

We don't know that Alice made 3 requests or the path thos requests took. All we can do is query the services and get back the results (a heightmap) shown below.


How many requests did client Alice make today?


<img src="https://github.com/raphaelreyna/go-comptop/raw/master/examples/microservices//images/heightMap.png" height="256px">

#### Solution:
We can integrate the height map over the network with respect to the Euler characteristic.
This gives us an estimate on the number of requests Alice made.
```go
package main

import (
  "fmt"
  "github.com/raphaelreyna/go-comptop"
)

func main() {
  c := &comptop.Complex{}
  
	loggingNetwork := []comptop.Base{
		{0, 1, 2},
		{5, 6, 7},
    {5, 9},
	}
	c.NewSimplices(loggingNetwork...)

	heightMap := map[comptop.Index]int{
		0: 2, 1: 2, 2: 1,
		5: 1, 6: 0, 7: 0,
		9: 1,
	}

	f := comptop.CF(func(idx comptop.Index) int {
		return heightMap[idx]
	})

	fmt.Println(c.EulerIntegral(0, 2, f))
}
```

We correctly get 3 as the answer.
