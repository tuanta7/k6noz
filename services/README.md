# Location

References:

- [System Design School | Design Uber, Lyft](https://systemdesignschool.io/problems/uber/solution)
- [Ably | Scaling WebSockets](https://ably.com/topic/the-challenge-of-scaling-websockets)

## 1. High Level Design

## 2. WebSocket

WebSocket services are used to maintain a live feed of driver locations for both drivers (update their location) and passengers (get the driver updates). In real-world production like Uber, QUIC/HTTP3 - a more modern technology is used instead. 

### 2.1. General Considerations

- **Client-side TCP Port Limit**: A single server listening on a single port (only one fixed IP) can theoretically handle up to 65,535 (16-bit) *concurrent* connections from distinct client IP addresses, as this is the maximum number of available TCP ports on the client side.
- **Server-side File Descriptor Limit**: Every open socket consumes a file descriptor. On Linux/macOS, the default soft limit (ulimit -n) is often 1024 and the recommend hard limit is 65,535. There is effectively no fixed port-based limit on the server side for WebSocket concurrency. 
- **Scaling**: For very high numbers of connections, horizontal scaling is typically employed to distribute connections across multiple servers, often with load balancers to manage incoming connections.

### 2.2. Multi Servers Sticky-session 

Sticky session, also known as session affinity is a load balancing strategy where a client consistently connects to the same server across a session or multiple reconnects. 

### 2.3. Fallback Transport Mechanisms

Some clients, due to restrictive firewalls, proxies, or legacy environments won't be able to establish a WebSocket connection at all.

## 3. Trip History Storage: ClickHouse

A time-series/column oriented database like TimescaleDB (PostgreSQL extension) or ClickHouse is recommended. This storage pattern allows efficient historical queries for analytics or trip reconstruction.

## 4. Latest Location Storage: Redis

A high-performance key-value store is preferred because only the latest update is needed, overwriting the record is sufficient.

## 5. Messaging Layer: Kafka

A location-tracking workload behaves like a high-frequency telemetry stream. Kafka's partitioned log design allows location updates from thousands of drivers can be processed without pressure on the broker.
