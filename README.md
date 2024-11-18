# Ultra

Ultra is a collection of high-performance, production-ready Go packages designed to help you build better backend services.

## Features
- 📦 Modular design - use only what you need
- 🛠️ Production-tested tooling
- ⚡ Focus on performance and reliability
- 🔄 Consistent APIs across packages

## Installing
```go
go get github.com/fmdunlap/ultra
```

## Available Packages
* `ultra/log` - A structured, concurrent logger with support for colorization, extensible formatting, terminal
  colorization, multi-writer support with asynchronous logging, separate formatters for each writer, and more!

More packages coming soon!

### On the Horizon

* `ultra/config` - An init-and-forget config library that allows you to access environment variables from anywhere in
  your app. Configs can be broken into logical chunks, with user-definable prefixes, and conversion functions.

### Eventually

* `ultra/cache` - A generic key-val store that can be configured as in-memory, or connected to redis without (too much)
  additional configuration. Goal here is to make write-through caching as easy as possible
* `ultra/router` - An http.ServeMux with subrouters, route walking, route-mapped middleware, param helpers, and more!
  
### Considering

* `ultra/environ` - An application-state store that supports event-driven signals. Contextual access to global
  application data and services. "Am I running in dev, staging, or prod?" Use `environ`. "How can I access my
  `UserService` struct in this api code?" Use `environ`. "When I get to N concurrent users I need to run this func". Use
  `environ`.
* `ultra/auth` - Set-and-forget auth functions with pre-defined middleware recipes.

## Versioning

Ultra uses [Semantic Versioning](https://semver.org/) for versioning. All packages within the module are versioned
together to ensure compatibility across packages. This means that a breaking change in any package will require a major
version bump, and a non-breaking change will require a minor version bump.

We're planning on maintaining a v0.x.x version for the entire module while new packages are added, and transitioning to
the first v1.0.0 once a critical mass of packages has been added. Until then, consider the package experimental.

## Contributing

Feel free to open an issue or PR if you have any suggestions or improvements!

## License

Ultra is licensed under the MIT License. See [LICENSE](LICENSE) for more details.