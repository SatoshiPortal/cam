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

const KEYS_FILE_ENV_KEY = "CYPHERNODE_KEYS_FILE"

var DockerVolumeWhitelist = []string{
  `^\$(\{ *|)GATEKEEPER_CERT_FILE(| *\})$`, // exact match
  `^\$(\{ *|)CLIGHTNING_RPC_SOCKET(| *\})$`, // exact match
  `^\$(\{ *|)APP_DATA(| *\})`, // everything beneath APP_DATA
}

var DockerVolumeElementBlacklist = []string{
 "..", // something fishy is going on, maybe trying sth like $APP_DATA/../../
}