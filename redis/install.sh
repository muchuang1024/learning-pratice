# !bin/bash

# 请安装redis，并创建redis-server redis-sentinel软连接到目录

mkdir /usr/local/var/db/redis/6379
mkdir /usr/local/var/db/redis/6380

# run redis

/usr/local/bin/redis-server ./conf/redis_6379.conf
/usr/local/bin/redis-server ./conf/redis_6380.conf

# run sentinel

/usr/local/bin/redis-sentinel ./conf/sentinel_26379.conf
/usr/local/bin/redis-sentinel ./conf/sentinel_26380.conf
/usr/local/bin/redis-sentinel ./conf/sentinel_26381.conf

