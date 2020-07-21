# Gaspa Bridge Protocol

## Overview

Gaspa Bridge offers a simple protocol based on TCP.

## First Packet

Gaspa Bridge will determine all the headers it needs by the first packet, ending with `|`.

The first byte indicates which service to request.

| Byte Content (in Hex) | Service    | Description                                                  |
| --------------------- | ---------- | ------------------------------------------------------------ |
| `0x72` (`r`)          | register   | Register a new node to this bridge node.                     |
| `0x6d` (`m`)          | meta       | Find all the nodes this bridge node have.                    |
| `0x71` (`q`)          | meta_query | Check if a node has had a connection with this bridge node.  |
| `0x63` (`c`)          | connect    | Connect to another node (will establish another connection). |
| `0x6a` (`j`)          | join       | Join a connection created by another node.                   |

And the the following bytes will be the arguments, usually ending with `!`.

## Services

### register

After the first packet, the second packet should contain the information of the client.

The packet will consist of 2 parts, with the first showing the node's UUID and the second showing its name.

The UUID should be 16 bytes, sended in plain text.

The packet should end with a '!'.

Example:

```text
# packet 1 : dispatch to register service
+-------------------------+-----------------+
|        HEX  VIEW        |   ASCII  VIEW   |
+-------------------------+-----------------+
| 72 .. .. .. .. .. .. .. | r . . . . . . . |
+-------------------------+-----------------+

# packet 2 : arguments to pass
+-------------------------+-----------------+
|        HEX  VIEW        |   ASCII  VIEW   |
+-------------------------+-----------------+
| 8D D5 35 77 EC 10 4E FF | t E . . V R K I |
| 81 3B CA 44 0F 82 CF 17 | . . . Q @ . . . |
| 74 65 73 74 2D 6E 6F 64 | t e s t - n o d |
| 65 21 .. .. .. .. .. .. | e ! . . . . . . |
+-------------------------+-----------------+

# combined
+-------------------------+-----------------+
|        HEX  VIEW        |   ASCII  VIEW   |
+-------------------------+-----------------+
| 72 8D D5 35 77 EC 10 4E | r t E . . V R K |
| FF 81 3B CA 44 0F 82 CF | I . . . Q @ . . |
| 17 74 65 73 74 2D 6E 6F | . t e s t - n o |
| 64 65 21 .. .. .. .. .. | d e ! . . . . . |
+-------------------------+-----------------+
```

Above example registers a machine with UUID `af602678-15f6-4ef2-a01b-af56346d8330` and name `test-node`.

The bridge will send one packet to response, containing:

```text
+-------------------------+-----------------+
|        HEX  VIEW        |   ASCII  VIEW   |
+-------------------------+-----------------+
| 41 21 .. .. .. .. .. .. | A ! . . . . . . |
+-------------------------+-----------------+
```

The established connection will be re-used for notifying. It will be called `node_conn` in this document.

### meta

Meta query has no arguments.

The result will be presented in JSON.
