# github-star-client

go run client.go  -r golang/go -r tinygo-org/tinygo-site -r golang/g

Naming covention of organization and repository in github:
number, character, dash(-), underscore(_), dot(.), 
comma(,) and spaces Any special character will be convert to dash(-). For example: @me will be -me


go test ./... -coverprofile cover.out

go tool cover -html=cover.out