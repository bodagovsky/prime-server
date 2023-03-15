## How to build
run `go run cmd/main.go` command inside root directory. Server starts to listen on localhost:8080

port can be specified through flag --port=8080

request can be sent like this:
`
curl --request POST 'http://localhost:8080/' \
--header 'Content-Type: application/json' \
--data '[1, 2, 3, 4, 5]'
`
server responds with Content-Type: application/json

tests are located inside internal/pkg/handlers directory
to run tests run `go test ./internal/handlers -v` inside root directory

## Decision explanation
I have designed a standard worker pool with limited number of workers to run simultaneously calculating the result. It is defined by constant MAX_WORKERS in handlers/handler.go 
Advantage of such solution is that response is calculated in parallel, assuming that there may be big numbers and large input size. The disadvantage is that it can consume a lot of resources. For that situation, only for large inputs result is calculated in parallel, otherwise it is done sequentially

There are a lot of ways to improve the perfomance of this solution, for example we could record the calculated values in cache or persistent storage so that out application does not estimate similar values over and over again. Another way is to launch a background task to precalculate those values, so that for some input we could extract answer in O(1) time (right now it is O(n*m) where n - length of the input and m is the maximum number in array) 
There is also a check for the size of the input so that we dont accept large message that does not fit into our memory

I also decided not to implement hexagonal architecture as it would appear in production environment because it is not clear how the project would grow, so I keep all the necessary logic close to handler.
Directory /pkg containes only one subdir called /math because it seemed to me that it is the only logic that can be safely being imported by other projects and reused. The reason why math is another directory is that there may appear other dirs that will reside inside pkg directory for reuse.


## Assumptions
I assumed client sends POST requests with Content-Type: application/json
Also I assumed that client won't close the connection until he gets the response, so context cancellation check was not implemented

Another assumption is that when there more handlers will be implemented it is good practice to tie them up to a single struct, for example, called service which holds some common parameters for all handlers like logger. It will look something like that:

`type service struct {
    log logger
} 

func New() *service{
    return &service{}
}

func (s *service) PrimeHandler() {}`
