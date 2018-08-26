package main

import (
  "fmt"
  "github.com/integrii/flaggy"
)

func init() {
  flaggy.SetName("compress")
  flaggy.SetDescription("Compress your images")
  flaggy.SetVersion("0.1")


}

func main() {
  var file = "default.img"
  flaggy.String(&file, "f", "file", "Image file to compress")
  flaggy.Parse()

  fmt.Printf(file)
}
