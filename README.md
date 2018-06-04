# passhash
simple sha512 password hasher

Build:

```
go build
```

Run:

```
./passhash
```

`/hash` endopint usage:

```
% curl -d "password=angryMonkey" localhost:8082/hash
12

% curl localhost:8082/hash/12

(wait 5 seconds)
% curl localhost:8082/hash/12
ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==
```

`/stats` endpoint reports total requests and average processing time in *nanoseconds*

```
% curl localhost:8082/stats
{"total":3,"average":5003883407}
```

`/shutdown` endpoint starts a shutdown process... no new requests will be accepted.  server will shut down after existing requests are processed.


Configuration:

The configuration file 'config.json' contains two configuration properties:

port - the port that the server should listen on (defaults to 8081 if not specified)
hashSize - the maximum size of the password hash map that will be saved in memory before being reused (defaults to 1 million if not specified)
