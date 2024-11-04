## Getting started

1. Install dependencies: golang 1.17+, golangci-lint, docker + docker-compose.
2. Default parameters are set in docker-compose.yml and can be inlined in command line when starting the application.
3. Makefile commands are available for convenience to build and run the project (see below).

### Makefile commands

Start server locally without docker:

`SERVER_HOST=:13371 CHALLENGE_COMPLEXITY=12 SOLUTION_TIMEOUT=15s READ_TIMEOUT=30s make run-server`

Start client locally without docker:

`READ_TIMEOUT=30s SERVER_HOST=:13371 make run-client`

Build and run docker images for server and client:

`make build-and-run-docker`

Run previously built docker images:

`make run-docker`

Run tests:

`make test`

Run linters:

`make lint`

### Reasoning behind the choice of Hashcash 
1. Easy to implement and to verify the code.
2. Solution is efficiently verifiable on the server-side.
3. Adjustable complexity of challenges.
4. Hashcash is used as a proof-of-work algorithm for various cryptocurrencies, e.g. Bitcoin.

### Possible future roadmap (depending on the product requirements)
1. Implement a RAM-bound algorithm (MBound) as it provides a smaller time difference to solve between new and old client hardware.
2. Allow switching between algorithms in settings.
3. Generate challenges for both algorithms to allow clients to choose which one to solve (possibly solving both concurrently).

Sources:

https://en.wikipedia.org/wiki/Hashcash

http://www.hashcash.org/hashcash.pdf

http://www.hashcash.org/papers/memory-bound-crypto.pdf

