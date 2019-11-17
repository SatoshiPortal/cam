#!/usr/bin/env sh

docker run -e CYPHERAPPS_INSTALL_DIR=/apps -v $(pwd)/apps:/apps -v $(pwd):/data --rm cam $*
