package version

var APPNAME = "unknown"
var VERSION = "v0.0.2"

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
