# Overview

Retry blocks of code, limited by context deadlines & cancellations as well as multiple back off strategies and retry attempt limits. Recover the last error that triggered a retry or errors that should not be retried immediately.

Any time the retry would sleep after an error, the sleep should never excessively exceed the deadline of the context. This is useful for use in servers which tend to have per-request deadlines. For example, if each endpoint in your server had 30 seconds to perform its tasks, and you have an exponential back off that would sleep and cause the request to take 45 seconds, the library will stop sleeping at the context deadline instead of allowing the sleep to continue for 15 additional seconds.

When the deadline is exceeded or the context is cancelled, your method _may_ be executed one more time. However, if just before the callback is invoked, the context is not yet Done, it will not be run. If your code also uses the same context, the error will be `context.DeadlineExceeded`. Never wrap in `retryError.Again()` otherwise retry could infinite loop (depending on which strategy you select). Should your callback complete and the context becomes Done, the wait duration will be curtailed as it is also restricted by contexts. This should ensure that retry-able blocks don't exceed the context by an unreasonable factor and make very reliable, and controllable code.

Example use-case: Most of the errors returned by MySQL/Postgres aren't retryable. Query formatting issues or missing data, for example. The only time I really wanted to retry is if there was a network timeout or disconnect. Therefore, the default is to stop retrying on any error, unless `retry.Again` is returned. In this case, it will be retried unless we're at the limits. If you return `retry.Success`, then iteration stops and returns no error. The library can be made to invert this, you just need to wrap any errors to retry in retryError.Again.

These errors are returned to the calling code, so you can take a specific action in response to a specific error value.

# Installing

```shell
go get -u github.com/wojnosystems/go-retry
```

# How do I use it?

Pick which strategy works for you and pass in your function to retry. You can find a list of officially supported methods under the "retry" package and below. See below for examples.

The function you want to retry can do anything it wants within it. However, you control the retry logic based on the return value of your function.

You can return 3 types of errors:

* **nil AKA retryError.StopSuccess:** this indicates that the attempt succeeded and should not be retried. nil is returned from the `Retry` method
* **retryError.Again(ErrSomeError):** wrap any errors in this method to trigger a retry. If you exceed the retries, the error passed to retryError.Again will be returned to the caller of `Retry` without the wrapper
* **any other error:** will indicate a non-retryable error. No retries will be attempted, this error will be returned immediately to the caller of `Retry` without any waiting

I opted to not retry for errors not explicitly marked to be retried in order to allow only certain errors to be retried. I think this makes this retry library a bit safer as we're only changing how the logic operates if the developer explicitly requests a retry.

Developers are encouraged to make functions that take in `error` and only wrap it in `retryError.Again()` should they decide that it's appropriate to retry.

## Strategies

There are several retry strategies available form this library for common retry patterns in increasing complexity. You can find all of these under the `retry` package

* **Skip**: will never execute the callback and will return retryError.StopSuccess (nil). Useful for testing and mocking things you want to test that don't actually do anything.
* **Never**: will never retry, but will execute the callback exactly once. Useful for testing and mocking when you just want to run something once
* **UpTo**: will run the callback until it succeeds, returns a non-retryable error, or runs out of maximum attempts, waiting the same time between attempts
* **Linear**: will run the callback until it succeeds or returns a non-retryable error, but waiting in increasing amount of time between attempts that grows linearly, see the struct's documentation for the formula
* **LinearUpTo**: same as Linear and UpTo, the wait time increases and the number of waits is now bounded
* **LinearMaxWaitUpTo**: same as LinearUpTo, but the maximum wait time is capped so that if your wait times grow too large, you can set a bound on the wait time's growth
* **Exponential**: same as Linear, but the wait time grows exponentially. See the struct's documentation for the formula
* **ExponentialUpTo**: Same as exponential, the wait time increases and the number of waits is now bounded
* **ExponentialMaxWaitUpTo**: same as ExponentialUpTo, but the maximum wait time is capped so that if your wait times grow too large, you can set a bound on the wait time's growth. This is very important for exponential because the wait times can grow very quickly

# Examples

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
	"github.com/wojnosystems/go-retry/retryError"
	"time"
)

func main() {
	// you can define your retry strategy and use it wherever you need
	// you can use the struct or retry.NewUpTo constructor, whichever you find cleaner
	retryStrategy := &retry.UpTo{
		// we will wait 10ms between every attempt that failed but could be retried
		WaitBetweenAttempts: 10 * time.Millisecond,
		// we will only attempt the callback 10 times before returning the last retryable error
		MaxAttempts: 10,
	}

	tries := 0
	timer := common.NewTimeSet()
	var err error
	duration := common.TimeThis(func() {
		err = retryStrategy.Retry(context.TODO(), func() (err error) {
			fmt.Println(timer.SinceLast())
			tries++
			return retryError.Again(errors.New("simulated error"))
		})
	})
	fmt.Println("tried", tries, "times taking", duration)
	fmt.Println("should get 'simulated error': ", err.Error())
}
```

### Outputs

```
487ns
10.248678ms
10.268813ms
10.234624ms
10.232686ms
10.213887ms
10.165354ms
10.206152ms
10.205407ms
10.186255ms
tried 10 times taking 91.997547ms
should get 'simulated error':  simulated error
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
	"github.com/wojnosystems/go-retry/retryError"
	"time"
)

func main() {
	tries := 0
	duration := common.TimeThis(func() {
		timer := common.NewTimeSet()

		_ = retry.NewExponentialMaxWaitUpTo(
			10*time.Millisecond,
			1.5,
			10,
			100*time.Millisecond,
		).Retry(context.TODO(), func() (err error) {
			tries++
			fmt.Println(timer.SinceLast())
			return retryError.Again(errors.New("simulated error"))
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
	"github.com/wojnosystems/go-retry/retryError"
	"time"
)

func main() {
	tries := 0
	duration := common.TimeThis(func() {
		timer := common.NewTimeSet()

		_ = retry.NewForever(
			10*time.Millisecond,
		).Retry(context.TODO(), func() (err error) {
			fmt.Println(timer.SinceLast())
			if tries < 10 {
				tries++
				return retryError.Again(errors.New("simulated error"))
			}
			return nil
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
	"github.com/wojnosystems/go-retry/retryError"
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
		return retryError.Again(errors.New("some error"))
	})

	strategy = retry.Never

	_ = strategy.Retry(context.TODO(), func() (err error) {
		fmt.Println("NEVER")
		return retryError.Again(errors.New("some error"))
	})

	strategy = normal
	_ = strategy.Retry(context.TODO(), func() (err error) {
		fmt.Println("normal")
		return retryError.Again(errors.New("some error"))
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

## Network I/O: Exponential Back off with a maximum wait cliff

Here's an example of how you can use this to connect with an HTTP request and retry on failure. This example shows you how to use an exponential back-off and what happens when your service fails permanently.

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/wojnosystems/go-retry/examples/common"
	"github.com/wojnosystems/go-retry/retry"
	"github.com/wojnosystems/go-retry/retryError"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	timer := common.NewTimeSet()

	dialerStrategy := &retry.ExponentialMaxWaitUpTo{
		InitialWaitBetweenAttempts: 50 * time.Millisecond,
		GrowthFactor:               1.0,
		MaxAttempts:                15,
		MaxWaitBetweenAttempts:     500 * time.Millisecond,
	}

	totalTime := common.TimeThis(func() {
		err := dialerStrategy.Retry(ctx, func() error {
			fmt.Println("getting", timer.SinceLast())
			req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/non-existent", nil)
			req = req.WithContext(ctx)
			_, getErr := http.DefaultClient.Do(req)
			if getErr != nil {
				getErr = errors.Unwrap(getErr)
				if getErr != context.DeadlineExceeded {
					getErr = retryError.Again(getErr)
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

# FAQ

## I want to write my own wait back-off strategy

Awesome. This library is intended to be extensible. You should be able to use the retryLoop.* methods to build your own strategies. You can basically just pass in a method that indicates how long each round should take. Provide your own method to calculate this with any state you need. I can easily see building a smarter exponential back-off using retryLoop as a base.

## Why retry.Skip and retry.Never?

I found that I was making these where I used this retry library to test. These are very nice mocks you can use in your tests from this library.

# Disclaimer

Use this library at your own risk. Christopher Wojno and any other authors are not liable for any damages. No warranty is provided or implied.
