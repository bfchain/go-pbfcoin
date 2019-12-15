#!/bin/bash -x

# creates the necessary docker images to run testrunner.sh locally

docker build --tag="pbfcoin/cppjit-testrunner" docker-cppjit
docker build --tag="pbfcoin/python-testrunner" docker-python
docker build --tag="pbfcoin/go-testrunner" docker-go
