#!/bin/bash

wget http://packages.couchbase.com/releases/4.6.2/couchbase-server-enterprise_4.6.2-ubuntu12.04_amd64.deb
dpkg-deb -x couchbase-server-enterprise_4.6.2-ubuntu12.04_amd64.deb $HOME
cd $HOME/opt/couchbase
./bin/install/reloc.sh `pwd`
./bin/couchbase-server -- -noinput -detached
sleep 20
./bin/couchbase-cli cluster-init -c 127.0.0.1:8091  --cluster-init-username=Administrator --cluster-init-password=password --cluster-init-port=8091 --cluster-init-ramsize=1024 --cluster-index-ramsize=1024 --services=data,index,query,fts
./bin/couchbase-cli bucket-create -c 127.0.0.1:8091 --bucket=test --bucket-type=couchbase --bucket-port=11211 --bucket-ramsize=512  --bucket-replica=1 -u Administrator -p password --wait
