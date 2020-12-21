# github-star-client
This is a simple Go http client to work with github-star server.
The client will take in a list of organization/reposiotry, server host http://ip:port (default http://localhost:8080) from command line, and print out the response from server or any internal client error. More details about input validation is [here](#2.input).

## Usage
### 1.Build client
```
$ make build
```
client should be created at root directory

### 2.Input
Command line arguments are handled by argparser from [akamensky]("https://github.com/akamensky/argparse").
This client supports 2: -r for input list, -t host and port.

```
$ ./client -h 
usage: client [-h|--help] -r|--request "<value>" [-r|--request "<value>" ...]
              [-t|--host "<value>"]

              sends a request and receives a response from github-star server

Arguments:

  -h  --help     Print help information
  -r  --request  List of organization/repossitory to send to server.
  -t  --host     IP address and port of the server. Default:
                 http://localhost:8080
```

Example: client takes in multiple lists
```
$ ./client -r me/e,teori/23423 -r 324324/43 -r golang/go
input list: {Input:[me/e teori/23423 324324/43 golang/go]}
Response: {
  "totalStars": 79997,
  "invalidRepos": [
    "me/e",
    "teori/23423",
    "324324/43"
  ],
  "validRepos": [
    {
      "name": "golang/go",
      "star(s)": 79997
    }
  ],
  "status": "success"
}
```

If github-stars server runs at localhost, you should not specify -t flag. This option is useful when you deploy github-stars sever to minikube, you can pass in the ip and port of minikube to test it.

### 3.Input validation
I found a fun thing about github naming convention for repository name. Valid inputs contains number, character, dash(-), underscore(_), dot(.). If you type in comma(,), spaces, any special characters, they will be converted to dash(-). For example: @me will be -me. That is how I wrote my regular expression for input list.

## Unit test
unit-test were written with go test.

- To run unit-test: 
```
$ make unit-test 
```
- To show code-coverage:
```
$ make code-coverage
```