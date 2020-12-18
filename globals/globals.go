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
  TRUST_ZONE_CORE                         = "core"
  CORE_NETWORK                            = "cyphernodenet"
  APP_NETWORK                             = "cyphernodeappsnet"
  TRUSTED_APP_NETWORK                     = "cyphernodetrustedappsnet"
)

var (
  DockerVolumeWhitelist = []string{
    `^\$(\{ *|)GATEKEEPER_CERTS_PATH(| *\})`,    // everything beneath GATEKEEPER_CERTS_PATH
    `^\$(\{ *|)TRUSTED__CLIGHTNING_PATH(| *\})`, // everything beneath UNSAFE__CLIGHTNING_PATH
    `^\$(\{ *|)APP_DATA(| *\})`,                 // everything beneath APP_DATA
  }
  // TODO: research if \.\. or \\.\\. is bad
  // TODO: check env files
  DockerVolumeElementBlacklist = []string{
    "..", // something fishy is going on, maybe trying sth like $APP_DATA/../../
  }
  ValidTrustZones = []string{
    TRUST_ZONE_UNTRUSTED, TRUST_ZONE_TRUSTED, TRUST_ZONE_CORE,
  }
  DefaultTrustZone = TRUST_ZONE_UNTRUSTED
)
