# Overview

The customToolkit provides a few utilities for developers who want to build their own retry logic from a common logic.

The most complicated parts are LoopForever and LoopUpTo. These should handle the most common retry patterns. This is also the easiest place to make a mistake with retry logic, and this handles all of that garbage for you.

See "retry" package for ample examples of how to use these basic building blocks to build your own.
