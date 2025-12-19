# Stateful

## 1. TCP

A TCP connection is considered persisted not through explicit coordination between machines, but through the continued validity of protocol state maintained independently at both endpoints.

For idle connections, TCP does not continuously exchange messages by default. The connection remains valid indefinitely unless one side sends a FIN or RST, or the network path fails and the failure is detected.

### Stateless HTTP

While built on top of TCP – a stateful protocol., HTTP is stateless itself.

## 2. Gateway/LB Behavior

The client establishes a TCP connection to the API gateway or load balancer. This is the only TCP connection visible to the client. Persistence is achieved through connection anchoring at the gateway layer, rather than through a direct end-to-end TCP connection between the client and the backend.

During the HTTP upgrade to WebSocket, the gateway selects a backend instance and opens (or reuses) a separate TCP connection to that backend. At this point, a connection pair is formed:

- Client ↔ Gateway (TCP)
- Gateway ↔ Backend (TCP)

After the upgrade completes, the gateway switches from HTTP request routing to stream proxy mode. Bytes received on one TCP connection are forwarded directly to the other TCP connection with minimal inspection.
