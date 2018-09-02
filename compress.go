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
  "github.com/spf13/pflag"
  "github.com/disintegration/imaging"
)
// Import png for initialisation
import _ "image/png"

var (
  file string
  disableResize bool
  halfResize bool
  help bool
  quality int
  suffix = "_compressed"
  targetPixels int
  files []string
  verbose = false
  allowedExtensions = [3]string{".jpg", ".jpeg", ".png"}
  flagSet *pflag.FlagSet
  programText = "Compress - Compress your images"
  version = 0.1
)

func init() {
  flagSet = pflag.NewFlagSet("compress", pflag.ExitOnError) 

  flagSet.BoolVarP(&disableResize, "no-resize", "n", false, "Keep image at original size")
  flagSet.BoolVarP(&halfResize, "half", "2", false, "Save image at half it's original size")
  flagSet.BoolVarP(&help, "help", "h", false, "Display help")
  flagSet.IntVarP(&quality, "quality", "q", 80, "Quality to save image at 0-100")
  flagSet.IntVarP(&targetPixels, "pixels", "p", 2073600, "Target pixel count for resized image")
  flagSet.StringVarP(&suffix, "suffix", "s", "_compressed", "Suffix to be appended to filenames")
  flagSet.BoolVarP(&verbose, "verbose", "v", false, "Print additional information during processing")
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

func processFile(file string) error {
  fmt.Printf("Processing %v\n", file)

  allowed := checkFileExtension(file)

  if !allowed {
    if verbose {
      fmt.Printf("Skipping %v not valid file extension\n", file)
    }
    return nil
  }
  comp_file := getNewFilename(file)

  reader, err := os.Open(file)
  if err != nil {
    return err
  }

  initImage, format, err := image.Decode(reader)
  if err != nil {
    return err
  } else if verbose {
    fmt.Println("Decoded " + format)
  }

  if !disableResize {
    resizeImage(&initImage)
  }

  writer, err := os.Create(comp_file)
  if err != nil {
    return err
  }

  options := jpeg.Options{Quality: quality}

  err = jpeg.Encode(writer, initImage, &options)
  if err != nil {
    return err
  } else if verbose {
    fmt.Printf("Saved %v with %v%% quality\n", comp_file, quality)
  }
  return nil
}

func printHelp() {
  usage := "compress path1 path2 ... [options]"
  fmt.Printf("%v\nVersion: %v\n\nUsage:\n  %v\n\nOptions:\n", programText, version, usage)
  flagSet.PrintDefaults()
}

func main() {
  flagSet.Parse(os.Args[1:])

  if help {
    printHelp()
    os.Exit(1)
  }

  positionals := flagSet.Args()
  if len(positionals) == 0 {
    fmt.Println("No path specified")
    fmt.Println("run compress -h for usage")
    os.Exit(1)
  } 
  files := getFiles(positionals)

  // Process files
  for _, file := range files {
    err := processFile(file)
    if err != nil {
      fmt.Printf("%v failed with error %v\n", file, err)
    }
  }

  fmt.Println("Done!")
  
}
