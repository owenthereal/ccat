package toolchain

// ConfigFilename is the filename of the toolchain configuration file. The
// presence of this file in a directory signifies that a srclib toolchain is
// defined in that directory.
const ConfigFilename = "Srclibtoolchain"

// Config represents a Srclibtoolchain file, which defines a srclib toolchain.
type Config struct {
	// Tools is the list of this toolchain's tools and their definitions.
	Tools []*ToolInfo
}
