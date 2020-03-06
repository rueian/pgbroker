# pgbroker

pgbroker is a golang library for building PostgreSQL proxy, which makes it easy to support from simple query logging to complex dynamic database mappings from an external resource controller and modification on data transferred between client and pg in streaming or per message manner.

## Usage

### static proxy

One example usage of pgbroker is just simple mapping multiple pg instances into one entry for centralizing management.

Please checkout the https://github.com/rueian/pgbroker-static project.

It is a production ready postgres proxy configured by a yaml file, and it is also has a pre-build docker image.

```shell
docker pull rueian/pgbroker-static:latest
```

### dynamic proxy

By implementing the `PGResolver` interface, the proxy is able to acquire different connections based on the client's `StartupMessage`.

```golang
type PGResolver interface {
	GetPGConn(ctx context.Context, clientAddr net.Addr, parameters map[string]string) (net.Conn, error)
}
```

Please check out the https://github.com/rueian/godemand-example project, which uses godemand as an external http resource controller for dynamic pg mapping.

In this way, the external resource controller is able to do any change to the postgres before connection estiblished, including creating new postgres instance on the fly.

### data modification

pgbroker provides `MessageHandler` and `StreamCallback` for each postgres protocol v3 messages. Developers can use them to easily modify the data transferred between client and pg in streaming or per message manner.

It is useful when implementing such as data obfuscation and permission control.
