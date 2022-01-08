# Catu

Catupiri framework core module

## usage

```go
// Start the new app:
app := catu.NewApp()
// Register your plugins
app.RegisterPlugin("mm", mm.GetPlugin())
// Start the bootstrap process, will load all resources, bind routes, middlewares ...
app.Bootstrap()

// here you can use the app resources like run a command or start a server...

// Start the http server if required:
app.StartHTTPServer()
```

## Core events

Powered by: https://github.com/gookit/event

- configuration
- bindMiddlewares
- bindRoutes
- setResponseFormats
- bootstrap

## Build with

- Go
- Time
- Unknow things that some call magic