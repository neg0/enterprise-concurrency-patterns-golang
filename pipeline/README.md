# Pipeline Concurrency Pattern
Every time that we write any function that performs some logic, we are writing 
a pipeline: If this then that, or else something else. Pipelines pattern can be 
made more complex by using a few functions that call to each other. They can even 
get looped in their out execution.

The Pipeline pattern in Go works in a similar fashion, but each step in the Pipeline 
will be in a different Goroutine and communication, and synchronizing will be done 
using channels.


## Objectives
When creating a Pipeline, we are mainly looking for the following benefits:

* We can create a concurrent structure of a Multi-step algorithm
* We can exploit the parallelism of Multi-core machines by decomposing an algorithm in 
different Goroutines

However, just because we decompose an algorithm in different Goroutines doesnt necessarily 
mean that it will execute the fastest. We are constantly talking about CPUs, so ideally the 
algorithm must be CPU-intensive to take advantage of a concurrent structure. The overhead 
of creating Goroutines and channels could make an algorithm smaller.


# Example
In this example we are going to create real Machine Learning application that enriches the 
bank transactions with Category and Merchant name. Transactions are modeled after [Open Banking 
specification v3.1](https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937558098/Transactions+v3.1) 
and endpoints for enrichment of category and merchant are randomly selected for purpose of training. 
Transactions should be enriched according to field `TransactionInformation` that specifies the brief 
information about transaction.


## Acceptance Criteria
blah blah


## Tests
There are two sets of tests for this pattern:
* Integration tests are included under `future_test` package to emulate import and creating real HTTP GET Calls
* Unit tests are included under `future` package as per classic Golang convention

In order to run the tests individually, please go inside `future` or `future_test` directory and run the following:

    ~$: go test -run=<NameOfTestMethod> -v .
    ~$: go test -run=TestPipeline -v .

if you wish to run the tests individually, please use following command:
    
    ~$: go test -run=TestPipeline/<test_description> -v .
    ~$: go test -run=TestPipeline/should_finish_enrichment_process_without_an_error -v .
