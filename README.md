## Installation

This has only been tested on ubuntu 22. It also requires palworld being run as a systemd service.

Download
```
wget https://palworld-manager-assets.s3.amazonaws.com/release/palworld-manager
chmod +x ./palworld-manager
cp ./palworld-manager /usr/local/bin/
rm palworld-manager
```

Add to crontab
```
echo "*/2 * * * * root /usr/local/bin/palworld-manager monitor -m 95 >> /var/log/palworld-manager.log 2>&1" > /etc/cron.d/palworld-monitor
```

Check logs
```
tail -f /var/log/palworld-manager.log
```

## RCON Config File

This is required to use the `rcon` command or to enable broadcasting to the server for restarts

Create a file ~/.palworld-manager.yaml
```
rcon:
    host: localhost
    port: <rcon_port>
    password: <password>
```

From there, verify it's working by calling `palworld-manager rcon Info`
