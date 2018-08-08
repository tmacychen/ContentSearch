#!/usr/bin/env bash

echo "***********INSTALL*************"
echo "***** Content  Search**********"
echo "*****    开始安装  ************"


command -v wvText >/dev/null 2>&1 || { echo >&2 "需要安装wv."; sudo apt install wv links -y; }
#command -v wps >/dev/null 2>&1 || { echo >&2 "需要安装wps."; sudo apt install wps-office -y; }

sudo apt-get install  libglib2.0-dev at-spi2-core

cd ../..
make
make install

