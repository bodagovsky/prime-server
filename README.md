## How to build
run `go run main.go` command inside root directory. Server starts to listen on localhost:80
request can be sent like this:
`
curl --location --request POST 'http://localhost:80/' \
--header 'Content-Type: text/plain' \
--data-raw '[1, 2, 3, 4, 5]'
`
server responds with Content-Type: text/plain 

tests are located inside ./tests directory
to run tests run `go test ./tests -v` inside roor directory

## Decision explanation
I have designed a standard worker pool with limited number of workers to run simultaneously calculating the result. It is defined by constant MAX_WORKERS in handlers/handler.go 
Advantage of such solution is that response is calculated in parallel, assuming that there may be big numbers and large input size. The disadvantage is that it can consume a lot of resources.

Although I added graceful shutdown I have not implemented context cancellation tracking for the sake of simplicity (which totally necessary in production). Drawbacks of such implementation may be goroutines that continue working after shutting down the connection by client.
There are a lot of ways to improve the perfomance of this solution, for example we could record the calculated values in cache or persistent storage so that out application does not estimate similar values over and over again. Another way is to launch a background task to precalculate those values, so that for some input we could extract answer in O(1) time (right now it is O(n*m) where n - length of the input and m is the maximum number in array) 
There is also a check for the size of the input so that we dont accept large message that does not fit into our memory


## Assumptions
I assumed client sends POST requests with Content-Type: plain/text and there is an arrangement for the syntax which allows only square brackets for denoting an array.
There is also an integer overflow check inside string conversion, which may be not so obviuos for the client. There are also different solutions to handle large numbers, in my implementation I decided to inform the user if he sent number that is too large
