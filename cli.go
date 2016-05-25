package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	ExitCodeOK = iota
	ExitCodeParserFlagError
)

const (
	HomeEnv = "HOME"
)

type Options struct {
	OptHelp    bool `short:"h" long:"help" description:"Show this help message and exit"`
	OptVersion bool `short:"v" long:"version" description:"Print the version and exit"`
}

type CLI struct {
	outStream, errStream io.Writer
}

// [WIP]
func downloadFile(url, tempDir string) (downloadedFile string, err error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	downloadedFilePath := filepath.Join(tempDir, fileName)
	fmt.Println("Downloading", url, "to", downloadedFilePath)

	file, err := os.Create(downloadedFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, err)
		return nil, err
	}
	defer response.Body.Close()

	n, err := io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, err)
		return nil, err
	}

	fmt.Println("Downloaded:", fileName)
	return file, nil
}

func (cli *CLI) showHelp() {
	// TODO: show usage
	fmt.Fprintf(cli.errStream, "Usage: \n")
}

func homeDir() string {
	homeDir := os.Getenv(HomeEnv)
	// TODO: return an error if HOME environment variablie is not set
	return homeDir
}

func binaryDir() string {
	return homeDir() + "/bin"
}

func prefixDir() string {
	return homeDir() + "/usr/local"
}

func setup() error {
	binaryDir := binaryDir()

	if err := os.MkdirAll(binaryDir, os.ModePerm); err != nil {
		return err
	}

	prefixDir := prefixDir()
	if err := os.MkdirAll(prefixDir, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func emacs(version string) error {
	const (
		DEFAULT_VERSION = "24.5"
	)
	if version == "" {
		version = DEFAULT_VERSION
	}

	content := "emacs-" + version
	// comment out temporarily to pass go compilation
	archFile := content + ".tar.gz"
	mirrorListUrl := "http://ftpmirror.gnu.org/emacs"
	// flags := "--without-x"

	dir, err := ioutil.TempDir(os.TempDir(), Name)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer os.RemoveAll(dir)

	// Download an archive file
	file, err := downloadFile(mirrorListUrl+"/"+archFile, dir)
	if err != nil {
		return err
	}

	os.Chdir(dir)

	// Build
	// TODO

	return nil
}

func (cli *CLI) parseOptions(args []string) (*Options, []string, error) {
	opts := &Options{}
	p := flags.NewParser(opts, flags.PrintErrors)
	args, err := p.ParseArgs(args)
	if err != nil {
		return nil, nil, err
	}

	return opts, args, nil
}

func (cli *CLI) Run(args []string) int {
	opts, parsedArgs, err := cli.parseOptions(args)
	if err != nil {
		return ExitCodeParserFlagError
	}

	if opts.OptHelp {
		cli.showHelp()
		return ExitCodeOK
	}

	if opts.OptVersion {
		fmt.Fprintf(cli.errStream, "%s: v%s\n", Name, Version)
		return ExitCodeOK
	}

	if len(parsedArgs) == 0 {
		fmt.Fprintf(cli.errStream, "[Error]: You must specify the target.\n")
		cli.showHelp()
		return ExitCodeParserFlagError
	}

	target := parsedArgs[0]
	// version := parsedArgs[1]

	setup()

	switch target {
	case "emacs":
		emacs("") // TODO: pass version if given as a cli argument
	default:
		// TODO: show usage
	}

	return ExitCodeOK
}
