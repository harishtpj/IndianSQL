# 🗃️ IndianSQL

IndianSQL is a small relational database engine written from scratch in Go. The project is primarily 
intended as a learning exercise to understand how database systems work internally by implementing the 
storage engine, pager, B+Tree index and SQL execution engine without relying on existing database 
libraries.

# 🌟 Features

* Basic Table Creation and CRUD operations
* SQL parser (via [sqlparser](https://github.com/xwb1989/sqlparser)) and execution engine
* Persistent, single file storage (like SQLite)
* Slotted page storage layout
* B+Tree indexes
* Extensible storage architecture
* Standalone server with MySQL-compatible client protocol (via [go-mysql](https://github.com/go-mysql-org/go-mysql))

# ⚙️ Building

Clone the repo and build it via:

```bash
$ git clone https://github.com/harishtpj/IndianSQL.git
$ cd IndianSQL
$ go build ./cmd/indsql
```

# 🚀 Usage

Create or open a database via:

```bash
$ indsql test.idb
```

If no filename is provided, an *in-memory* transient database will be loaded. The DB engine supports
basic SQL DML statements and a few helpers.

You can also start a server via:

```bash
$ indsql server test.idb
```

The default server runs on `localhost:4405` with `user=root` and no password. This can be configured via
commandline flags or via a [`config.yml`](https://github.com/harishtpj/IndianSQL/blob/master/examples/config.yml).

> [!NOTE]
> IndianSQL currently implements the legacy MySQL authentication protocol. Recent MySQL clients default
> to newer authentication methods, so you may need to use an older client or explicitly select a
> compatible authentication plugin.


The examples/ directory contains sample SQL scripts, including the classic SCOTT schema 
(.sql and prebuilt .idb) for testing.

IndianSQL is primarily intended as a proof-of-concept database engine. The focus of the project is to 
understand and implement the core components of a relational database rather than to compete with 
production systems such as SQLite or MySQL. Consequently, only a subset of SQL is currently supported, 
although the architecture is designed to be extended over time.

# 📝 License

#### Copyright © 2026 [M.V.Harish Kumar](https://github.com/harishtpj). <br>
#### This project is [MIT](https://github.com/harishtpj/IndianSQL/blob/master/LICENSE) licensed.
