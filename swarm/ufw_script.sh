#!/bin/bash

echo "Opening firewall ports"
ufw allow 2376/tcp && \
ufw allow 2377/tcp && \
ufw allow 7946/tcp && \
ufw allow 7946/udp && \
ufw allow 4789/udp && \
ufw reload && \
ufw --force enable
echo "Finished opening firewall ports"
