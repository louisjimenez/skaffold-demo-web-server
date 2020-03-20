# Web Server

An example Web Server microservice for evaluating Skaffold. 

## Development

### Run the Skaffold dev loop
`skaffold dev --default-repo <path-to-registry> --port-forward`

### Compile Protos
`protoc -I todo/ todo/todo.proto --go_out=plugins=grpc:todo`