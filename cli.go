package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
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

func homeDir() string {
	homeDir := os.Getenv(HomeEnv)
	// TODO: return an error if HOME environment variablie is not set
	return homeDir
}

func binaryDir() string {
	return homeDir() + "/bin"
}

func basePrefixDir() string {
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

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func setup() error {
	binaryDir := binaryDir()

	if err := os.MkdirAll(binaryDir, os.ModePerm); err != nil {
		return err
	}

	prefixDir := basePrefixDir()
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
	archFile := content + ".tar.gz"
	mirrorListUrl := "http://ftpmirror.gnu.org/emacs"
	flags := []string{"--without-x", "--with-gnutls"}

	dir, err := ioutil.TempDir(os.TempDir(), Name)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer os.RemoveAll(dir)

	// Download an archive file
	filePath, err := DownloadFile(mirrorListUrl+"/"+archFile, dir)
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
		return err
	}

	// Build
	fmt.Fprintln(cli.outStream, "Building...")
	runner := Runner{outStream: cli.outStream}

	prefixDir := filepath.Join(basePrefixDir(), content)
	flags = append([]string{"--prefix=" + prefixDir}, flags...)

	runner.Run(exec.Command("./configure", flags...))
	runner.Run(exec.Command("make"))
	runner.Run(exec.Command("make", "install"))

	if err = runner.Err(); err != nil {
		fmt.Fprintln(cli.errStream, err)
		return err
	}

	// Create a symbolic link
	originalBinary := filepath.Join(prefixDir, "bin", "emacs")
	binaryLink := filepath.Join(binaryDir(), "emacs")
	if exists(binaryLink) {
		os.Rename(binaryLink, binaryLink+".org")
	}
	if err := os.Symlink(originalBinary, binaryLink); err != nil {
		fmt.Fprintln(cli.errStream, err)
	}

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
