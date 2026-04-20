package cmd

// Version is set at build time via ldflags.
// Falls back to "dev" when built without GoReleaser.
var Version = "dev"
