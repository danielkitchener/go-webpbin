package webpbin

import (
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/image/webp"
)

func init() {
	DetectUnsupportedPlatforms()
	downloadFile("https://upload.wikimedia.org/wikipedia/commons/e/e3/Avola-Syracuse-Sicilia-Italy_-_Creative_Commons_by_gnuckx_%283858115914%29.jpg", "source.jpg")
	downloadFile("https://upload.wikimedia.org/wikipedia/commons/d/d1/Snail_in_Forest_on_Lichtscheid_2.webp", "source.webp")
}

func downloadFile(url, target string) {
	_, err := os.Stat(target)

	if err != nil {
		client := &http.Client{
			
		}
		req, err := http.NewRequest("GET", url, nil)
		
		if err != nil {
			fmt.Printf("Error while creating request to download test image: %v\n", err)
			panic(err)
		}
		req.Header.Add("User-Agent", "CWebp-Go-Binwrapper-Test/1.0")

		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf("Error while downloading test image: %v\n", err)
			panic(err)
		}

		defer resp.Body.Close()

		f, err := os.Create(target)

		if err != nil {
			panic(err)
		}

		defer f.Close()

		_, err = io.Copy(f, resp.Body)

		if err != nil {
			panic(err)
		}
	}
}

func TestEncodeImage(t *testing.T) {
	c := NewCWebP()
	f, err := os.Open("source.jpg")
	assert.Nil(t, err)
	img, err := jpeg.Decode(f)
	assert.Nil(t, err)
	c.InputImage(img)
	c.OutputFile("target.webp")
	err = c.Run()
	assert.Nil(t, err)
	validateWebp(t)
}

func TestEncodeReader(t *testing.T) {
	c := NewCWebP()
	f, err := os.Open("source.jpg")
	assert.Nil(t, err)
	c.Input(f)
	c.OutputFile("target.webp")
	err = c.Run()
	assert.Nil(t, err)
	validateWebp(t)
}

func TestEncodeFile(t *testing.T) {
	c := NewCWebP()
	c.InputFile("source.jpg")
	c.OutputFile("target.webp")
	err := c.Run()
	assert.Nil(t, err)
	validateWebp(t)
}

func TestEncodeWriter(t *testing.T) {
	f, err := os.Create("target.webp")
	assert.Nil(t, err)
	defer f.Close()

	c := NewCWebP()
	c.InputFile("source.jpg")
	c.Output(f)
	err = c.Run()
	assert.Nil(t, err)
	f.Close()
	validateWebp(t)
}

func TestVersionCWebP(t *testing.T) {
	c := NewCWebP()
	r, err := c.Version()
	assert.Nil(t, err)

	if _, ok := os.LookupEnv("DOCKER_ARM_TEST"); !ok {
		assert.Contains(t, r, "1.4.0")
	}
}

func validateWebp(t *testing.T) {
	defer os.Remove("target.webp")
	fSource, err := os.Open("source.jpg")
	assert.Nil(t, err)
	imgSource, err := jpeg.Decode(fSource)
	assert.Nil(t, err)
	fTarget, err := os.Open("target.webp")
	assert.Nil(t, err)
	defer fTarget.Close()
	imgTarget, err := webp.Decode(fTarget)
	assert.Nil(t, err)
	assert.Equal(t, imgSource.Bounds(), imgTarget.Bounds())
}
