#!/usr/bin/env sh

docker run -e CYPHERAPPS_INSTALL_DIR=/install -v $(pwd)/install:/install -v $(pwd):/data --rm cna $*
