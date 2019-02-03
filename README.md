# AWS Lambda Workshop
Go Days Berlin workshop for building AWS Lambda using Golang and Serverless platform.

### Description

TBD 

### Requirements

* Go
* [Dep](https://github.com/golang/dep) 
* AWS account, [AWS CLI](https://aws.amazon.com/cli/), security credentials
* NodeJS
* [Serverless platform](https://github.com/serverless/serverless)


## Installation 

`dep ensure` - installs required packages

`serverless` - to view list of commands (btw endpoints can be invoked locally)

`make build` - builds binaries (check **Makefile**)

`make deploy` - deploy binaries to S3 and then create Lambda functions.

## Examples 

* **world** - lambda execution with response and message;

* **hello** - with query parameter and response as JSON;

* **short-url** - single endpoint for POST URL aliases to Dynamo DB, and GET full URL by alias.