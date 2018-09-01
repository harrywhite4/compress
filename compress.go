package main

import (
  "os"
  "io/ioutil"
  "fmt"
  "strings"
  "image"
  "image/jpeg"
  "path/filepath"
  "math"
  "github.com/integrii/flaggy"
  "github.com/disintegration/imaging"
)
// Import png for initialisation
import _ "image/png"

var (
  file string
  disableResize = false
  halfResize = false
  quality = 80
  suffix = "_compressed"
  targetPixels = 2073600
  files []string
  verbose = false
  allowedExtensions = [3]string{".jpg", ".jpeg", ".png"}
)

func init() {
  flaggy.SetName("compress")
  flaggy.SetDescription("Compress your images")
  flaggy.SetVersion("0.1")

  flaggy.AddPositionalValue(&file, "file", 1, true, "Image file to compress, if passed a directory all images in the directory will be processed")
  flaggy.Bool(&disableResize, "n", "no-resize", "Keep image at original size")
  flaggy.Bool(&halfResize, "2", "half", "Save image at half it's original size")
  flaggy.Int(&quality, "q", "quality", "Quality to save image at 0-100")
  flaggy.Int(&targetPixels, "p", "pixels", "Target pixel count for resized image")
  flaggy.String(&suffix, "s", "suffix", "Suffix to be appended to filenames")
  flaggy.Bool(&verbose, "v", "verbose", "Print additional information during processing")
  flaggy.Parse()
}

func getNewFilename(path string) string {
  dir, filename := filepath.Split(path)
  splitname := strings.Split(filename, ".")
  return dir + splitname[0] + suffix + ".jpg"
}

func resizeImage(initImage *image.Image) {
  performResize := false
  curSize := (*initImage).Bounds().Size()
  width := curSize.X
  height := curSize.Y
  currentPixels := width*height
  var newWidth, newHeight int

  if verbose {
    fmt.Println("Resizing...")
  }

  if halfResize {
    // If resizing by half
    newHeight = height / 2
    newWidth = width / 2
    performResize = true
  } else if currentPixels > targetPixels {
    // If not resizing by half and above target pixels
    ratio := float64(height) / float64(width)
    if verbose {
      fmt.Printf("Ratio: %v\n", ratio)
    }
    newWidth = int(math.Sqrt(float64(targetPixels)/ratio))
    newHeight = int(float64(newWidth)*ratio)
    performResize = true
  }

  if performResize {
    if verbose {
      fmt.Printf("Width: %v Height: %v\n", newWidth, newHeight)
    }
    *initImage = imaging.Resize(*initImage, newWidth, newHeight, imaging.MitchellNetravali)
  } else if verbose {
    fmt.Println("No resize required")
  }
}

func getFiles(file string) []string {
  var files []string
  info, err := os.Lstat(file)
  if err != nil {
    fmt.Printf("Could not inspect %v\n", file)
    fmt.Println(err)
    os.Exit(1)
  }
  
  if info.IsDir() {
    fmt.Printf("Processing files in directory %v\n", file)
    dir := file
    // Get all files in directory
    fileInfos, err := ioutil.ReadDir(dir)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    files = make([]string, 0, len(fileInfos))
    for _, info := range fileInfos {
      if !info.IsDir() {
        fullPath := filepath.Join(dir, info.Name())
        files = append(files, fullPath)
      }
    }
  } else {
    files = make([]string, 1)
    files[0] = file
  }
  return files
}

func checkFileExtension(file string) bool {
  extension := filepath.Ext(file)
  allowed :=  false
  for _, allowedExt := range allowedExtensions {
    if extension == allowedExt {
      allowed = true
      break
    }
  }
  return allowed
}

func processFile(file string) bool {
  fmt.Printf("Processing %v\n", file)

  allowed := checkFileExtension(file)

  if !allowed {
    if verbose {
      fmt.Printf("Skipping %v not valid file extension\n", file)
    }
    return true
  }
  comp_file := getNewFilename(file)

  reader, err := os.Open(file)
  if err != nil {
    fmt.Println(err)
    return false
  }

  initImage, format, err := image.Decode(reader)
  if err != nil {
    fmt.Println(err)
    return false
  } else if verbose {
    fmt.Println("Decoded " + format)
  }

  if !disableResize {
    resizeImage(&initImage)
  }

  writer, err := os.Create(comp_file)
  if err != nil {
    fmt.Println(err)
    return false
  }

  options := jpeg.Options{Quality: quality}

  err = jpeg.Encode(writer, initImage, &options)
  if err != nil {
    fmt.Println(err)
    return false
  } else if verbose {
    fmt.Printf("Saved %v with %v%% quality\n", comp_file, quality)
  }
  return true
}

func main() {

  files := getFiles(file)

  // Process files
  for _, file := range files {
    errored := processFile(file)
    if !errored {
      fmt.Printf("%v could not be processed\n", file)
    }
  }

  fmt.Println("Done!")
  
}
