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

const CYPHERAPPS_REPO = "git://github.com/SatoshiPortal/cypherapps.git"

const VERSION = "0.1.0"
const AUTHOR = "SKP <skp@skp.rocks>"
const NAME = "cam - Cyphernode apps management tool"
const DESCRIPTION = "A tool to manager your cypherapps"

const DATA_DIR = ".cam"
const REPO_DIR = "repo"
const STATE_FILE = "state.json"
const SOURCE_FILE = "sources.list"
const REPO_INDEX_FILE = REPO_DIR+string(filepath.Separator)+"index.json"
const LOCK_FILE = "state.lock"

const APP_DESCRIPTION_FILE = "app.json"
const APP_VERSIONS_DIR = "versions"
const CANDIDATE_DESCRIPTION_FILE = "candidate.json"

const INSTALL_DIR_ENV_KEY = "CYPHERAPPS_INSTALL_DIR"
const INSTALLED_APPS_FILE = "index.json"
const INSTALL_DIR = "apps"
const KEYS_FILE_ENV_KEY = "CYPHERNODE_KEYS_FILE"
const CYPHERNODE_INFO_FILE_ENV_KEY = "CYPHERNODE_INFO_FILE"
const DOCKER_COMPOSE_TEMPLATE_REGEXP_TEMPLATE = `<%%= *%s *%%>`


var DockerVolumeWhitelist = []string{
  `^\$(\{ *|)GATEKEEPER_CERTS_PATH(| *\})`, // everything beneath GATEKEEPER_CERTS_PATH
  `^\$(\{ *|)UNSAFE__CLIGHTNING_PATH(| *\})`, // everything beneath UNSAFE__CLIGHTNING_PATH
  `^\$(\{ *|)APP_DATA(| *\})`, // everything beneath APP_DATA
}

var DockerVolumeElementBlacklist = []string{
 "..", // something fishy is going on, maybe trying sth like $APP_DATA/../../
}