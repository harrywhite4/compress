package main

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/spf13/pflag"
	"image"
	"image/jpeg"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Import png for initialisation
import _ "image/png"

const (
	programName        string = "compress"
	programDescription string = "Compress your images"
	programVersion     string = "1.0.1"
)

var (
	allowedExtensions = [3]string{".jpg", ".jpeg", ".png"}
)

type arguments struct {
	disableResize bool
	halfResize    bool
	quality       int
	suffix        string
	targetPixels  int
	verbose       bool
	files         []string
}

func getNewFilename(path string, suffix string) string {
	dir, filename := filepath.Split(path)
	splitname := strings.Split(filename, ".")
	return dir + splitname[0] + suffix + ".jpg"
}

func resizeImage(initImage *image.Image, args *arguments) {
	performResize := false
	curSize := (*initImage).Bounds().Size()
	width := curSize.X
	height := curSize.Y
	currentPixels := width * height
	var newWidth, newHeight int

	if args.verbose {
		fmt.Println("Resizing...")
	}

	if args.halfResize {
		// If resizing by half
		newHeight = height / 2
		newWidth = width / 2
		performResize = true
	} else if currentPixels > args.targetPixels {
		// If not resizing by half and above target pixels
		ratio := float64(height) / float64(width)
		if args.verbose {
			fmt.Printf("Ratio: %v\n", ratio)
		}
		newWidth = int(math.Sqrt(float64(args.targetPixels) / ratio))
		newHeight = int(float64(newWidth) * ratio)
		performResize = true
	}

	if performResize {
		if args.verbose {
			fmt.Printf("Width: %v Height: %v\n", newWidth, newHeight)
		}
		*initImage = imaging.Resize(*initImage, newWidth, newHeight, imaging.MitchellNetravali)
	} else if args.verbose {
		fmt.Println("No resize required")
	}
}

func checkFileExtension(file string) bool {
	extension := filepath.Ext(file)
	for _, allowedExt := range allowedExtensions {
		if extension == allowedExt {
			return true
		}
	}
	return false
}

func processFile(file string, args *arguments) error {
	fmt.Printf("Processing %v\n", file)

	allowed := checkFileExtension(file)

	if !allowed {
		if args.verbose {
			fmt.Printf("Skipping %v not valid file extension\n", file)
		}
		return nil
	}
	comp_file := getNewFilename(file, args.suffix)

	reader, err := os.Open(file)
	if err != nil {
		return err
	}

	initImage, format, err := image.Decode(reader)
	if err != nil {
		return err
	} else if args.verbose {
		fmt.Println("Decoded " + format)
	}

	if !args.disableResize {
		resizeImage(&initImage, args)
	}

	writer, err := os.Create(comp_file)
	if err != nil {
		return err
	}

	options := jpeg.Options{Quality: args.quality}

	err = jpeg.Encode(writer, initImage, &options)
	if err != nil {
		return err
	} else if args.verbose {
		fmt.Printf("Saved %v with %v%% quality\n", comp_file, args.quality)
	}
	return nil
}

// Given a list of files or directories return files + files in directories
func getFiles(items []string) []string {
	var files []string
	for _, file := range items {
		info, err := os.Stat(file)
		if err != nil {
			fmt.Printf("Could not inspect %v\n", file)
			continue
		}

		if info.IsDir() {
			fmt.Printf("Processing files in directory %v\n", file)
			dir := file
			// Get all files in directory
			fileInfos, err := ioutil.ReadDir(dir)
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, info := range fileInfos {
				if !info.IsDir() {
					fullPath := filepath.Join(dir, info.Name())
					allowed := checkFileExtension(fullPath)
					if allowed {
						files = append(files, fullPath)
					}
				}
			}
		} else {
			allowed := checkFileExtension(file)
			if allowed {
				files = append(files, file)
			} else {
				fmt.Printf("Ignoring %v invalid file extension\n", file)
			}
		}
	}
	return files
}

func validateArgs(args *arguments) error {
	if args.quality < 0 || args.quality > 100 {
		return errors.New("Quality must be between 0 and 100")
	}

	if args.targetPixels < 0 {
		return errors.New("Target pixels can not be negative")
	}

	return nil
}

func printHelp(flagSet *pflag.FlagSet) {
	usage := `Usage: %v [OPTION]... [PATH]...

%v

Path can be an image file or a directory.
If it is a directory, all images within that directory will be processed.

Options:
%v
`
	options := flagSet.FlagUsages()
	fmt.Printf(usage, programName, programDescription, options)
}

// Parse command line arguments
func parseArgs() *arguments {
	args := new(arguments)
	flagSet := pflag.NewFlagSet("compress", pflag.ExitOnError)
	flagSet.BoolVarP(&args.disableResize, "no-resize", "n", false, "Keep image at original size")
	flagSet.BoolVarP(&args.halfResize, "half", "2", false, "Save image at half it's original size")
	flagSet.IntVarP(&args.quality, "quality", "q", 80, "Quality to save image at 0-100")
	flagSet.IntVarP(&args.targetPixels, "pixels", "p", 2073600, "Target pixel count for resized image")
	flagSet.StringVarP(&args.suffix, "suffix", "s", "_compressed", "Suffix to be appended to filenames")
	flagSet.BoolVarP(&args.verbose, "verbose", "v", false, "Print additional information during processing")
	help := flagSet.BoolP("help", "h", false, "Display help")
	version := flagSet.Bool("version", false, "Show version")

	flagSet.Parse(os.Args[1:])
	// Handle help and version before any further processing
	if *help {
		printHelp(flagSet)
		os.Exit(1)
	}
	if *version {
		fmt.Printf("%v %v\n", programName, programVersion)
		os.Exit(1)
	}

	// Process positional arguments
	positionals := flagSet.Args()
	if len(positionals) == 0 {
		fmt.Println("No path specified")
		fmt.Println("Run compress --help for usage")
		os.Exit(1)
	}
	args.files = getFiles(positionals)
	return args
}

func main() {
	args := parseArgs()
	err := validateArgs(args)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	// Process files
	for _, file := range args.files {
		err := processFile(file, args)
		if err != nil {
			fmt.Printf("%v failed with error %v\n", file, err)
		}
	}

	fmt.Println("Done!")
}
