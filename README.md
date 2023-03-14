## How to build
run `go run main.go` command inside root directory. Server starts to listen on localhost:8080

port can be specified through flag --port=8080

request can be sent like this:
`
curl --request POST 'http://localhost:8080/' \
--header 'Content-Type: application/json' \
--data '[1, 2, 3, 4, 5]'
`
server responds with Content-Type: application/json

tests are located inside internal/pkg/handlers directory
to run tests run `go test ./internal/pkg/handlers -v` inside root directory

## Decision explanation
I have designed a standard worker pool with limited number of workers to run simultaneously calculating the result. It is defined by constant MAX_WORKERS in handlers/handler.go 
Advantage of such solution is that response is calculated in parallel, assuming that there may be big numbers and large input size. The disadvantage is that it can consume a lot of resources. For that situation, only for large inputs result is calculated in parallel, otherwise it is done sequentially

There are a lot of ways to improve the perfomance of this solution, for example we could record the calculated values in cache or persistent storage so that out application does not estimate similar values over and over again. Another way is to launch a background task to precalculate those values, so that for some input we could extract answer in O(1) time (right now it is O(n*m) where n - length of the input and m is the maximum number in array) 
There is also a check for the size of the input so that we dont accept large message that does not fit into our memory


## Assumptions
I assumed client sends POST requests with Content-Type: application/json
Also I assumed that client won't close the connection until he gets the response, so context cancellation check was not implemented
