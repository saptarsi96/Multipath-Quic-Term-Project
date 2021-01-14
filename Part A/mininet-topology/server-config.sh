#!/bin/bash

ifconfig server-eth0 100.0.0.1
route add default gw 100.0.0.2 server-eth0