## ️ Roadmap

### **M1 – Core Broker (single node)**
- [ ] Define simple binary protocol (`[HEADER][PAYLOAD]` with magic, timestamp, key, value)  
- [ ] Basic TCP connection (client → broker)  
- [ ] Message publishing on a topic  
- [ ] Subscription and message delivery  
- [ ] Basic acknowledgments  

*Goal*: minimal working broker  

---

### **M2 – Persistence & Reliability**
- [ ] Append-only log file persistence  
- [ ] Configurable retention (time-based, size-based)  
- [ ] Replay messages to new subscribers  
- [ ] Recovery after restart  

*Goal*: basic resilience  

---

### **M3 – Clustering & Gossip**
- [ ] Gossip protocol (Scuttlebutt-style) for node discovery  
- [ ] Message replication across nodes  
- [ ] Automatic failover  
- [ ] Zero-touch configuration: add node → cluster discovers it  

*Goal*: distributed broker  

---

### **M4 – Advanced Features**
- [ ] Topic partitions (parallel consumption)  
- [ ] Consumer groups (Kafka-style)  
- [ ] Message compression (Snappy/LZ4)  
- [ ] Prometheus metrics + Grafana dashboard  

*Goal*: production-ready lightweight broker  

---

### **M5 – Ecosystem & Developer Experience**
- [ ] Client SDK in Go (then JS, Python, Java)  
- [ ] CLI `gregorctl` for topics, cluster, metrics  
- [ ] Minimal Docker image (scratch/alpine)  
- [ ] Helm chart for Kubernetes deploy  

*Goal*: developer-friendly & cloud-native  

---

## ️ Repository Structure
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

