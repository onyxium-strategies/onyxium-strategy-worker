# Onyxium Strategy Worker
An end user specifies the conditions and actions of his strategy in the frontend (see [onyxium-webapp-client](https://github.com/onyxium-strategies/onyxium-webapp-client). When the end user is done this strategy is send to the collector which collects all jobs. All jobs are parsed to a K-ary tree structure. Then the collector send it to the dispatcher which dispatches the job to a worker. The worker keeps running the job untill it is finished.

# Getting started
Install go: https://golang.org/dl/
Make sure your $GOPATH is correct. In order to clone a private bitbucket repo we need to enable default cloning with SSH instead of https.
`go get https://github.com/onyxium-strategies/o yxiumstrategy-worker.git`
`cd $GOPATH/bin`
`./onyxium-strategy-worker`

# Development
Dependancy management during development is done with [golang/dep](https://golang.github.io/dep/docs/introduction.html).
Note that this is not intended for end users who are installing Go software - that's what `go get` does.

Install with homebrew: `brew install dep`
Pull all dependencies: `dep ensure`
Add new dependency to project: `dep ensure -add github.com/foo/bar`

# Resources
* [Go tour](https://tour.golang.org/welcome/1)
* [Worker queue example](http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html)
* [Chat server with websockets example](https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets)
* [Go with websockets example](https://jacobmartins.com/2016/03/07/practical-golang-using-websockets/)
* [Why Go?](https://medium.com/@kevalpatel2106/why-should-you-learn-go-f607681fad65)
* [Go concurrency dev talk](https://www.youtube.com/watch?v=f6kdp27TYZs)
* [mgo query objectid for range of time values](https://stackoverflow.com/questions/31502195/mgo-query-objectid-for-range-of-time-values)
* [Dependency Injection in Golang](https://medium.com/@zach_4342/dependency-injection-in-golang-e587c69478a8)
* [Golang Logger](https://github.com/Sirupsen/logrus)
* [Test Driven Development](https://leanpub.com/golang-tdd/read#leanpub-auto-test-driven-development)
* [TDD](https://www.binpress.com/tutorial/getting-started-with-go-and-test-driven-development/160)
* [Organising DB access](http://www.alexedwards.net/blog/organising-database-access)
* [K-ary tree data structure](https://en.m.wikipedia.org/wiki/Left-child_right-sibling_binary_tree)
