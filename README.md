## APIGateway

APIGateway is a MusicStreaming entrypoint for API request to service.
Uses GRPC and REST server to listen requests

### Dependencies

This service depend on so-called [Platform](https://github.com/CuriosityMusicStreaming/Platform).
It provides local environment and necessary devtools(like [apisynchronizer](https://github.com/UsingCoding/ApiSynchronizer) to sync api files between services)

This service initiate connection to other services of Platform and passes request there.

Other libraries:
* [ComponentsPool](https://github.com/CuriosityMusicStreaming/ComponentsPool) - common library with components
* [ApiStore](https://github.com/CuriosityMusicStreaming/ApiStore) - repository that provides services api that synced by apisynchronizer
* [Protobuf](https://github.com/protocolbuffers/protobuf) - provides protobuf api codegen
* [GrpcGateway](https://github.com/grpc-ecosystem/grpc-gateway) - v1 only - provides rest proxy to grpc server
* Other code dependencies in `go.mod`

### Build

**To have ability to build service download [Platform](https://github.com/CuriosityMusicStreaming/Platform) and make installation steps**

Run make

```shell
make
```

Command build all dependencies and put binary file to `bin/`

Run `make publish` to dockerize service

### Test

This service has only linter because there no sense to write unit or integration tests 

You can run linter
```shell
make check
```