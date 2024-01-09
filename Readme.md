## Development/Testing/Networking/Entertainment tools
<hr>

## BST Load Balancer

Binary Search Tree (BST) based load balancer designed to distribute incoming network traffic across multiple servers more efficiently

## Network Listener Program

The Network Listener Program is a lightweight utility that monitors network activity.
<br>
In plans, to made it firewall.

## Open Port Scanner

Unexpectedly, scans target systems for open ports.

## Forwarding/Reverse Proxy Servers

The forwarding proxy forwards client requests to the destination server, while the reverse proxy handles requests on behalf of servers.

## Unit Testing Tool

Currently working on it.

## IP Tracker
Simply enter an IP, and the exact location of the target will be printed

## DDos attack
Very very basic ddos attack tool. Unit testing tool may be used as DDos attack tool more efficiently that this :D

## Man-In-The-Middle attack
Firstly, retrieve mac addresses of both machines. <br>
Initial configurations, optional configurations, and start sending ARP packets.
<br>
Example of usage (don't forget to change mode to u+x) 
```shell
./main.py --interface <INTERFACE_NAME> --target1 <TARGET1_IP> --target2 <TARGET2_IP> 
```
where `INTERFACE_NAME` retrieved from `ifconfig`, and `TARGET1, TARGET2` is ip addresses representing server & client.


### You can simply copy each component right into Your code, or if You want to test:
```go
cd {folder name}
go run .
```

## Contibutions are much welcome! Please, share in `Issues` what tools would You like to see here.
