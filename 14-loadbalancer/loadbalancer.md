# Load Balancer

Basic load balancer util, using Round-Robin algorithm.

## Implementation

### Structs

`Loadbalancers` - `Port`, `Servers`, `RoundRobinCount` 

`Servers` - `Address`, `Proxy`

### Interfaces

`Server`

`Address` - Get the address of a server

`IsAlive` - Get if a server is still alive

`Serve` - Serve the server

### Path

Create server list -> Server proxy function -> Get next available server (Round robin algorithm)

