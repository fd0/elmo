package list

import "strings"

const helpShort = "Display information about repositories on the server"

var helpLong = strings.TrimSpace(`
The 'list' command connects to a gitolite server and displays information about
the repositories accessible to the user.

If a pattern is given, only matching repositories are printed.
`)

const helpExamples = ""
