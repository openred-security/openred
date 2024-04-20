#!/bin/bash

chmod +x gradlew
#Generate binary tarballs at build/distribution/local/opensearch-*-SNAPSHOT/
./gradlew localDistro

#Compress if required
#tar -cvf opensearch-2.6.0-SNAPSHOT.tar.gz opensearch-2.6.0-SNAPSHOT/

