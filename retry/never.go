package retry

// Never will not retry. It will try the callback once, and then stop, whether it succeeds or fails, retryAgain or not
var Never = &UpTo{
	WaitBetweenAttempts: 0,
	MaxAttempts:         0,
}
