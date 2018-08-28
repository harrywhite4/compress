package main

import (
  "os"
  "fmt"
  "strings"
  "image"
  "image/jpeg"
  "path/filepath"
  "strconv"
  "github.com/integrii/flaggy"
)

var (
  file string
  disable_resize = true
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
    // Get all files in directory
    os.Exit(1)
  } else {
    files = make([]string, 1)
    files[0] = file
  }

  for _, file := range files {
    comp_file := getNewFilename(file)
    fmt.Println(comp_file)

    reader, err := os.Open(file)
    if err != nil {
      os.Exit(1)
    }
    image, format, err := image.Decode(reader)
    fmt.Println("Decoded " + format)
    if err != nil {
      os.Exit(1)
    }

    if !disable_resize {
      fmt.Println("Resizing")
    }

    writer, err := os.Create(comp_file)
    if err != nil {
      os.Exit(1)
    }

    options := jpeg.Options{Quality: quality}

    err = jpeg.Encode(writer, image, &options)
    fmt.Println("Success!")
  }
  
}
