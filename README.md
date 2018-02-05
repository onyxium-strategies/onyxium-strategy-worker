### Strategy workflow
An end user specifies the conditions and actions of his strategy on the website. When the end user is done this strategy is send to the collector which collects all jobs. Then the collector send it to the dispatcher which dispatches the job to a worker. The worker keeps running the job untill it is finished.

The workflow will be implemented in (Go)[https://golang.org]. More resources can be found here: (Go tour)[https://tour.golang.org/welcome/1], (Worker queue example)[http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html], (Chat server with websockets example)[https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets], (Go with websockets example)[https://jacobmartins.com/2016/03/07/practical-golang-using-websockets/], (Why Go?)[https://medium.com/@kevalpatel2106/why-should-you-learn-go-f607681fad65], (Go concurrency dev talk)[https://www.youtube.com/watch?v=f6kdp27TYZs]

### GO Lang installation
Visit the get started page by GO: https://golang.org/doc/install
To build an applicate use the command `go build` and then run it with `./application-name`.


