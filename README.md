# Why this?

Just a little test to do something with Go that I normally do with python.  
The usual stuff. Some string parsing, get / post from api and such.

# Using

Nothin really to use :)

Download the latest "release" and start the thingie..

## linux

`./cli-go-test-linux-amd64 --help`

## windows

`cli-go-test-windows-amd64.exe --help`

## the other one

`cli-go-test-darwin-amd64 --help`

# TODO

## Create random stuff to csv

`generate --records 1000 --filename output.csv`

- [x] macaddresses 
- [x] hostnames
- [x] IPv4 
- [x] groups


## Import csv

`import --filename output.csv`

- [x] read csv
- [x] convert csv to struct
- [ ] post to ise
 - [ ] macaddresses
 - [ ] hostname
 - [ ] description
 - [ ] group

