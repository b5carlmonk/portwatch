package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/portwatch/portwatch/internal/monitor"
)

const version = "0.1.0"

func main() {
	// CLI flags
	interval := flag.Int("interval", 30, "Polling interval in seconds")
	configFile := flag.String("config", "", "Path to config file (optional)")
	showVersion := flag.Bool("version", false, "Print version and exit")
	once := flag.Bool("once", false, "Run a single scan and exit")
	verbose := flag.Bool("verbose", false, "Enable verbose output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "portwatch v%s - Monitor open ports and detect service changes\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: portwatch [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  portwatch --once\n")
		fmt.Fprintf(os.Stderr, "  portwatch --interval 60 --verbose\n")
		fmt.Fprintf(os.Stderr, "  portwatch --config /etc/portwatch/config.yaml\n")
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("portwatch v%s\n", version)
		os.Exit(0)
	}

	// Build monitor options from flags
	opts := monitor.Options{
		Interval:   *interval,
		ConfigFile: *configFile,
		Verbose:    *verbose,
		RunOnce:    *once,
	}

	m, err := monitor.New(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initializing monitor: %v\n", err)
		os.Exit(1)
	}

	if err := m.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running monitor: %v\n", err)
		os.Exit(1)
	}
}
