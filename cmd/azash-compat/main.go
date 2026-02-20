package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuanying/azash/internal/compat/pipeline"
	"github.com/yuanying/azash/internal/resource"
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
	case "resource-diagnose":
		runResourceDiagnose(os.Args[2:])
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
	var cliBase string
	var bundleDir string

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
	fs.StringVar(&cliBase, "resource-base", "", "resource base directory (highest priority)")
	fs.StringVar(&bundleDir, "bundle-dir", "", "bundled resource directory (defaults to executable dir)")
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
	logRuntimeResourceResolution(cliBase, bundleDir)

	mode := pipeline.Mode(sub)
	r, err := pipeline.Execute(context.Background(), mode, opts)
	fmt.Printf("mode=%s samples=%d files_changed=%d report=%s\n", sub, r.SampleCount, r.DiffSummary.FilesChanged, opts.ReportDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runResourceDiagnose(args []string) {
	var cliBase string
	var bundleDir string
	var webSite string

	fs := flag.NewFlagSet("resource-diagnose", flag.ExitOnError)
	fs.StringVar(&cliBase, "resource-base", "", "resource base directory (highest priority)")
	fs.StringVar(&bundleDir, "bundle-dir", "", "bundled resource directory (defaults to executable dir)")
	fs.StringVar(&webSite, "web-site", "", "site name for web/<site>/extract.txt check")
	_ = fs.Parse(args)

	resolver, err := resource.New(resource.Options{
		CLIBaseDir: cliBase,
		BundleDir:  bundleDir,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: init resolver: %v\n", err)
		os.Exit(1)
	}

	printCoreResourceResolution(resolver)
	printWebExtractResolution(resolver, webSite)
}

func usage() {
	fmt.Println("azash-compat <golden-java|compare|run|resource-diagnose> [flags]")
	fmt.Println("  --sample-id can be repeated, e.g. --sample-id sample-001 --sample-id sample-002")
	fmt.Println("  --resource-base/--bundle-dir are available on all subcommands")
	fmt.Println("  resource-diagnose: print resolved paths for template/chuki/replace*/gaiji/presets/web resources")
}

type csvList []string

func (c *csvList) String() string {
	return strings.Join(*c, ",")
}

func (c *csvList) Set(v string) error {
	*c = append(*c, strings.TrimSpace(v))
	return nil
}

func logRuntimeResourceResolution(cliBase, bundleDir string) {
	resolver, err := resource.New(resource.Options{
		CLIBaseDir: cliBase,
		BundleDir:  bundleDir,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "resource\tinit-error\t%v\n", err)
		return
	}
	fmt.Fprintln(os.Stderr, "resource\tbegin")
	printCoreResourceResolutionTo(resolver, os.Stderr)
	fmt.Fprintln(os.Stderr, "resource\tend")
}

func printCoreResourceResolution(resolver *resource.Resolver) {
	printCoreResourceResolutionTo(resolver, os.Stdout)
}

func printCoreResourceResolutionTo(resolver *resource.Resolver, out io.Writer) {
	checks := []struct {
		path string
		dir  bool
	}{
		{path: "template", dir: true},
		{path: "gaiji", dir: true},
		{path: "presets", dir: true},
		{path: "chuki_tag.txt"},
		{path: "chuki_tag_suf.txt"},
		{path: "chuki_utf.txt"},
		{path: "chuki_ivs.txt"},
		{path: "chuki_alt.txt"},
		{path: "chuki_latin.txt"},
	}

	for _, chk := range checks {
		if chk.dir {
			res, err := resolver.ResolveDir(chk.path)
			if err != nil {
				_, _ = fmt.Fprintf(out, "%s\tmissing\t%v\n", chk.path, err)
				continue
			}
			_, _ = fmt.Fprintf(out, "%s\t%s\t%s\n", chk.path, res.Source, res.Path)
			continue
		}
		res, err := resolver.ResolveFile(chk.path)
		if err != nil {
			_, _ = fmt.Fprintf(out, "%s\tmissing\t%v\n", chk.path, err)
			continue
		}
		_, _ = fmt.Fprintf(out, "%s\t%s\t%s\n", chk.path, res.Source, res.Path)
	}

	replaces, err := resolver.ResolveFileGlob("replace*.txt")
	if err != nil {
		_, _ = fmt.Fprintf(out, "replace*.txt\tmissing\t%v\n", err)
		return
	}
	if len(replaces) == 0 {
		_, _ = fmt.Fprintln(out, "replace*.txt\tnot-found\t(optional)")
		return
	}
	sort.Slice(replaces, func(i, j int) bool { return replaces[i].Resource < replaces[j].Resource })
	for _, res := range replaces {
		_, _ = fmt.Fprintf(out, "%s\t%s\t%s\n", res.Resource, res.Source, res.Path)
	}
}

func printWebExtractResolution(resolver *resource.Resolver, webSite string) {
	if webSite != "" {
		printWebExtractForSite(resolver, webSite, os.Stdout)
		return
	}

	webDir, ok, err := resolver.ResolveOptionalDir("web")
	if err != nil {
		fmt.Printf("web\tmissing\t%v\n", err)
		return
	}
	if !ok {
		fmt.Println("web\tnot-found\t(optional)")
		return
	}
	fmt.Printf("web\t%s\t%s\n", webDir.Source, webDir.Path)

	entries, _, err := resolver.ReadDir("web")
	if err != nil {
		fmt.Printf("web/*/extract.txt\tmissing\t%v\n", err)
		return
	}
	siteDirs := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			siteDirs = append(siteDirs, e.Name())
		}
	}
	if len(siteDirs) == 0 {
		fmt.Println("web/*/extract.txt\tnot-found\t(no site directories)")
		return
	}
	sort.Strings(siteDirs)
	for _, site := range siteDirs {
		printWebExtractForSite(resolver, site, os.Stdout)
	}
}

func printWebExtractForSite(resolver *resource.Resolver, site string, out io.Writer) {
	webExtract := filepath.ToSlash(filepath.Join("web", site, "extract.txt"))
	res, err := resolver.ResolveFile(webExtract)
	if err != nil {
		_, _ = fmt.Fprintf(out, "%s\tmissing\t%v\n", webExtract, err)
		return
	}
	_, _ = fmt.Fprintf(out, "%s\t%s\t%s\n", webExtract, res.Source, res.Path)
}
