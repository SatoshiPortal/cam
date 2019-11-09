package errors

import (
  goErrors "errors"
)

var DATADIR_IS_NOT_A_DIRECTORY = goErrors.New( "data directory is not a directory")
var DATADIR_DOES_NOT_EXIST = goErrors.New( "data directory does not exist")
var DATADIR_IS_LOCKED = goErrors.New( "data directory is locked" )
var DATADIR_IS_INVALID = goErrors.New( "invalid data directory" )

var NO_SUCH_SOURCE_TYPE = goErrors.New( "no such source type" )
var NO_SUCH_SOURCE = goErrors.New( "no such source" )
var DUPLICATE_SOURCE = goErrors.New( "duplicate source" )

var SOURCE_ADD_NO_SOURCE = goErrors.New( "source add: no source" )
var SOURCE_DELETE_NO_SOURCE = goErrors.New( "source delete: no source" )