package main

import (
  "os"
  "io/ioutil"
  "fmt"
  "strings"
  "image"
  "image/jpeg"
  "path/filepath"
  "strconv"
  "math"
  "github.com/integrii/flaggy"
  "github.com/disintegration/imaging"
)

var (
  file string
  disable_resize = false
  quality = 80
  suffix = "_compressed"
  files []string
)

func init() {
  flaggy.SetName("compress")
  flaggy.SetDescription("Compress your images")
  flaggy.SetVersion("0.1")

  flaggy.AddPositionalValue(&file, "file", 1, true, "Image file to compress")
  flaggy.Bool(&disable_resize, "n", "no-resize", "Keep image at original size")
  flaggy.Int(&quality, "q", "quality", "Quality to save image at 0-100")
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
  // Simple half size for now
  width := curSize.X
  height := curSize.Y
  ratio := float64(height) / float64(width)
  fmt.Printf("Ratio: %v\n", ratio)
  pixcount := 2073600.0
  newWidth := math.Sqrt(pixcount/ratio)
  newHeight := newWidth*ratio
  fmt.Printf("Width: %v Height: %v\n", int(newWidth), int(newHeight))
  *initImage = imaging.Resize(*initImage, int(newWidth), int(newHeight), imaging.MitchellNetravali)
}

func main() {
  fmt.Println(file)
  fmt.Println(strconv.FormatBool(disable_resize))
  fmt.Println(strconv.Itoa(quality))

  info, err := os.Lstat(file)
  if err != nil {
    os.Exit(1)
  }
  
  if info.IsDir() {
    fmt.Println("Passed directory")
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

  // Process files
  for _, file := range files {
    comp_file := getNewFilename(file)
    fmt.Println(comp_file)

    reader, err := os.Open(file)
    if err != nil {
      os.Exit(1)
    }
    initImage, format, err := image.Decode(reader)
    fmt.Println("Decoded " + format)
    if err != nil {
      os.Exit(1)
    }

    if !disable_resize {
      fmt.Println("Resizing")
      resizeImage(&initImage)
    }

    writer, err := os.Create(comp_file)
    if err != nil {
      os.Exit(1)
    }

    options := jpeg.Options{Quality: quality}

    err = jpeg.Encode(writer, initImage, &options)
    fmt.Println("Success!")
  }
  
}
