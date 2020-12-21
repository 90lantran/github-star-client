# github-star-client
This is a simple Go http client to send a POST request, and receive a response for [github-star](https://github.com/90lantran/github-star) server.
The client will take in a list of organization/reposiotry, server host `http://ip:port` (default http://localhost:8080) from command line, and print out the response from server or any internal client error. More details about input validation is [here](#2.input).

## Usage
### 1.Build client
```
$ make build
```
./client excutable should be created at root directory

### 2.Input
Command line arguments are handled by argparser from [akamensky]("https://github.com/akamensky/argparse").
This client supports 2 flags: -r for input list, -t host and port.

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
  "payload": {
    "totalStars": 79999,
    "invalidRepos": [
      "me/e",
      "teori/23423",
      "324324/43"
    ],
    "validRepos": [
      {
        "name": "golang/go",
        "star(s)": 79999
      }
    ]
  },
  "error": "At least one of the input is not valid",
  "status": "success"
}
```

If github-star server runs at localhost, you `should not` specify -t flag. This option is useful when you deploy github-star sever to minikube, you can pass in the ip and port of minikube to test it. Here is an example.

```
$ ./client -r tinygo-org/tinygo-site,golang/go -r 4534/433 -t http://192.168.99.107:30000
input list: {Input:[tinygo-org/tinygo-site golang/go 4534/433]}
Response: {
  "payload": {
    "totalStars": 80032,
    "invalidRepos": [
      "4534/433"
    ],
    "validRepos": [
      {
        "name": "tinygo-org/tinygo-site",
        "star(s)": 22
      },
      {
        "name": "golang/go",
        "star(s)": 80010
      }
    ]
  },
  "error": "At least one of the input is not valid",
  "status": "success"
}
```

### 3.Input validation
I found a fun thing about github naming convention for repository name. Valid inputs contains numbers, characters, dashes(-), underscores(_), dots(.). If you type in comma(,), whitespace, any special characters, they will be converted to dash(-). For example: @me will be -me. That is how I wrote my regular expression for input list.

## Unit test
unit-test was written with go test.

- To run unit-test: 
```
$ make unit-test
cd pkg/client && go test ./... -coverprofile cover.out
ok      github.com/90lantran/github-star-client/pkg/client      0.019s  coverage: 81.8% of statements 
```
- To show code-coverage:
```
$ make code-coverage
```