A minimal and distributed key value store built on top of Bolt.
Bolt is an embedded key value store and this project build on top of it adding features like replicas, and sharding.
Furthermore, the project provides HTTP APIs to `set` and `get` values to and from the distributed key value store.
This project was helpful for me in practically learning distributed system concepts like partitioning, replication, etc.

### Sample requests

curl -Lkv 'http://localhost:8081/get?key=a'
curl -Lkv 'http://localhost:8081/set?key=a&value=58'
