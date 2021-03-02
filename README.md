Name: Chow Man Chun

New to golang

# Approach
All the code is residing in the `accountapi-client` folder

I used TDD to implement the test

The library code is located in `accountapi-client/client.go`

The method `NewClient` returns a new client struct, and access to CREATE, FETCH and DELETE methods

All the methods share the same parse response method to handle response of the API servers

There is a config package to read the config file and set the domain and port of the API server.