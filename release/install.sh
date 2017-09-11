#!/usr/bin/env bash

echo "***********INSTALL*************"
echo "***** Content  Search**********"
echo "*****    开始安装  ************"


command -v wvText >/dev/null 2>&1 || { echo >&2 "需要安装wv."; sudo apt install wv links -y; }
command -v wps >/dev/null 2>&1 || { echo >&2 "需要安装wps."; sudo apt install wps-office -y; }

sudo apt-get install libgtk-3-dev libcairo2-dev libglib2.0-dev at-spi2-core

cd ../..
make
chmod u+x ./cs
mv ./cs release/cs

ExPATH=$(pwd)

cat <<EOF> ~/Desktop/Csearch.desktop
[Desktop Entry]
Version=0.3
Name=ContentSearch
Comment=An application for searching the content in doc or docx files
Exec=${ExPATH}/cs
Icon=./search.png
Terminal=false
Type=Application
X-Desktop-File-Install-Version=0.1
X-Deepin-AppID=ContentSearch
EOF
