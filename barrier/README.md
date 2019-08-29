# Barrier Concurrency Pattern
The Barrier pattern opens the door of microservices programming with its composable nature. 
It could be considered a Structural pattern, as you can imagine. The Barrier pattern is not 
only useful to make network requests; we could also use it to split some task into multiple 
Goroutines. For example, an expensive operation could be split into a few smaller operations 
distributed in different Goroutines to maximize parallelism and achieve better performance.

Imagine the situation where we have a microservices application where one service needs to 
compose its response by merging the responses of another three microservices. This is where 
the Barrier pattern can help us.

Our Barrier pattern could be a service that will block its response until it has been composed 
with the results returned by one or more different Goroutines (or services). And what kind of 
primitive do we have that has a blocking nature? Well, we can use a lock, but it's more idiomatic 
in Go to use an unbuffered channel.


## Objectives
As its name implies, the Barrier pattern tries to stop an execution so it doesnt finish before it's ready 
to finish. The Barrier pattern's objectives are as follows:

* Compose the value of a type with the data coming from one or more Goroutines.
* Control the correctness of any of those incoming data pipes so that no inconsistent data is returned. 
We don't want a partially filled result because one of the pipes has returned an error.

# Example
For our example, we are going to write a very typical situation in a microservices 
application; an app that performs two HTTP GET calls and joins them in a single response 
that will be printed on the console.

Our small app must perform each request in a different Goroutine and print the result on 
the console if both responses are correct. If any of them returns an error, then we print 
just the error.

The design must be concurrent, allowing us to take advantage of our multi-core CPUs 
to make the calls in parallel:

//todo: Add a diagram

In the preceding diagram, the solid lines represent calls and the dashed lines represent channels. 
The balloons are Goroutines, so we have two Goroutines launched by the main function (which could 
also be considered a Goroutine). These two functions will communicate back to the main function by 
using a common channel that they received when they were created on the makeRequest calls.


## Acceptance criteria
Our main objective in this app is to get a merged response of two different calls, so we can describe 
our acceptance criteria like this:

* Print on the console the merged result of the two calls to an endpoint to get encrypted private key, 
and then Google Cloud Key Management Service to decrypt the private key.
* If any of the calls fails, it must not print any result but the error message (or error messages 
if both calls have failed).
* The output must be returned as array of strings when both calls have finished. It means that we 
cannot return the result of one call and then the other.


## Tests
There are two sets of tests for this pattern:
* Integration tests are included under `barrier_test` package to emulate import and creating real HTTP GET Calls
* Unit tests are included under `barrier` package as per classic Golang convention

In order to run the tests individually, please go inside `barrier` or `barrier_test` directory and run the following:

    ~$: go test -run=<NameOfTestMethod> -v .
    ~$: go test -run=TestBarrier -v .

if you wish to run the tests individually, please use following command:
    
    ~$: go test -run=TestBarrier/<test_description> -v .
    ~$: go test -run=TestBarrier/with_correct_endpoint -v .
