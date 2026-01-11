# Distributed Disk Register with gRPC

This project is a **distributed, fault-tolerant disk-based register system**
implemented in **Go (Golang)** using **gRPC**.

It is designed as part of a *Systems Programming* course and follows a
**Leader–Member architecture** with dynamic membership support.

---

## 🎯 Project Goals
- Leader–Member based distributed system
- Fault-tolerant message replication
- Persistent disk storage
- Dynamic member join capability
- Clear separation between binaries and core logic

---

## 🏗️ System Architecture

### Client
- Sends `SET` and `GET` requests to the Leader
- Communicates via a simple text-based protocol

### Leader
- Accepts client requests
- Distributes messages to Members via gRPC
- Ensures fault tolerance based on `tolerance.conf`
- Responds with `OK` or `ERROR` to the Client

### Member
- Receives replicated messages from the Leader
- Persists messages on disk
- Periodically reports its state to the Leader

---

## 🔌 Communication Model
- **Client ↔ Leader**: Text-based protocol
- **Leader ↔ Member**: gRPC (Protocol Buffers)

---

## 📁 Project Structure

```
cmd/
  main.go → Unified binary (auto-detects leader or follower role)

internal/
  leader/
    leader_coordinator.go → Message coordination and replication
    leader_grpc.go → gRPC communication to followers
    leader_tcp.go → TCP client handler
    leader_utils.go → Utilities (printer, etc)
  node/
    family_service.go → gRPC service implementation
    registry.go → Cluster membership tracking
    join.go → Dynamic join protocol
    health.go → Health checking
  common/
    command.go → Command parsing (SET/GET)
  config/
    config.go → Tolerance configuration
  discovery/
    discovery.go → Available port discovery
  storage/
    disk.go → Persistent message storage

proto/
  family.proto → gRPC protobuf definitions
  family/ → Generated Go code from protobuf

messages/
  {id}.msg → Persisted messages on disk
```

> Single binary (`main.go`) determines role dynamically:
> - If port 5555 is available → **LEADER**
> - If port 5555 is occupied → **FOLLOWER** (joins cluster)

---

## 🧠 Development Workflow
- Project management is handled via **GitHub Projects**
- Each task is tracked as a **ToDo item**
- Development is done on **feature branches**
- Completed features are merged via **Pull Requests**
- Tasks are moved to **Done** after successful merge

---

## ⚙️ Technologies Used
- Go (Golang)
- gRPC
- Protocol Buffers
- Git & GitHub Projects

---

## 🚀 Build & Run

### Build
```bash
go build ./cmd/main
```

### Run

**Terminal 1 - Start Leader (automatically detects):**
```bash
./main
# Output: [ROLE] LEADER on port 5555
```

**Terminal 2 - Start Follower (automatically joins):**
```bash
./main
# Output: [ROLE] FOLLOWER on port 5556 (or next available)
```

**Terminal 3 - Send Commands via TCP:**
```bash
# SET command
echo "SET 1 hello" | nc localhost 6666

# GET command
echo "GET 1" | nc localhost 6666
```

### Configuration
Edit `tolerance.conf` to set replication tolerance:
```
tolerance=2
```
(Messages will be replicated to N followers based on tolerance)