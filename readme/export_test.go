package readme

// Export internal functions for testing.
var (
	GetVersion  = (*VersionClient).getVersion
	ParseUUID   = parseUUID
	HasNextPage = hasNextPage
)
