#!/bin/bash

name=goipmi

[ -f $name ] && rm -f $name
go build -o $name
