# GregorMQ 🪳
*A lightweight message broker in Go — in honor of Gregor Samsa from Kafka’s Metamorphosis.*

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/gregormq)](https://goreportcard.com/report/github.com/yourusername/gregormq)
[![Build Status](https://github.com/yourusername/gregormq/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/gregormq/actions)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Coverage Status](https://img.shields.io/codecov/c/github/yourusername/gregormq/main.svg)](https://codecov.io/gh/yourusername/gregormq)

---

## Overview
GregorMQ is a **lightweight, cloud-friendly message broker** written in Go.  
It aims to provide the power of Kafka and RabbitMQ in a **simpler, faster, and smaller** package.

⚠️ **Note:** The project is in its **early stages of development**. Many features are planned but not yet implemented.

- **Pub/Sub model**
- **Persistence optional**
- **Clustering via gossip protocol**
- **Zero configuration** (works out of the box)
- **Docker & Kubernetes ready**

---

## Quickstart

### Run with Docker
```bash
docker run -d --name gregormq -p 7777:7777 ghcr.io/yourusername/gregormq:latest
```

### Publish a message
```bash
gregorctl publish --topic greetings --message "Hello Gregor!"
```

### Subscribe to a topic
```bash
gregorctl subscribe --topic greetings
# Output:
# [greetings] Hello Gregor!
```

---

## Roadmap
See [docs/roadmap.md](docs/roadmap.md) for the full feature roadmap.

---

## Repository Structure
```
gregormq/
├── cmd/
│   └── broker/        # broker main
│   └── gregorctl/     # CLI client
├── pkg/
│   ├── protocol/      # binary protocol definition
│   ├── broker/        # core pub/sub
│   ├── storage/       # log persistence
│   ├── discovery/     # gossip
│   ├── metrics/       # prometheus integration
├── client/
│   └── go/            # Go SDK
├── docs/
│   ├── architecture.md
│   ├── roadmap.md
│   └── protocol.md
└── tests/
    └── integration/
```

---

## Contributing
Contributions are welcome!
- Open an issue for ideas, bugs, or feature requests
- Submit a PR following the [contributing guide](CONTRIBUTING.md)

---

## License
GregorMQ is released under the [Apache 2.0 License](LICENSE).

---

## Links
- Documentation: coming soon
- Roadmap: [docs/roadmap.md](docs/roadmap.md)
- Blog post: coming soon

---
