# Registry Usage Example

This example demonstrates how to use the send registry to register and send packets.

## Overview

The example shows:

1. How to create a HandlerRegistry
2. How to register all handlers
3. How to configure a server type
4. How to send packets

## Components

### SimpleLogger

A simple implementation of the core.Logger interface that logs messages to stdout.

### MockConnection

A mock implementation of a connection that prints the packets it would send.

### Main Function

The main function demonstrates:

1. Creating a hook manager
2. Creating a logger
3. Creating a BaseSend instance
4. Creating a HandlerRegistry
5. Registering all handlers
6. Getting packet definitions
7. Configuring a server type
8. Setting a connection
9. Sending packets

## Running the Example

To run the example:

```bash
cd goKore/network/send/examples/registry_usage
go run main.go
```

## Expected Output

```
[INFO] Registered all send handlers
[DEBUG] Registered login send handlers
[DEBUG] Registered game send handlers
[DEBUG] Registered additional game send handlers
[DEBUG] Registered server-specific send handlers
Sending packet: [0 0]
Sending packet: [0 0]
[SUCCESS] Successfully sent packets
```

## Notes

- This example uses mock implementations of the logger and connection.
- In a real application, you would use actual implementations of these interfaces.
- The packets sent are not actual packets, but placeholders for demonstration purposes.