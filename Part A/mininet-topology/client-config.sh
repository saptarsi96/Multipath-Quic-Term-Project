#!/bin/bash

ifconfig client-eth0 10.0.0.1
ifconfig client-eth1 11.0.0.1

ip rule add from 10.0.0.1 table 1
ip rule add from 11.0.0.1 table 2

ip route add 10.0.0.0/8 dev client-eth0 scope link table 1
ip route add default via 10.0.0.2 dev client-eth0 table 1

ip route add 11.0.0.0/8 dev client-eth1 scope link table 2
ip route add default via 11.0.0.2 dev client-eth1 table 2

ip route add default scope global nexthop via 10.0.0.2 dev client-eth0

