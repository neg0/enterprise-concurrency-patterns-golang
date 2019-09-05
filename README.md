<h1 align="center">Trinity Concurrency Patterns in Golang</h1>
<p align="center"><img src="https://travis-ci.org/neg0/enterprise-concurrency-patterns-golang.svg?branch=master"> <a href="https://coveralls.io/github/neg0/enterprise-concurrency-patterns-golang?branch=master" target="_blank"><img src="https://coveralls.io/repos/github/neg0/enterprise-concurrency-patterns-golang/badge.svg?branch=master"></a></p>

<p>Concurrency patterns mostly manage the timing execution and order execution of applications that has more 
than one flow. Examples in this repository are done in TDD, integration tests are included under appended `_test` package.</p>

* __Barrier:__ is a common pattern, especially when we have to wait for more than one response 
from different Goroutines before letting the program continue

* __Future:__ pattern allows us to write an algorithm that will be executed eventually in time 
(or not) by the same Goroutine or a different one

* __Pipeline:__ is a powerful pattern to build complex synchronous flows of Goroutines that are 
connected with each other according to some logic

 
## Docker & Makefile
If you have a docker you could use `docker-compose` file to run your Go environment using commands 
provided in `Makefile`. For more information about available commands please run following at the 
root of this repository.

    ~$ make help


## Documentation
Each pattern contains a single documentation explaining the pattern and the example I've created for 
each pattern. Please ensure you run the integration while having test container up and running using 
commands in make file for Docker Compose file inside this project.


## Credits
Built with :heart: & :coffee: at the heart of beautiful London


## License
Reuse of examples in this repository is prohibited for commercial use and OpenSource only by citation and my written approval.