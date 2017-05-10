# whoami - simple service discovery

Checks and sends host information (IP) to a key-value store server. 

   - The key-value store uses [kvstore](https://github.com/peteretelej/kvstore)


## Use case
Get the IP address of nodes on the network that have dynamic or unknown IPs.


# Usage 
- Install the [kvstore](https://github.com/peteretelej/kvstore) server on the server to store info
- Install `whoami`: 
```
go get github.com/peteretelej/whoami
```
- Run whoami:
```
whoami
# checks if the ip has changed and submits the new one to the kvstore, that's all
```

