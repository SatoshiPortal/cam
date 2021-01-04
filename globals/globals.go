/*
 * MIT License
 *
 * Copyright (c) 2020 schulterklopfer/__escapee__
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILIT * Y, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package globals

import "path/filepath"

const (
  VERSION                                 = "0.1.0"
  AUTHOR                                  = "SKP <skp@skp.rocks>"
  NAME                                    = "cam - Cyphernode apps management tool"
  DESCRIPTION                             = "A tool to manager your cypherapps"
  DATA_DIR                                = ".cam"
  REPO_DIR                                = "repo"
  STATE_FILE                              = "state.json"
  SOURCE_FILE                             = "sources.list"
  REPO_INDEX_FILE                         = REPO_DIR + string(filepath.Separator) + "index.json"
  LOCK_FILE                               = "state.lock"
  APP_DESCRIPTION_FILE                    = "app.json"
  APP_VERSIONS_DIR                        = "versions"
  CANDIDATE_DESCRIPTION_FILE              = "candidate.json"
  INSTALL_DIR_ENV_KEY                     = "CYPHERAPPS_INSTALL_DIR"
  INSTALLED_APPS_FILE                     = "index.json"
  INSTALL_DIR                             = "apps"
  KEYS_FILE_ENV_KEY                       = "CYPHERNODE_KEYS_FILE"
  CYPHERNODE_INFO_FILE_ENV_KEY            = "CYPHERNODE_INFO_FILE"
  DOCKER_COMPOSE_TEMPLATE_REGEXP_TEMPLATE = `<%%= *%s *%%>`
  TRUST_ZONE_UNTRUSTED                    = "untrusted"
  TRUST_ZONE_TRUSTED                      = "trusted"
  TRUST_ZONE_SERVICE                      = "service"
  TRUST_ZONE_CORE                         = "core"
  CORE_NETWORK                            = "cyphernodenet"
  APPS_NETWORK                            = "cyphernodeappsnet"
  SERVICE_NETWORK                         = "cyphernodeservicenet"
)

const DOCKER_COMPOSE_MIDDLEWARE_PATTERN = "^traefik\\.http\\.middlewares\\.(<%= *APP_ID *%>-\\w+?)\\."

var (
  DOCKER_COMPOSE_ALLOWED_MAIN_SERVICE_LABELS = []string{
    DOCKER_COMPOSE_MIDDLEWARE_PATTERN,
    "^traefik\\.http\\.services\\.<%= *APP_ID *%>\\.loadbalancer\\.server\\.port=\\d+",
  }
)

const (
  DOCKER_COMPOSE_LABEL_TRAEFIK_ENABLE   = "traefik.enable=true"
  DOCKER_COMPOSE_LABEL_MOUNTPOINT_RULE  = "traefik.http.routers.<%= APP_ID %>.rule=PathPrefix(`/<%= APP_MOUNTPOINT %>`)"
  DOCKER_COMPOSE_LABEL_ENTRYPOINTS      = "traefik.http.routers.<%= APP_ID %>.entrypoints=web,websecure"
  DOCKER_COMPOSE_LABEL_MIDDLEWARES      = "traefik.http.routers.<%= APP_ID %>.middlewares="
  DOCKER_COMPOSE_LABEL_ROUTER_SERVICE   = "traefik.http.routers.<%= APP_ID %>.service=<%= APP_ID %>"
  DOCKER_COMPOSE_LABEL_MW_STRIPPREXIX   = "traefik.http.middlewares.<%= APP_ID %>-stripprefix.stripprefix.prefixes=/<%= APP_MOUNTPOINT %>,/<%= APP_MOUNTPOINT %>/"
  DOCKER_COMPOSE_LABEL_FORCE_SLASH      = "traefik.http.middlewares.<%= APP_ID %>-stripprefix.stripprefix.forceSlash=true"
  DOCKER_COMPOSE_LABEL_PASS_HOST_HEADER = "traefik.frontend.passHostHeader=true"
  DOCKER_COMPOSE_STRIPPREFIX_MIDDLEWARE = "<%= APP_ID %>-stripprefix@docker"
)

var (
  TRUST_ZONE_CORE_PATTERN    = `^\$(\{ *|)CORE__.+?(| *\})` // everything beneath TRUSTED__<ANYTHING>_PATH
  TRUST_ZONE_TRUSTED_PATTERN = `^\$(\{ *|)TRUSTED__.+?(| *\})` // everything beneath TRUSTED__<ANYTHING>_PATH
  TRUST_ZONE_SERVICE_PATTERN = `^\$(\{ *|)SERVICE__.+?(| *\})` // everything beneath SERVICE__<ANYTHING>_PATH

  DockerVolumeWhitelist = []string{
    `^\$(\{ *|)GATEKEEPER_CERTS_DATAPATH(| *\})`,    // everything beneath GATEKEEPER_CERTS_PATH
    `^\$(\{ *|)APP_DATAPATH(| *\})`,            // everything beneath APP_DATA_PATH
    TRUST_ZONE_CORE_PATTERN,
    TRUST_ZONE_TRUSTED_PATTERN,
    TRUST_ZONE_SERVICE_PATTERN,
  }

  // TODO: research if \.\. or \\.\\. is bad
  // TODO: check env files
  DockerVolumeElementBlacklist = []string{
    "..", // something fishy is going on, maybe trying sth like $APP_DATA/../../
  }
  ValidTrustZones = []string{
    TRUST_ZONE_UNTRUSTED, TRUST_ZONE_TRUSTED, TRUST_ZONE_SERVICE, TRUST_ZONE_CORE,
  }
  DefaultTrustZone = TRUST_ZONE_UNTRUSTED
)
