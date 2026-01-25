#!/bin/bash

# Get the internal IP of the active interface
INTERNAL_IP=$(hostname -I | awk '{print $1}')

if [ -z "$INTERNAL_IP" ]; then
    echo "Disconnected"
else
    echo "$INTERNAL_IP"
fi
