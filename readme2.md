# Introduction

This is the official golang client for **cella**.

This client can be used in your golang projects which want to interact with **cella server**.

An example can be found under `example/cli/main.go` where we implemented a little cella CLI using the client library.

View the [changelog](CHANGELOG.md) for the latest updates and changes by version.

# Getting Started

## Installation

To start using this client run `go get`:
``` 
go get github.com/check24/cella-client-go
```
This will retrieve the library
## Usage

The client library offers different clients, which use different protocols to transfer the data to the server.
If you want to use e.g. UDP and TCP at the same time, two clients (one for UDP and one for TCP) have to be created.

### UDP Client
The UDP Client sends data to the **cella server** using UDP. Data is sent directly without any queuing.
``` go
// create the client
client, err := cella.NewUDPClient("127.0.0.1:1401")
// send message(s)
err = client.Put("teststream", []byte("payload"))
// close connection after sending all messages
err = client.Close()
```

### TCP Client
The TCP Client sends data to the **cella server** using TCP. For performance reasons data is queued.
``` go
// create the client
client, err := cella.NewTCPClient("127.0.0.1:1402")	
// send message(s)
err = client.Put("teststream", []byte("payload"))
// close connection after sending all messages
err = client.Close()
```
You can control the queue by passing **ClientOptions** to `NewTCPClient()`. The **ClientOptions** contain 3
values. Only set these values if you want to override the default values.
- `queueSize` -  Size of queue (golang channel), if queue is full all calls to `Put()` will block.
- `flushInterval`- This variable controls the flushing interval
- `bulkSize`- If the queue holds bulkSize messages the queue will be flushed

Messages will be sent whatever comes first: reaching `flushInterval` or `bulkSize`.

``` go
client, err := cella.NewTCPClient(*address,
    cella.NewClientOptions().
      SetQueueSize(1000).
      SetFlushInterval(time.Second).
      SetBulkSize(20),
)
```

### Debug Client
The debug client is for local debugging only as the name suggests. It prints all data to `STDOUT`.

``` go
client, err := cella.NewDebugClient()
```

# Development

This project uses a Makefile where all tasks are configured. `make help` will print out all commands and their function. Some tasks will not work on windows!

### Linting
- Use `make lint-install` to install golangcli
- Use `make lint` to run the linter

### Testing / Benchmarking
- Use `make unittest` to run all unittests
- Use `make systemtest` run all systemtests. In order to execute this command the cellad server has to run in the background.
- Use `make benchmark` to execute the benchmarks
### cellad
If you want to run systemtests, you need to install and run the cellad server
- Use `make cellad-install` to install the cellad server in a new directory called cellad
- Use `make cellad-run` to start the cellad server
- Use `make cellad-clean` to remove the cellad/cellaFilesink.dat file