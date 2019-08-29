<h1 align="center">Enterprise Concurrency Patterns in Golang</h1>
<p align="center"><img src="https://travis-ci.org/neg0/docfony.svg?branch=master" alt="build:passed"></p>
<p>Concurrency patterns mostly manage the timing execution and order execution of applications that has more 
than one flow. Examples in this repository are done in TDD, integration tests are included under appended `_test` package.</p>


> This is a work in progress, which will be updated weekly bases till end of September 2019 (compilation date).

* __Barrier:__ is a common pattern, especially when we have to wait for more than one response 
from different Goroutines before letting the program continue <small>_(Done)_</small>

* __Future:__ pattern allows us to write an algorithm that will be executed eventually in time 
(or not) by the same Goroutine or a different one <small>_(In Progress)_</small>

* __Pipeline:__ is a powerful pattern to build complex synchronous flows of Goroutines that are 
connected with each other according to some logic <small>_(In Progress)_</small>

* __Publisher/Subscriber:__ <small>_(Pending)_</small>

* __Workers Pool:__ <small>_(Pending)_</small>

 
## Docker & Makefile
If you have a docker you could use `docker-compose` file to run your Go environment using commands 
provided in `Makefile`. For more information about available commands please run following at the 
root of this repository.

    ~$ make help

## Credits
Built with :heart: & :coffee: at the heart of beautiful London


## License
Reuse of examples in this repository is prohibited for commercial use and OpenSource only by citation and my written approval.