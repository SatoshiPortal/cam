package output

import "github.com/fatih/color"

var Errorf  = color.New(color.Bold, color.FgRed).PrintfFunc()
var Noticef = color.New().PrintfFunc()
var Warningf= color.New(color.FgYellow).PrintfFunc()
var Successf= color.New(color.Bold, color.FgGreen).PrintfFunc()

var Error  = func( text string ) { Errorf( "%s\n", text) }
var Notice = func( text string ) { Noticef( "%s\n", text) }
var Warning= func( text string ) { Warningf( "%s\n", text) }
var Success= func( text string ) { Successf( "%s\n", text) }

