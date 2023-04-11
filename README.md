[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/EPA507O7)
# Assignment: Ethernet Switch

![Ethernet switch animation](http://www.fiber-optic-solutions.com/wp-content/uploads/2018/05/ethernet-switch.gif)
![Ethernet hub animation](http://www.fiber-optic-solutions.com/wp-content/uploads/2018/05/hub.gif)

## Overview

Ethernet resides in the L1 and L2 of the network stack.  An IP packet is sent over the network as the payload of an ethernet frame.  Initially the world had ethernet hubs, which is a L1 device.  When an ethernet frame was received by a hub it was broadcast to all other ports on the hub.  This is inefficient since often the frame was only intend for one destination.  An advance in networking technology came when the ethernet switch was developed, a L2 device.  The ethernet switch keeps a record of what destinations are connected to which ports.  Therefore when a ethernet frame arrives at a port on an ethernet switch, it is only sent on the port that is connected to the destination.

In this assignment you will write a software implementation of an ethernet switch.  The switch will have an arbitrary number of ports and will use store and forward packet switching with a configurable send queue buffer size.  Store and forward simply means that the packets is fully read by the port on the switch and then forwarded to the appropriate port by looking at the forwarding table.  Each port has a send queue (of configurable finite length) that is used to buffer frames as they are sent over the wire.  The send queue is useful for providing a little buffer for slower ports (running at a lower line speed) and for the case when lots of incoming ports have frames destined for the same outgoing port.  Ethernet switches do best effort delivery therefore if the send queue for a port is full then you must drop the packet and move on to processing the next frame.

## Learning Objectives

- Manage concurrency in Go with channels

## Requirements

- Your switch design needs to be non-blocking.  This means that your ports must always read data from the outside without blocking on other activities (e.g., sending frames) within the switch.
- The design must be free of deadlock, races, and livelock.
- The switch must support IEEE 802.3 ethernet frames (L2) defined [here](https://en.wikipedia.org/wiki/Ethernet_frame#Preamble_and_start_frame_delimiter).  You only need to support payloads less than or equal to 1500 B (so no jumbo frame support).
- Upon reading the frame you must validate the checksum in the frame.  If the frame is invalid do not forward the frame.  Instead simply drop it.
- The switch must handle the special broadcast frame properly.
- All ports must be full duplex.

- Your code will be graded on completeness and form.
- You must only edit `pkg/eth/types.go` and `cmd/eth/switch_test.go`.

- Do not use any library other than the Go standard library.
- The source code must compile with the most recent version of the Go compiler.
- The program must not panic under any circumstances.
- Make sure your code is "gofmt'd".  See "gofmt" or better use "goimports" or better yet configure IDE to do this formatting on file save.
- Commit and push your working code to your GIT repository.

## Hints

- The package "encoding/binary" and "hash/crc32" might come in handy when parsing and formatting ethernet frames.
- Use the `Makefile`.  For example, you can run tests with `make test -B`.  To make sure your code is properly formatted run `make fix`.

## Submission

- Commit and push your working code to your GIT repository.
- Ensure all tests pass otherwise you will receive no credit.
