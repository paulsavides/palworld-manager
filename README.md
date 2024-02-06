## Installation

Only linux x64 binary available right now

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
