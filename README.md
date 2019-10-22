# go-http-api

A package that provides an HTTP REST API for a Filecoin implementation written in go.

[![CircleCI](https://circleci.com/gh/filecoin-project/go-http-api.svg?style=svg)](https://circleci.com/gh/filecoin-project/go-http-api)
[![codecov](https://codecov.io/gh/filecoin-project/go-http-api/branch/master/graph/badge.svg)](https://codecov.io/gh/filecoin-project/go-http-api)

## Features
* Response body is JSON.
* POST request bodies are expected to be JSON.
* SSL/TLS supported
* Bearer auth scheme supported (SOON)

## Install
```
go get github.com/filecoin-project/go-http-api
```

## Implement

### 1. Set up `vN.Callbacks`
The core of the API is the `Callbacks` struct.  It is a typestruct in each API version package containing named callback functions, which should call into your code.  The server is then instantiated with your desired callbacks. Each version of the API has its own `Callbacks` struct: 
```go
package v1
import(...)
type Callbacks struct {
	GetActorByID func(string) (*types.Actor, error)
	GetActorNonce func(string) (*big.Int, error)
	GetActors func() ([]*types.Actor, error)
    // ... etc
}
```
Because it is a struct and not an interface, implementers are free to support as much of the API as they like; a default handler will be used for nil `Callbacks` in each, for example:
```bash
curl http://localhost:5000/api/filecoin/v1/actors
    /api/filecoin/v1/actors is not implemented
``` 
A 404 response will be sent for endpoints that don't (and can't) exist:
```bash
curl http://localhost:5000/api/filecoin/v1/atcorz
    curl: (22) The requested URL returned error: 404 Not Found
```
 This standardizes unimplemented endpoint responses for every node, ensures the API endpoints are compliant with the [API spec](https://github.com/filecoin-project/filecoin-http-api), and more easily allows the API consumer to know what functionality is implemented -- or at least, what is exposed to the API -- by each node. 

In order to be implementation-agnostic, this package uses its own Filecoin-based typestructs for callbacks and serialized responses.

### 2. Instantiate and run the server
Place the following code somewhere in your node startup sequence:

```go
cb := &v1.Callbacks {
    GetActorByID: cbs.MyGetActorByIDFunc,
    GetActorNonce: cbs.MyGetActorNonceFunc,
    // ...
    SendSignedMessage: cbs.MySendSignedMessageFunc,
    WaitForMessage: cbs.MyWaitForMessageFunc,
}

cfg := server.Config{
    Port: 5001,
    TLSCertPath os.Getenv("TLS_CERT_PATH")
    TLSKeyPath  os.Getenv("TLS_KEY_PATH")
}

s := server.NewHTTPAPI(context.Background(), cb, cfg).Run()
```

### 3. Launch your node and test the endpoints
First launch your filecoin node. Then, to verify only that you have correctly launched the HTTP API server, use the `hello` endpoint:

```bash
curl http://localhost:5000/api/filecoin/v1/hello
    HELLO
```

Then attempt to retrieve information from your node, by sending a request to an endpoint for one of the callbacks you implemented:
```bash
curl http://localhost:5000/api/filecoin/v1/actors
```

Assuming you have correctly implemented your callbacks, you should see familiar output.

Please see test files for more details. 

# References
* For a current map of endpoints to callbacks, see [`server.Route()`](https://github.com/filecoin-project/go-http-api/blob/516f52ea8f6e13c708613c42e087c346c1f39e2b/server.go#L84)
* [Filecoin HTTP API Specification](https://github.com/filecoin-project/filecoin-http-api)
* filecoin-project/go-filecoin [example implementation in a branch](https://github.com/filecoin-project/go-filecoin/tree/feat/rest-api-part1)