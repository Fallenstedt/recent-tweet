.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./get-tweet/get-tweet
	
build:
	GOOS=linux GOARCH=amd64 go build -o get-tweet/get-tweet ./get-tweet

package:
	sam package --s3-bucket sam-cli-bucket-fallenstedt --template-file template.yaml --output-template-file packaged-template.yaml

deploy:
	sam deploy --template-file packaged-template.yaml --stack-name hello-sam  --capabilities CAPABILITY_IAM