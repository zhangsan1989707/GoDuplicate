package core

// ScannerEngine defines the minimal capabilities shared by CLI and GUI.
// Concrete implementations can optimize for different modes while sharing the interface.
type ScannerEngine interface {
	// Scan performs the scan according to the provided configuration and returns duplicate groups.
	// Engines SHOULD invoke config.OnProgress if not nil to report progress.
	Scan(config ScanConfig) ([]DuplicateGroup, error)
}


