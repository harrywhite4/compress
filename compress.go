package main

import (
  "fmt"
  "strconv"
  "github.com/integrii/flaggy"
)

var file = "default.img"
var disable_resize = true
var quality = 80

func init() {
  flaggy.SetName("compress")
  flaggy.SetDescription("Compress your images")
  flaggy.SetVersion("0.1")

  flaggy.AddPositionalValue(&file, "file", 1, true, "Image file to compress")
  flaggy.Bool(&disable_resize, "n", "no-resize", "Keep image at original size")
  flaggy.Int(&quality, "q", "quality", "Quality to save image at 0-100")
  flaggy.Parse()
}

func main() {
  fmt.Printf(file + "\n")
  fmt.Printf(strconv.FormatBool(disable_resize) + "\n")
  fmt.Printf(strconv.Itoa(quality) + "\n")
}
