# Recent Tweet

A simple way to show your latest tweet. Every 24hrs it fetches your latest tweet and stores it in DynamoDB. A simple GET endpoint fetches your tweet. Now you can display your latest tweet where you need it!

## Requirements

- AWS CLI already configured with Administrator permission
- [AWS SAM CLI](https://github.com/awslabs/aws-sam-cli)
- [Docker installed for local development](https://www.docker.com/community-edition)
- [Golang](https://golang.org)

## Setup process

Copy and paste the `template.example.yaml` and turn it into `template.yaml`. You will need your access and consumer tokens from a twitter app for you to fetch your latest tweet.
