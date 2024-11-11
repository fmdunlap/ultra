# Ultra

Ultra is a collection of high-performance, production-ready Go packages designed to help you build better backend services.

## Features
- 📦 Modular design - use only what you need
- 🛠️ Production-tested tooling
- ⚡ Focus on performance and reliability
- 🔄 Consistent APIs across packages

## Installing
```go
import "github.com/fmdunlap/ultra/log" // for the logger
```

## Available Packages
* log - A structured, concurrent logger with support for colorization, extensible formatting, terminal colorization,
  multi-writer support with asynchronous logging, separate formatters for each writer, and more!

More packages coming soon!

## Versioning

Ultra uses [Semantic Versioning](https://semver.org/) for versioning. All packages within the module are versioned
together to ensure compatibility across packages. This means that a breaking change in any package will require a major
version bump, and a non-breaking change will require a minor version bump.

We're planning on maintaining a 1.x.x version for the entire module while new packages are added.

## Contributing

Feel free to open an issue or PR if you have any suggestions or improvements!

## License

Ultra is licensed under the MIT License. See [LICENSE](LICENSE) for more details.