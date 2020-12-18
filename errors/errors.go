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

package errors

import (
  goErrors "errors"
)

var DIR_IS_NOT_A_DIRECTORY = goErrors.New( "directory is not a directory ;-)")
var DIR_DOES_NOT_EXIST = goErrors.New( "directory does not exist")
var DATADIR_IS_LOCKED = goErrors.New( "data directory is locked" )
var DATADIR_IS_INVALID = goErrors.New( "invalid data directory" )
var NO_SUCH_SOURCE_TYPE = goErrors.New( "no such source type" )
var NO_SUCH_SOURCE = goErrors.New( "no such source" )
var DUPLICATE_SOURCE = goErrors.New( "duplicate source" )
var SOURCE_ADD_NO_SOURCE = goErrors.New( "source add: no source" )
var SOURCE_DELETE_NO_SOURCE = goErrors.New( "source delete: no source" )
var NO_SUCH_APP = goErrors.New( "no such app" )
var NO_SUCH_KEY = goErrors.New( "no such key" )
var DUPLICATE_APP = goErrors.New( "duplicate app" )
var DUPLICATE_KEY = goErrors.New( "duplicate key" )
var SOURCE_ADD_NO_APP = goErrors.New( "source add: no app" )
var SOURCE_DELETE_NO_APP = goErrors.New( "source delete: no app" )
var REPO_INDEX_DOES_NOT_EXIST = goErrors.New( "repo index does not exist" )
var INSTALLED_APPS_INDEX_DOES_NOT_EXIST = goErrors.New( "installed apps index does not exist" )
var INSTALL_DIR_DOES_NOT_EXIST = goErrors.New( "install dir does not exist" )
var APP_SEARCH_NO_SEARCH_TERM = goErrors.New( "app search: no search term" )
var APP_INSTALL_NO_APP_ID = goErrors.New( "app install: no app id" )
var APP_ALREADY_INSTALLED = goErrors.New( "app already installed" )
var APP_NOT_INSTALLED = goErrors.New( "app is not installed" )
var NO_SUCH_VERSION = goErrors.New( "no such version")
var VOLUME_HAS_ILLEGAL_ELEMENTS = goErrors.New( "volume has illegal elements" )
var VOLUME_NOT_IN_WHITELIST = goErrors.New( "volume not in whitelist" )
var NO_KEYS_FILE = goErrors.New( "no cyphernode keys file found" )
var CYPHERNODE_INFO_FILE_DOES_NOT_EXIST = goErrors.New( "cyphernode info file does not exist" )
var APP_VERSION_IS_NOT_COMPATIBLE = goErrors.New( "app version is not compatible" )
var APP_MOUNTPOINT_BLOCKED = goErrors.New( "app mount point is used by other app" )
var SERVICE_NAME_NOT_UNIQUE = goErrors.New( "service name not unique" )
var APP_HAS_WRONG_TRUST_ZONE = goErrors.New( "app has wrong trust zone" )