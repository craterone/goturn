
```
usage: gostun [<flags>]

实时猫 TURN/STUN 服务器

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=config.json,/etc/goturn/config.json
                                 Configuration file location
  -p, --port=3478                Server port
  -x, --external_ip=EXTERNAL_IP  TURN Server public/private address mapping, if the server is behind NAT.
  -r, --relay_ip=RELAY_IP        Relay endpoint ip
      --min_port=49152           Lower bound of the UDP port range for relay endpoints allocation.
      --max_port=65535           Upper bound of the UDP port range for relay endpoints allocation.
  -v, --version                  Show application version.
```