#!/bin/bash
go get "github.com/gin-gonic/gin"
go get --insecure "scc-gitlab-1.dev.octanner.net/octanner/octvault"
go build main.go

