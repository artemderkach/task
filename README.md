## Description
Service retrieves random name and random joke, after combining both, you will receive the result

## Prerequirements
Go compiler

## Implementation reasoning
Giving that task if fairly simple, project structure will be minimal.

## running service
`$ go run main.go` - to run  
`$ go run main.go -p 7777` - to run on specific port  
`$ go run main.go --help` - to call for help  

## using service
Send request to retrieve random joke with random name `GET /`  
`GET /ping` - to test if service is running.

## testing
`$ go test .` - to run tests
