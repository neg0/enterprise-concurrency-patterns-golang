# Future (Promise) Concurrency Pattern
The Future design pattern (also called Promise) is a quick and easy way to achieve concurrent 
structures for asynchronous programming. We will take advantage of first class functions in Go 
to develop Futures.

we will define each possible behavior of an action before executing them in different Goroutines. 
Node.js uses this approach, providing event-driven programming by default. The idea here is to 
achieve a fire-and-forget that handles all possible results in an action.

<p align="center">
  <img src="https://user-images.githubusercontent.com/5575209/64031089-a4d03f00-cb3f-11e9-9bd1-2c3941d73a96.png">
</p>

In the preceding diagram, the main function launches a Future within a new Goroutine. It won't 
wait for anything, nor will it receive any progress of the Future. It really fires and forgets it.

Like Promises in JavaScript; we can launch a new Future within a Future and embed as many 
Futures as we want in the same Goroutine (or new ones). The idea is to take advantage of the result 
of one Future to launch the next.


## Objectives
With the Future pattern, we can launch many new Goroutines, each with an action and its own handlers. 
This enables us to do the following:

* Delegate the action handler to a different Goroutine
* Stack many asynchronous calls between them (an asynchronous call that calls another asynchronous 
call in its results)

This is a kind of lazy programming, where a Future could be calling to itself indefinitely or just until 
some rule is satisfied. The idea is to define the behavior in advance and let the future resolve the 
possible solutions.


# Example
We are going to develop a very simple example to try to understand how a Future works. In this example, 
we will have a method that returns a string or an error, but we want to execute it concurrently. 
We have learned ways to do this already. Using a channel, we can launch a new Goroutine and handle the 
incoming result from the channel.

But in this case, we will have to handle the result (string or error), and we don't want this. Instead, 
we will define what to do in case of success and what to do in case of error and fire-and-forget the Goroutine.

We could also have done asynchronous programming without functions by setting an interface with Success, Fail, 
and Execute methods and the types that satisfy them, and using the Template pattern to execute them asynchronously


## Acceptance criteria
We don't have functional requirements for this task. Instead, we will have technical requirements for it:

* Delegate the function execution to a different Goroutine
* The function will return a string (maybe) or an error
* The handlers must be already defined before executing the function
* The design must be reusable


## Tests
There are two sets of tests for this pattern:
* Integration tests are included under `future_test` package to emulate import and creating real HTTP GET Calls
* Unit tests are included under `future` package as per classic Golang convention

In order to run the tests individually, please go inside `future` or `future_test` directory and run the following:

    ~$: go test -run=<NameOfTestMethod> -v .
    ~$: go test -run=TestFuture -v .

if you wish to run the tests individually, please use following command:
    
    ~$: go test -run=TestFuture/<test_description> -v .
    ~$: go test -run=TestFuture/when_result_is_successful -v .
