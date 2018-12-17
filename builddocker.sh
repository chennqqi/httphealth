#!/bin/bash

go build -v -ldflags "-s -w" 
sudo docker build -t 'sort/httphealth' -f Dockerfile  .
