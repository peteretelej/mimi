# mimi - simple service discovery

Checks and sends host information (IP) to a key-value store server. 

   - The key-value store uses [kvstore](https://github.com/peteretelej/kvstore) (simple key-value store over http)


## Use case
Get the IP address of nodes on the network that have dynamic or unknown IPs.
__mimi__: Swahili for __me__.


# Usage 
- Install the [kvstore](https://github.com/peteretelej/kvstore) server on the server to store information from the nodes/machines
- Install `mimi`: 
```
go get github.com/peteretelej/mimi
```
- Set KVSTORE and KVCRED environment variables to access the kvstore:
```
export KVSTORE=http://10.0.1.11:8080
export KVCRED=credential1
# assuming kvstore server is at 10.0.1.11 and has credential "credential1"
```

- Run mimi: launches agent that checks machine/node info every 5m and sends to kvstore
```
mimi
# checks if the ip has changed and submits the new one to the kvstore
# If KVSTORE or KVCRED are not set, this will simply return the IP

mimi -interval 1h 
# specify custom check interval of 1hour (default 5minutes)

mimi -myname machineXYZ
# specify custom name to identify this machine (default: machine's hostname)
```

