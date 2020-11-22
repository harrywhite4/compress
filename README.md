# Compress

Command line image compression tool written in GO

## Description

Command line program to downsize images and save them as jpegs.
By default it will resize an image down to ~2M pixels (preserving
aspect ratio) and save it as a 80% quality jpeg

The default number of pixels is the same as a 1920 x 1080 image.
If an image already has less pixels, no resize will be done

## Usage

compress [OPTION]... [PATH]...

Path can be an image file or a directory. If it is a direcotory, all images within that directory
will be processed

See `compress --help` for options.
