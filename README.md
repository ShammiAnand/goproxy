# GoProxy 🚀

<p align="center">
   <a href="#">
    <img src="https://raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher.png" alt="GoProxy Gopher" width="30"/>
  </a>
  <a href="https://golang.org/doc/go1.21">
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  </a>
  <a href="https://github.com/yourusername/goproxy/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/License-BSD--3--Clause-blue?style=for-the-badge" alt="License">
  </a>
</p>

GoProxy is a high-performance, feature-rich reverse proxy written in Go. Designed with extensibility in mind, it aims to provide a robust solution for modern web architectures.

## 📦 Releases

Latest release: [v0.0.5](https://github.com/shammianand/goproxy/releases/tag/v0.0.5)

## 🌟 Features

### Feature Roadmap

- ✅ Basic HTTP reverse proxying
- ✅ Configurable via YAML
- ✅ Structured logging with slog
- ✅ Customizable log levels and formats
- ✅ Load balancing (Round Robin)
- 🔜 Additional load balancing algorithms (Least Connections)
- 🔜 TLS/SSL support
- 🔜 Request/Response manipulation
- 🔜 Caching
- 🔜 Rate limiting
- 🔜 Metrics and monitoring (Prometheus integration)
- 🔜 Health checking
- 🔜 Circuit breaking

## 📋 Prerequisites

- Go 1.21 or higher

## 🛠 Installation

1. Clone the repository:
   ```
   git clone https://github.com/shammianand/goproxy.git
   ```
2. Change to the project directory:
   ```
   cd goproxy
   ```
3. Build the project:
   ```
   make build
   ```

## ⚙ Configuration

GoProxy uses a YAML configuration file. Here's a sample configuration:

```yaml
server:
  listen_addr: ":8080"
  read_timeout: 5
  write_timeout: 10
  idle_timeout: 120

proxy:
  target_addr: "http://localhost:8000"
  max_idle_conns: 100
  dial_timeout: 10

load_balancing:
  enabled: true
  algorithm: "round_robin"
  backends:
    - "http://backend1:8080"
    - "http://backend2:8080"
    - "http://backend3:8080"

logging:
  level: "info"
  format: "json"

# ... (other configuration options)
```

For a full list of configuration options, see the [Configuration Guide](docs/configuration.md).

## 🚀 Usage

1. Start GoProxy:
   ```
   ./goproxy -config=config.yaml
   ```
2. The proxy will start and listen on the configured address.

## 📊 Monitoring

(Coming soon) GoProxy will expose Prometheus metrics on a configurable endpoint.

## 🧪 Testing

Run the test suite:

```
make test
```

This includes unit tests, integration tests, and performance benchmarks.

### Performance Test Results

Our latest performance tests show impressive results:

- Throughput: ~932 requests per second
- Average Latency: 11.8ms
- Min Latency: 31.4µs
- Max Latency: 38.27ms
- Success Rate: 98.94%

These results were achieved with 100 concurrent connections and simulated backend latency of 1-20ms. The proxy demonstrates excellent performance and stability under load.

## 🛣 Roadmap

1. Phase 1 (Completed): Basic proxying, configuration, and logging
2. Phase 2 (Current): Load balancing and health checking
   - ✅ Round Robin load balancing
   - 🔜 Least Connections load balancing
   - 🔜 Health checking implementation
3. Phase 3: TLS support and request/response manipulation
4. Phase 4: Caching and rate limiting
5. Phase 5: Metrics, monitoring, and advanced features (circuit breaking)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📜 License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.

## 👏 Acknowledgments

- The Go team for the amazing language and standard library
- The authors of the third-party libraries used in this project

---

<p align="center">
  Made with ❤️ by Shammi Anand
</p>
