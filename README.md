# Overview

A thread-safe, minimalistic yet flexible retry library for GoLang with context support.

I wanted to use this for database connections, so I'd need to attempt a query or connection with the same retry configuration multiple times. I'd configure this once and re-use it for each request. I also wanted it to be re-usable and thread-safe.

I also wanted to be sure that any time the retry would sleep after an error, the sleep would never exceed the deadline of the context, which is useful for servers. For example, if each endpoint in your server had 30 seconds to perform its tasks, and you have an exponential back off that would sleep and cause the request to take 45 seconds, the library will stop sleeping at the context deadline instead of allowing the sleep to continue for 15 additional seconds.

When the deadline is exceeded, your method will be executed one more time. If your code also uses the same context, the error will be `context.DeadlineExceeded`, which you should never wrap in `retryAgain.Error`.

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

Here's an example of how you can use this to connect with an HTTP request and retry on failure. This example shows you how to use an exponential back-off and what happens when your service fails permanently.

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryAgain"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	dialer := &retry.ExponentialMaxWaitUpTo{
		InitialWaitBetweenAttempts: 50 * time.Millisecond,
		GrowthFactor:               1.0,
		MaxAttempts:                15,
		MaxWaitBetweenAttempts:     500 * time.Millisecond,
	}

	timer := common.NewTimeSet()

	totalTime := common.TimeThis(func() {
		err := dialer.Retry(ctx, func() error {
			fmt.Println("getting", timer.SinceLast())
			req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/non-existent", nil)
			req = req.WithContext(ctx)
			_, getErr := http.DefaultClient.Do(req)
			if getErr != nil {
				getErr = errors.Unwrap(getErr)
				if getErr != context.DeadlineExceeded {
					getErr = retryAgain.Error(getErr)
				}
				return getErr
			}
			return nil
		})

		// Outputs the http error because we ran out of retries
		fmt.Println(err)
	})

	fmt.Println("total time", totalTime)
}
```

common.TimeThis and common.NewTimeSet are helper methods that record time differences. They're not involved in the retry logic and only serve to help you understand how attempts and delays between attempts work.

The context controls how long the retry will wait as well. If the last request failed and the library would have slept, the sleep should not sleep much longer than the context deadline. It will not, of course, be perfect. However, it should help prevent the retry library from sleeping for an unreasonably long time after your context expires.

### Outputs

```
getting 573ns
getting 50.650745ms
getting 100.387843ms
getting 200.851694ms
getting 400.830581ms
getting 500.871094ms
getting 500.916346ms
getting 245.670663ms
context deadline exceeded
total time 2.000234103s
```

Because I have no service running on port 9999 on my localhost, this emulates a network timeout. You can see that this retries 10 times, exponentially backing off until it reaches 500ms, at which point it caps out and will not exceed the MaxWaitBetweenAttempts.

As you can see, after the context deadline of 2 seconds is exceeded, a 500ms sleep is interrupted and reduced to 245.66ms. The retrier then continues to attempt to execute the code again. Since Request is context-aware, the http.Client will not execute the request. We catch the DeadlineExceeded error and return it immediately to prevent the retry library from executing again.

## Retry With Cap

Retries something up to MaxAttempts times, waiting the same amount of time between each request

```go
package main

import (
	"context"
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
		}).Retry(context.TODO(), func() (err error) {
			fmt.Println(timer.SinceLast())
			tries++
			return retryAgain.Error(errors.New("simulated error"))
		})

	})
	fmt.Println("tried", tries, "times taking", duration)
}
```

### Outputs

```
453ns
10.214257ms
10.251075ms
10.161697ms
10.229908ms
10.211922ms
10.150657ms
10.228492ms
10.142921ms
10.162766ms
tried 10 times taking 91.767877ms
```

## Retry Exponential With Max Time Between Request and Cap

Retries something up to MaxAttempts times, waiting exponentially longer times between requests until it hits a certain time limit.

```go
package main

import (
	"context"
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
		}).Retry(context.TODO(), func() (err error) {
			tries++
			fmt.Println(timer.SinceLast())
			return retryAgain.Error(errors.New("simulated error"))
		})

	})
	fmt.Println("tried", tries, "times taking", duration)
}
```

### Outputs

```
149ns
10.208347ms
25.296481ms
63.249355ms
100.283166ms
100.27994ms
100.273478ms
100.330776ms
100.256766ms
100.269877ms
tried 10 times taking 700.496801ms
```

## Retry Forever

This will retry without limit until you either return `Success` or you return a non-retryable error.

```go
package main

import (
	"context"
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
		}).Retry(context.TODO(), func() (err error) {
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

### Outputs

```
150ns
10.174455ms
10.220227ms
10.202718ms
10.212075ms
10.140283ms
10.182415ms
10.165957ms
10.144869ms
10.16555ms
10.139993ms
tried 10 times taking 101.768423ms
```

# Swappable Controls

Because the configuration of the retries are all interfaces, you can swap them out if you want to control how things are being retried. This is useful if, for example, you have a circuit breaker that should not retry while in the open state.

```go
package main

import (
	"context"
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
		MaxAttempts:                5,
	}

	var strategy retry.Retrier
	strategy = normal

	_ = strategy.Retry(context.TODO(), func() (err error) {
		fmt.Println("normal")
		return retryAgain.Error(errors.New("some error"))
	})

	strategy = retry.Never

	_ = strategy.Retry(context.TODO(), func() (err error) {
		fmt.Println("NEVER")
		return retryAgain.Error(errors.New("some error"))
	})

	strategy = normal
	_ = strategy.Retry(context.TODO(), func() (err error) {
		fmt.Println("normal")
		return retryAgain.Error(errors.New("some error"))
	})
}
```

### Outputs

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
