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

cmd/
client/ → Client binary entrypoint
leader/ → Leader binary entrypoint
member/ → Member binary entrypoint

internal/
client/ → Client core logic
leader/ → Leader core logic (gRPC server, coordination)
member/ → Member core logic
common/ → Shared utilities and configuration

proto/
family/ → gRPC protobuf definitions and generated code



> All executable binaries live under `cmd/`  
> All application logic lives under `internal/`

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

## 🚀 Build & Run (Example)

```bash
go build ./cmd/leader
go build ./cmd/member
go build ./cmd/client
```