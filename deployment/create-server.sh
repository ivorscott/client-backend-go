#!/bin/bash
set -x

# create server
docker-machine create \
  --driver=digitalocean \
  --digitalocean-access-token=`cat ./secrets/do_token` \
  --digitalocean-size="2gb" \
  --digitalocean-ssh-key-fingerprint=`cat ./secrets/ssh_fingerprint` \
  --digitalocean-tags=gdr \
  --digitalocean-private-networking=true \
  --digitalocean-region="lon1" \
  gdr1 &
