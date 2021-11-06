# Overview

The retryCore provides a few utilities for developers who want to build their own retry logic from common logic.

The most complicated parts is LoopUntil. This should handle the most common retry patterns. This is also the easiest place to make a mistake with retry logic, and this handles all of that garbage for you.

# LoopUntil

This is the core of this library. It performs callback execution, waits, checks for retry conditions, returns errors, etc. I suspect that all retryable logic can be implemented with this function. It separates the following concerns:

* context
* callback
* waiting mechanism
* should retry logic

The only state this method keeps is the minimum number of times the callback has been called.

# Examples

See "retry" package for ample examples of how to use these basic building blocks to build your own.
