#!/bin/bash 

# Suited for Ubuntu 20.04 

## Temurin 17 JDK
mkdir -p /etc/apt/keyrings
wget -O - https://packages.adoptium.net/artifactory/api/gpg/key/public | tee /etc/apt/keyrings/adoptium.asc
echo "deb [signed-by=/etc/apt/keyrings/adoptium.asc] https://packages.adoptium.net/artifactory/deb $(awk -F= '/^VERSION_CODENAME/{print$2}' /etc/os-release) main" | tee /etc/apt/sources.list.d/adoptium.list
sudo apt update
sudo apt install apt-transport-https temurin-17-jdk zip unzip

## Install SDKman
curl -s "https://get.sdkman.io" | bash
source "$HOME/.sdkman/bin/sdkman-init.sh"

## Install Gradle 8.1.1
sdk install gradle 8.1.1
export PATH=$PATH:/opt/gradle/gradle-8.1.1/bin