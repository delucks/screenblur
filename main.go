package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	stackblur "github.com/esimov/stackblur-go"
	"github.com/kbinani/screenshot"
)

func main() {
	var (
		fileName, cpuProf, memProf string
		blurRadius, displayIndex   uint
	)
	flag.StringVar(&fileName, "out", "/tmp/screenshotblur.png", "output `file` for blurred screenshot")
	flag.UintVar(&blurRadius, "radius", 15, "`radius` of blur to apply")
	flag.UintVar(&displayIndex, "display", 0, "display `index` to screenshot")
	// Gotta go fast
	flag.StringVar(&cpuProf, "cpuprofile", "", "write CPU profile to `file`")
	flag.StringVar(&memProf, "memprofile", "", "write memory profile to `file`")
	flag.Parse()

	if cpuProf != "" {
		f, err := os.Create(cpuProf)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	bounds := screenshot.GetDisplayBounds(int(displayIndex))
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		log.Fatalf("could not take screenshot in bounds %#v of index %d: %v\n", bounds, displayIndex, err)
	}
	// Converting to RGBA takes ~70% of runtime
	blurred := stackblur.Process(img, uint32(blurRadius))
	file, _ := os.Create(fileName)
	defer file.Close()
	// Compressing the PNG is the most expensive computation performed by the naive version
	encoder := png.Encoder{CompressionLevel: png.NoCompression}
	encoder.Encode(file, blurred)
	fmt.Printf("Blurred screenshot saved to %s\n", fileName)

	if memProf != "" {
		f, err := os.Create(memProf)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
