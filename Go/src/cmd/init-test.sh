#!/usr/bin/env bash


cd /tmp
mkdir e
touch a.bat
echo `date` > a.bat

touch e/c.bat
echo `date` > e/c.bat

cp /home/ycl/Applications/Distribute/src/cmd/distribute.json .

tar -czf distribute03222.tar.gz e distribute.json a.bat


mkdir {app,web}