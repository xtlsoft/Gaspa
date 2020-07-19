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
