# Compress
Command line image compression tool written in GO

## Description
Command line program to downsize images and save them as jpegs.
By default it will resize an image down to ~2M pixels (preserving
aspect ratio) and save it as a 80% quality jpeg

The default number of pixels is the same as a 1920 x 1080 image.
If an image already has less pixels, no resize will be done

## Usage
compress path1 path2 ... [options]

Path can be an image file or a directory. 
If it is a direcotory, all images within that directory will be processed

### Options
Short flag | Long flag | Description
--- | --- | ---
-2 | --half | Save image at half it's original size
-h | --help | Display help
-n | --no-resize | Keep image at original size
-p | --pixels |Target pixel count for resized image (default 2073600)
-q | --quality | Quality to save image at 0-100 (default 80)
-s | --suffix | Suffix to be appended to filenames (default "_compressed")
-v | --verbose | Print additional information during processing

