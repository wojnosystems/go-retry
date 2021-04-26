# Overview

A thread-safe, minimalistic yet flexible retry library for GoLang.

I wanted to use this for database connections, so I'd need to attempt a query or connection with the same retry configuration multiple times. I'd configure this once and re-use it for each request. I also wanted it to be re-usable and thread-safe.

Most of the errors returned by MySQL/Postgres aren't retryable, like query formatting issues or missing data. The only time I really wanted to retry is if there was a networking timeout. Therefore, the default is to stop retrying on any error, unless `retry.Again` is returned. In this case, it will be retried unless we're at the limits. If you return `retry.Success`, then iteration stops and returns no error.

# Installing

```shell
go get github.com/wojnosystems/go-retry
```

# How do I use it?

Pretty simple, you simply pick which delay method works for you and pass in your function to retry. You can find a list of officially supported methods under the "retry" package. See below for examples.

The function you want to retry can do anything it wants within it. However, you control the retry logic based on the return value of your function.

You can return 3 types of errors:

* **nil AKA retryStop.Success:** this indicates that the attempt succeeded and should not be retried. nil is returned from the Retry method
* **retryAgain.Error(ErrSomeError):** wrap any errors in this method to trigger a retry. If you exceed the retries, the error passed to retryAgain.Error will be returned to the caller
* **any other error:** will indicate a non-retryable error. No retries will be attempted, this error will be returned immediately to the caller without any delays

# Examples

## Network I/O

Here's an example of how you can use this to connect and send data to a TCP socket and retry on failure. This example shows you how to use an exponential back-off and what happens when your service fails permanently.

```go
package main

import (
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryAgain"
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
			socket, dialErr := net.Dial("tcp", "localhost:9999")
			if dialErr != nil {
				// all dialErrs are retried
				return retryAgain.Error(dialErr)
			}

			// Write errors are NOT retried
			_, writeErr := socket.Write([]byte("some payload"))

			// if writeErr is nil, success!
			// if writeErr is not wrapped in retry.Error, retry will stop retrying and return the
			// error to the caller
			return writeErr
		})

		// Outputs the Dial error because we ran out of retries
		fmt.Println(err)
	})

	fmt.Println("total time", totalTime)
}
```

common.TimeThis and common.NewTimeSet are helper methods that record time differences. They're not involved in the retry logic and only serve to help you understand how attempts and delays between attempts work.

In the above example, dialer is a retry configuration. It tells the library how it should retry your function.

Outputs:

```
dialing 296ns
dialing 50.47666ms
dialing 100.40726ms
dialing 200.546019ms
dialing 400.66924ms
dialing 500.740532ms
dialing 500.799248ms
dialing 500.573812ms
dial tcp 127.0.0.1:9999: connect: connection refused
total time 2.254906137s
```

Because I have no service running on port 9999 on my localhost, this emulates a network timeout. You can see that this retries 8 times, exponentially backing off until it reaches 500ms, at which point it caps out and will not exceed the MaxWaitBetweenAttempts

## Retry With Cap

Retries something up to MaxAttempts times, waiting the same amount of time between each request

```go
package main

import (
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryAgain"
	"time"
)

func main() {
	tries := 0
	duration := common.TimeThis(func() {

		timer := common.NewTimeSet()

		_ = (&retry.UpTo{
			WaitBetweenAttempts: 10 * time.Millisecond,
			MaxAttempts:         10,
		}).Retry(func() (err error) {
			fmt.Println(timer.SinceLast())
			tries++
			return retryAgain.Error(errors.New("simulated error"))
		})

	})
	fmt.Println("tried", tries, "times taking", duration)
}
```

Outputs:

```
162ns
10.187236ms
10.155135ms
10.14498ms
10.163035ms
10.163972ms
10.166715ms
10.139126ms
10.160669ms
10.158786ms
tried 10 times taking 91.453743ms
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
	"github.com/wojnosystems/go-retry/retryAgain"
	"time"
)

func main() {
	tries := 0
	duration := common.TimeThis(func() {

		timer := common.NewTimeSet()

		_ = (&retry.ExponentialMaxWaitUpTo{
			InitialWaitBetweenAttempts: 10 * time.Millisecond,
			GrowthFactor:               1.5,
			MaxAttempts:                10,
			MaxWaitBetweenAttempts:     100 * time.Millisecond,
		}).Retry(func() (err error) {
			tries++
			fmt.Println(timer.SinceLast())
			return retryAgain.Error(errors.New("simulated error"))
		})

	})
	fmt.Println("tried", tries, "times taking", duration)
}
```

Outputs:

```
207ns
10.152764ms
25.173204ms
63.17935ms
100.326124ms
100.279685ms
100.344036ms
100.340817ms
100.194532ms
100.22337ms
tried 10 times taking 700.228581ms
```

## Retry Forever

This will retry without limit until you either return a success or you return an non-retryable error.

```go
package main

import (
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryAgain"
	"github.com/wojnosystems/go-retry/retryStop"
	"time"
)

func main() {
	tries := 0
	duration := common.TimeThis(func() {

		timer := common.NewTimeSet()

		_ = (&retry.Forever{
			WaitBetweenAttempts: 10 * time.Millisecond,
		}).Retry(func() (err error) {
			fmt.Println(timer.SinceLast())
			if tries < 10 {
				tries++
				return retryAgain.Error(errors.New("simulated error"))
			}
			return retryStop.Success
		})
	})
	fmt.Println("tried", tries, "times taking", duration)
}
```

Outputs:

```
273ns
10.145787ms
10.199774ms
10.141523ms
10.176248ms
10.290984ms
10.232297ms
10.21017ms
10.164733ms
10.158418ms
10.189141ms
tried 10 times taking 101.930106ms
```

# Swappable Controls

Because the configuration of the retries are all interfaces, you can swap them out if you want to control how things are being retried. This is useful if, for example, you have a circuit breaker that should not retry while in the open state.

```go
package main

import (
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryAgain"
	"time"
)

func main() {
	normal := &retry.LinearUpTo{
		InitialWaitBetweenAttempts: 100 * time.Millisecond,
		GrowthFactor:               1,
		MaxAttempts: 				5,
	}

	var strategy retry.Retrier
	strategy = normal

	_ = strategy.Retry(func() (err error) {
		fmt.Println("normal")
		return retryAgain.Error(errors.New("some error"))
	})

	strategy = retry.Never

	_ = strategy.Retry(func() (err error) {
		fmt.Println("NEVER")
		return retryAgain.Error(errors.New("some error"))
	})

	strategy = normal
	_ = strategy.Retry(func() (err error) {
		fmt.Println("normal")
		return retryAgain.Error(errors.New("some error"))
	})
}
```

Outputs:

```
normal
normal
normal
normal
normal
NEVER
normal
normal
normal
normal
normal

```
