A minimal and distributed key value store built on top of Bolt.
Bolt is an embedded key value store and this project builds on top of it adding features like replicas, and sharding.
Furthermore, the project provides HTTP APIs to `set` and `get` values to and from the distributed key value store.
This project was helpful for me in practically learning distributed system concepts like partitioning, replication, etc.

### Sample requests

```bash
# Get a key
curl -Lkv 'http://localhost:8081/get?key=a'
```

```bash
# Set a key
curl -Lkv 'http://localhost:8081/set?key=a&value=58'
```
