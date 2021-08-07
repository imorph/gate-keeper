package version

var appname = "unknown"
var version = "v1.0.10"

var revision = "unknown"

// GetVersion returns version
func GetVersion() string {
	return version
}

// GetRevision returns revision
func GetRevision() string {
	return revision
}

// GetAppName returns revision
func GetAppName() string {
	return appname
}
