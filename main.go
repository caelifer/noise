package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"

	ico "github.com/biessek/golang-ico"
	"github.com/mazznoer/colorgrad"
	"github.com/ojrac/opensimplex-go"
)

const (
	// IconHgt favicon height
	IconHgt = 32
	// IconWdt favicon width
	IconWdt = 32
	// ImageHgt generated image y-dimension
	ImageHgt = 480
	// ImageWdt generated image x-dimension
	ImageWdt = 600

	// EnvHTTPPortVar environment variable name for HTTP port parameter
	EnvHTTPPortVar = "APP_HTTP_PORT"
	// DefaultHTTPPort default HTTP port
	DefaultHTTPPort = "8080"
)

var (
	grad         = colorgrad.Rainbow().Sharp(7, 0)
	debugEnabled = len(os.Getenv("DEBUG")) != 0
	favIcon      = createFavIcon()
)

func main() {
	httpPort, ok := os.LookupEnv(EnvHTTPPortVar)
	if !ok {
		httpPort = DefaultHTTPPort
	}

	// Basic routing table

	// FavIcon GET request
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		if _, err := w.Write(favIcon); err != nil {
			log.Printf("failed to write 'favicon.ico' to HTTP stream: %v", err)
		}
	})

	// Main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Generate colored noise image
		w.Header().Set("Content-Type", "image/png")
		if err := writePNGImage(w, createImage(ImageWdt, ImageHgt)); err != nil {
			log.Fatal(err)
		}
	})

	// Start HTTP server
	addr := ":" + httpPort
	lsn, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalf("failed to open listening socket: %v", err)
	}
	log.Printf("Starting HTTP server: %q", fmt.Sprintf("http://%s/", lsn.Addr()))
	log.Fatal(http.Serve(lsn, nil))
}

func writePNGImage(w http.ResponseWriter, img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Printf("failed to encode PNG image: %v", err)
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, "Failed to encode image: %v", err)
		return err
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf.Bytes())))
	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Printf("failed to write image to HTTP stream: %v", err)
		return err
	}
	return nil
}

func createImage(wdt, hgt int) image.Image {
	const scale = 0.02
	img := image.NewRGBA(image.Rect(0, 0, wdt, hgt))
	noise := opensimplex.NewNormalized(rand.Int63n(1000))

	for y := 0; y < hgt; y++ {
		for x := 0; x < wdt; x++ {
			t := noise.Eval2(float64(x)*scale, float64(y)*scale)
			img.Set(x, y, grad.At(t))
		}
	}

	if debugEnabled {
		log.Printf("DEBUG: generated new noise image [%dx%d]", wdt, hgt)
	}
	return img

}

func createFavIcon() []byte {
	var buf bytes.Buffer
	if err := ico.Encode(&buf, createImage(IconWdt, IconHgt)); err != nil {
		log.Fatalf("failed to generate 'favicon.ico': %v", err)
	}
	return buf.Bytes()
}
