# Overview

A thread-safe, minimalistic yet flexible retry library for GoLang.

I wanted to use this for database connections, so I'd need to attempt a query or connection with the same retry configuration multiple times. I'd configure this once and re-use it for each request. I also wanted it to be re-usable and thread-safe.

Most of the errors returned by MySQL/Postgres aren't retryable, like query formatting issues or missing data. The only time I really wanted to retry is if there was a networking timeout. Therefore, the default is to stop retrying on any error, unless `retry.Again` is returned. In this case, it will be retried unless we're at the limits. If you return `retry.Success`, then iteration stops and returns no error.

# Installing

```shell
go get github.com/wojnosystems/go-retry
```

# Examples

## Network I/O

```go
package main

import (
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"net"
	"time"
)

func main() {
	dialer := &retry.ExponentialMaxWaitUpTo{
		InitialWaitBetweenAttempts: 50 * time.Millisecond,
		GrowthFactor:               1.0,
		MaxAttempts:                8,
		MaxWaitBetweenAttempts:     500 * time.Millisecond,
	}

	timer := common.NewTimeSet()

	totalTime := common.TimeThis(func() {
		err := dialer.Retry(func() error {
			fmt.Println("dialing", timer.SinceLast())
			socket, dialErr := net.Dial("tcp","localhost:9999" )
			if dialErr != nil {
				// all dialErrs are retried
				return retry.Again(dialErr)
			}

			// Write errors are NOT retried
			_, writeErr := socket.Write([]byte("some payload"))

			// if writeErr is nil, success!
			// if writeErr is not wrapped in retry.Again, retry will stop retrying and return the
			// error to the caller
			return writeErr
		})

		// Outputs the Dial error because we ran out of retries
		fmt.Println(err)
	})

	fmt.Println("total time", totalTime)
}
```

Outputs:

```
dialing 143ns
dialing 50.521628ms
dialing 100.403192ms
dialing 200.451139ms
dialing 400.775837ms
dialing 501.011315ms
dialing 501.030898ms
dialing 500.591912ms
dial tcp 127.0.0.1:9999: connect: connection refused
total time 2.254955255s

```

## Retry With Cap

Retries something up to MaxAttempts times, waiting the same amount of time between each request

```go
package main

import (
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"time"
)

func main() {
	tries := 0
	duration := common.TimeThis(func() {

		timer := common.NewTimeSet()

		_ = (&retry.UpTo{
			WaitBetweenAttempts: 10 * time.Millisecond,
			MaxAttempts: 10,
		}).Retry(func() (err error) {
			fmt.Println(timer.SinceLast())
			tries++
			return retry.Again(errors.New("simulated error"))
		})

	})
	fmt.Println("tried", tries, "times taking", duration)
}
```

Outputs:

```
142ns
10.175253ms
10.190037ms
10.211729ms
10.217857ms
10.163333ms
10.137719ms
10.139279ms
10.145087ms
10.138009ms
tried 10 times taking 91.530497ms
```

## Retry Exponential With Max Time Between Request and Cap

Retries something up to MaxAttempts times, waiting exponentially longer times between requests until it hits a certain time limit.

```go
package main

import (
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"time"
)

func main() {
	tries := 0
	duration := common.TimeThis(func() {

		timer := common.NewTimeSet()

		_ = (&retry.ExponentialMaxWaitUpTo{
			InitialWaitBetweenAttempts: 10 * time.Millisecond,
			GrowthFactor: 1.5,
			MaxAttempts: 10,
			MaxWaitBetweenAttempts: 100 * time.Millisecond,
		}).Retry(func() (err error) {
			tries++
			fmt.Println(timer.SinceLast())
			return retry.Again(errors.New("simulated error"))
		})

	})
	fmt.Println("tried", tries, "times taking", duration)
}
```

Outputs:

```
134ns
10.180866ms
25.207167ms
63.159911ms
100.26247ms
100.253ms
100.375093ms
100.368574ms
100.166419ms
100.308456ms
tried 10 times taking 700.31928ms
```

## Retry Forever

```go
package main

import (
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"time"
)

func main() {
	tries := 0
	duration := common.TimeThis(func() {

		timer := common.NewTimeSet()

		_ = (&retry.Forever{
			WaitBetweenAttempts: 10*time.Millisecond,
		}).Retry(func() (err error) {
			fmt.Println(timer.SinceLast())
			if tries < 10 {
				tries++
				return retry.Again(errors.New("simulated error"))
			}
			return retry.Success
		})
	})
	fmt.Println("tried", tries, "times taking", duration)
}
```

Outputs:

```
88ns
10.150113ms
10.32446ms
10.274221ms
10.205891ms
10.22437ms
10.183359ms
10.137373ms
10.22798ms
10.214048ms
10.200296ms
tried 10 times taking 102.159196ms
```
