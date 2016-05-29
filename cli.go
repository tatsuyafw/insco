package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
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

func (cli *CLI) showHelp() {
	// TODO: show usage
	fmt.Fprintf(cli.errStream, "Usage: \n")
}

// [WIP]
func downloadFile(url, tempDir string) (filePath string, err error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	downloadedFilePath := filepath.Join(tempDir, fileName)
	fmt.Println("Downloading", url, "to", downloadedFilePath)

	file, err := os.Create(downloadedFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, err)
		return "", err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, err)
		return "", err
	}

	fmt.Println("Downloaded:", fileName)
	return downloadedFilePath, nil
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

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func gunzip(src, dir string) (dst string, err error) {
	reader, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	dst = filepath.Join(dir, gzipReader.Name)
	writer, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer writer.Close()

	_, err = io.Copy(writer, gzipReader)
	if err != nil {
		return "", err
	}

	return dst, nil
}

func untar(tarball, dir string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		path := filepath.Join(dir, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err := os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
		file.Close()
	}
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

func (cli *CLI) emacs(version string) error {
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
	flags := "--without-x"

	dir, err := ioutil.TempDir(os.TempDir(), Name)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer os.RemoveAll(dir)

	// Download an archive file
	filePath, err := downloadFile(mirrorListUrl+"/"+archFile, dir)
	if err != nil {
		return err
	}

	// Unzip and Untar
	tarball, err := gunzip(filePath, dir)
	if err != nil {
		fmt.Fprintln(cli.errStream, err)
		return err
	}
	err = untar(tarball, dir)
	if err != nil {
		fmt.Fprintln(cli.errStream, err)
		return err
	}
	ext := filepath.Ext(tarball)
	contentDir := strings.TrimRight(tarball, ext)

	err = os.Chdir(contentDir)
	if err != nil {
		fmt.Fprintln(cli.errStream, err)
	}
	current, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(cli.errStream, err)
	}
	fmt.Fprintln(cli.outStream, current)

	// Build
	fmt.Fprintln(cli.outStream, "Building...")
	out, err := exec.Command("./configure", "--prefix="+prefixDir()+" "+flags).Output()
	if err != nil {
		fmt.Fprintln(cli.errStream, err)
	}
	fmt.Fprintln(cli.outStream, string(out))

	out, err = exec.Command("make").Output()
	if err != nil {
		fmt.Fprintln(cli.errStream, err)
	}
	fmt.Fprintln(cli.outStream, string(out))

	out, err = exec.Command("make", "install").Output()
	if err != nil {
		fmt.Fprintln(cli.errStream, err)
	}
	fmt.Fprintln(cli.outStream, string(out))

	fmt.Fprintln(cli.outStream, "Finished.")

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
		cli.emacs("") // TODO: pass version if given as a cli argument
	default:
		// TODO: show usage
	}

	return ExitCodeOK
}
