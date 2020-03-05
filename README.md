# xlspaint

A spreadsheet-painting tool, "Hoganshi" renderer.

![lenna_rendered.png](./docs/lenna_rendered.png)

# prerequisites

- macOS, Linux, Windows
- viewer for .xlsx file

- Accepted image file format are {png, jpeg, gif}

# usage

## run

download latest release, or build it by yourself, then 

```console
xlspaint [your-favorite-image-file]
```

## get packages & build

```console
go get
go build
```

# How it works

1. Read an image file specified in the first commandline argument
2. Crop the image (make it squared)
3. Resize the image to 256 x 256 
4. Make 256-color-palette 
5. Read the template .xls
6. Render each pixels
7. Write book


# License

MIT

