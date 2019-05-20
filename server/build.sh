#!/bin/bash

GOOS=linux go build *.go
zip deployment.zip main
mv deployment.zip ../terraform/publish/
rm main