# Caching Proxy
A simple caching proxy server implemented in Go. Utilises a LRU cache in order to implement a max size limit. 

## Usage
```bash
go run cmd/caching-proxy/main.go proxy --port <PORT> --origin <ORIGIN_URL> --ttl <TTL> --memory <MEMORY_SIZE>
```
