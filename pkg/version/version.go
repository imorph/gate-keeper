package version

var APPNAME = "gk"
var VERSION = "0.0.1"

var REVISION = "unknown"

// GetVersion returns version
func GetVersion() string {
	return VERSION
}

// GetRevision returns revision
func GetRevision() string {
	return REVISION
}

// GetAppName returns revision
func GetAppName() string {
	return APPNAME
}
