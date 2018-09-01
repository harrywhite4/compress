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

var (
  file string
  disableResize = false
  halfResize = false
  quality = 80
  suffix = "_compressed"
  targetPixels = 2073600
  files []string
)

func init() {
  flaggy.SetName("compress")
  flaggy.SetDescription("Compress your images")
  flaggy.SetVersion("0.1")

  flaggy.AddPositionalValue(&file, "file", 1, true, "Image file to compress")
  flaggy.Bool(&disableResize, "n", "no-resize", "Keep image at original size")
  flaggy.Bool(&halfResize, "2", "half", "Save image at half it's original size")
  flaggy.Int(&quality, "q", "quality", "Quality to save image at 0-100")
  flaggy.Int(&targetPixels, "p", "pixels", "Target pixel count for resized image")
  flaggy.String(&suffix, "s", "suffix", "Suffix to be appended to filenames")
  flaggy.Parse()
}

func getNewFilename(path string) string {
  dir, filename := filepath.Split(path)
  splitname := strings.Split(filename, ".")
  return dir + splitname[0] + suffix + ".jpg"
}

func resizeImage(initImage *image.Image) {
  curSize := (*initImage).Bounds().Size()
  width := curSize.X
  height := curSize.Y
  currentPixels := width*height

  if currentPixels > targetPixels {
    ratio := float64(height) / float64(width)
    fmt.Printf("Ratio: %v\n", ratio)
    newWidth := math.Sqrt(float64(targetPixels)/ratio)
    newHeight := newWidth*ratio
    fmt.Printf("Width: %v Height: %v\n", int(newWidth), int(newHeight))
    *initImage = imaging.Resize(*initImage, int(newWidth), int(newHeight), imaging.MitchellNetravali)
  }
}

func getFiles(file string) []string {
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

func processFile(file string) {
  comp_file := getNewFilename(file)
  fmt.Println(comp_file)

  reader, err := os.Open(file)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  initImage, format, err := image.Decode(reader)
  fmt.Println("Decoded " + format)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  if !disableResize {
    fmt.Println("Resizing...")
    resizeImage(&initImage)
  }

  writer, err := os.Create(comp_file)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  options := jpeg.Options{Quality: quality}

  err = jpeg.Encode(writer, initImage, &options)
  fmt.Println("Success!")
}

func main() {

  files := getFiles(file)

  // Process files
  for _, file := range files {
    processFile(file)
  }
  
}
