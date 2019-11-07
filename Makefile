.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./get-tweet/get-tweet
	
build:
	GOOS=linux GOARCH=amd64 go build -o get-tweet/get-tweet ./get-tweet

debug:
	GOARCH=amd64 GOOS=linux go build -gcflags='-N -l' -o get-tweet/get-tweet ./get-tweet

test:
	go test -v ./...
	
package:
	sam package --s3-bucket sam-cli-bucket-fallenstedt --template-file template.yaml --output-template-file packaged-template.yaml

deploy:
	sam deploy --template-file packaged-template.yaml --stack-name latest-tweet  --capabilities CAPABILITY_IAM