package globals

import "path/filepath"

const CYPHERAPPS_REPO = "git://github.com/SatoshiPortal/cypherapps.git"

const VERSION = "0.1.0"
const AUTHOR = "SKP <skp@skp.rocks>"
const NAME = "cna - Cyphernode apps management tool"
const DESCRIPTION = "A tool to manager your cypherapps"

const DATA_DIR = ".cna"
const REPO_DIR = "repo"
const STATE_FILE = "state.json"
const SOURCE_FILE = "sources.list"
const REPO_INDEX_FILE = REPO_DIR+string(filepath.Separator)+"index.json"
const LOCK_FILE = "state.lock"

const APP_DESCRIPTION_FILE = "app.json"
const APP_VERSIONS_DIR = "versions"
const CANDIDATE_DESCRIPTION_FILE = "candidate.json"