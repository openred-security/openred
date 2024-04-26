#!/bin/bash

chmod +x gradlew
#Generate binary tarballs at build/distribution/local/opensearch-*-SNAPSHOT/
#./gradlew localDistro
./gradlew localDistro
tar -czvf opensearch-2.13.0.tar.gz -C distribution/archives/linux-tar/build/install/opensearch-2.13.0-SNAPSHOT/ actions/docker/release
#Generate docker image (use opensearch-build repository)
cd actions/docker/release/
sudo ./build-image-single-arch.sh -v 2.13.0 -f ./dockerfiles/opensearch.al2.dockerfile -p opensearch -a x64 -t opensearch-2.13.0.tar.gz

#Compress if required
#tar -cvf opensearch-2.6.0-SNAPSHOT.tar.gz opensearch-2.6.0-SNAPSHOT/

