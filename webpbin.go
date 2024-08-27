package webpbin

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/belphemur/go-binwrapper"
)

var skipDownload bool
var dest = GetPath()
var libwebpVersion = "1.4.0"

func GetPath() string {
	return filepath.Join(map[string]string{
		"windows": filepath.Join(os.Getenv("APPDATA")),
		"darwin":  filepath.Join(os.Getenv("HOME"), ".cache"),
		"linux":   filepath.Join(os.Getenv("HOME"), ".cache"),
	}[runtime.GOOS], "webp", libwebpVersion, "bin")
}

type OptionFunc func(binWrapper *binwrapper.BinWrapper) error

func SetSkipDownload(isSkipDownload bool) OptionFunc {
	return func(binWrapper *binwrapper.BinWrapper) error {
		skipDownload = isSkipDownload
		return nil
	}
}

func SetVendorPath(path string) OptionFunc {
	return func(binWrapper *binwrapper.BinWrapper) error {
		dest = path
		return nil
	}
}

func SetLibVersion(version string) OptionFunc {
	return func(binWrapper *binwrapper.BinWrapper) error {
		libwebpVersion = version
		dest = GetPath()
		return nil
	}

}

func loadDefaultFromENV(binWrapper *binwrapper.BinWrapper) error {
	if os.Getenv("SKIP_DOWNLOAD") == "true" {
		skipDownload = true
	}

	if path := os.Getenv("VENDOR_PATH"); path != "" {
		dest = path
	}

	if version := os.Getenv("LIBWEBP_VERSION"); version != "" {
		libwebpVersion = version
	}

	return nil
}

// DetectUnsupportedPlatforms detects platforms without prebuilt binaries (alpine and arm).
// For this platforms libwebp tools should be built manually.
// See https://github.com/belphemur/go-webpbin/blob/master/docker/Dockerfile and https://github.com/belphemur/go-webpbin/blob/master/docker/Dockerfile.arm for details
func DetectUnsupportedPlatforms() {
	if runtime.GOARCH == "arm" {
		skipDownload = true
	} else if runtime.GOOS == "linux" {
		output, err := os.ReadFile("/etc/issue")

		if err == nil && bytes.Contains(bytes.ToLower(output), []byte("alpine")) {
			skipDownload = true
		}
	}
}

func createBinWrapper(optionFuncs ...OptionFunc) *binwrapper.BinWrapper {
	base := "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/"

	b := binwrapper.NewBinWrapper().AutoExe()

	loadDefaultFromENV(b)

	for _, optionFunc := range optionFuncs {
		optionFunc(b)
	}

	if !skipDownload {
		b.Src(
			binwrapper.NewSrc().
				URL(base + "libwebp-" + libwebpVersion + "-mac-arm64.tar.gz").
				Os("darwin").
				Arch("arm64")).
			Src(
				binwrapper.NewSrc().
					URL(base + "libwebp-" + libwebpVersion + "-mac-x86-64.tar.gz").
					Os("darwin").
					Arch("x64")).
			Src(
				binwrapper.NewSrc().
					URL(base + "libwebp-" + libwebpVersion + "-linux-x86-32.tar.gz").
					Os("linux").
					Arch("x86")).
			Src(
				binwrapper.NewSrc().
					URL(base + "libwebp-" + libwebpVersion + "-linux-x86-64.tar.gz").
					Os("linux").
					Arch("x64")).
			Src(
				binwrapper.NewSrc().
					URL(base + "libwebp-" + libwebpVersion + "-windows-x64.zip").
					Os("win32").
					Arch("x64")).
			Src(
				binwrapper.NewSrc().
					URL(base + "libwebp-" + libwebpVersion + "-windows-x86.zip").
					Os("win32").
					Arch("x86"))
	}

	return b.Strip(2).Dest(dest)
}

func createReaderFromImage(img image.Image) (io.Reader, error) {
	enc := &png.Encoder{
		CompressionLevel: png.NoCompression,
	}

	var buffer bytes.Buffer
	err := enc.Encode(&buffer, img)

	if err != nil {
		return nil, err
	}

	return &buffer, nil
}

func version(b *binwrapper.BinWrapper) (string, error) {
	b.Reset()
	err := b.Run("-version")

	if err != nil {
		return "", err
	}

	version := string(b.StdOut())
	version = strings.Replace(version, "\n", "", -1)
	version = strings.Replace(version, "\r", "", -1)
	return version, nil
}
