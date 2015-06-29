# go-rest-api
Playing around with go

## To build & run
```
go get github.com/emicklei/go-restful
go get github.com/stretchr/testify
go test
go build
```
You should now have a nice binary with the same name as the folder you checked out the project into (go-rest-api by default).

Run it like any other executable:
```
./go-rest-api
```
(It needs to be run in the same directory as the swagger resources, didn't feel like fiddling around with bindata package)
