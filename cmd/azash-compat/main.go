package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yuanying/azash/internal/compat/pipeline"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	sub := os.Args[1]
	switch sub {
	case "golden-java", "compare", "run":
		runSub(sub, os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", sub)
		usage()
		os.Exit(2)
	}
}

func runSub(sub string, args []string) {
	opts := pipeline.DefaultOptions()
	var sampleIDs csvList

	fs := flag.NewFlagSet(sub, flag.ExitOnError)
	fs.StringVar(&opts.SamplesFile, "samples-file", opts.SamplesFile, "path to samples csv")
	fs.StringVar(&opts.InputDir, "input-dir", opts.InputDir, "input directory")
	fs.StringVar(&opts.GoldenDir, "golden-dir", opts.GoldenDir, "java golden output directory")
	fs.StringVar(&opts.GoOutDir, "go-out-dir", opts.GoOutDir, "go epub output directory")
	fs.StringVar(&opts.WorkDir, "work-dir", opts.WorkDir, "work directory for unpacked/normalized files")
	fs.StringVar(&opts.ReportDir, "report-dir", opts.ReportDir, "report output directory")
	fs.StringVar(&opts.JavaCmd, "java-cmd", opts.JavaCmd, "java executable")
	fs.StringVar(&opts.JavaCP, "java-cp", opts.JavaCP, "java classpath")
	fs.StringVar(&opts.JavaMain, "java-main", opts.JavaMain, "java main class")
	fs.StringVar(&opts.INIPath, "ini", opts.INIPath, "AozoraEpub3 ini path")
	fs.Var(&sampleIDs, "sample-id", "sample id filter (repeatable)")
	fs.BoolVar(&opts.FailOnDiff, "fail-on-diff", false, "return non-zero if diff is found")
	_ = fs.Parse(args)

	if len(sampleIDs) > 0 {
		opts.SampleIDs = map[string]struct{}{}
		for _, id := range sampleIDs {
			opts.SampleIDs[id] = struct{}{}
		}
	}
	opts.LogDir = opts.ReportDir + "/logs"

	mode := pipeline.Mode(sub)
	r, err := pipeline.Execute(context.Background(), mode, opts)
	fmt.Printf("mode=%s samples=%d files_changed=%d report=%s\n", sub, r.SampleCount, r.DiffSummary.FilesChanged, opts.ReportDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("azash-compat <golden-java|compare|run> [flags]")
	fmt.Println("  --sample-id can be repeated, e.g. --sample-id sample-001 --sample-id sample-002")
}

type csvList []string

func (c *csvList) String() string {
	return strings.Join(*c, ",")
}

func (c *csvList) Set(v string) error {
	*c = append(*c, strings.TrimSpace(v))
	return nil
}
