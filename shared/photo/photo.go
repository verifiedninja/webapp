package photo

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/wujiang/imaging"
)

const (
	jpgQuality = 95
)

// PNGToImage converts PNG from POST request to image resource
func PNGToImage(file multipart.File) (image.Image, error) {
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return img, nil
}

// GIFToImage converts GIF from POST request to image resource
func GIFToImage(file multipart.File) (image.Image, error) {
	img, err := gif.Decode(file)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return img, nil
}

// JPGToFile converts JPG from POST request to file on disk
func JPGToFile(file multipart.File, output string) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}

	return nil
}

// ImageToJPGFile converts image resource to file on disk
func ImageToJPGFile(img image.Image, output string) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	var opt jpeg.Options
	opt.Quality = jpgQuality
	err = jpeg.Encode(out, img, &opt)
	if err != nil {
		return err
	}
	out.Close()

	return nil
}

// FileDimensions returns the image dimensions of the file
func FileDimensions(imagePath string) (image.Config, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return image.Config{}, err
	}
	defer file.Close()

	i, _, err := image.DecodeConfig(file)
	if err != nil {
		return image.Config{}, err
	}

	return i, nil
}

// ImageDimensions returns the image dimensions of the image
func ImageDimensions(file multipart.File) (image.Config, error) {
	i, _, err := image.DecodeConfig(file)
	if err != nil {
		return image.Config{}, err
	}

	// Reset the reader
	file.Seek(0, 0)

	return i, nil
}

// FixRotation rotates the image in place to correct from cell phones
func FixRotation(src string) error {
	ff, err := os.Open(src)
	if err != nil {
		return err
	}
	defer ff.Close()

	x, err := exif.Decode(ff)
	if err != nil {
		// This throws an EOF when it does not have any EXIF data
		return err
	}

	tag, err := x.Get(exif.Orientation)
	if err != nil {
		return err
	}

	ot, err := tag.Int(0)
	if err != nil {
		return err
	}

	img, err := imaging.Open(src)
	if err != nil {
		return err
	}

	// exif standard
	// http://www.daveperrett.com/articles/2012/07/28/exif-orientation-handling-is-a-ghetto/
	switch {
	case ot == 2:
		img = imaging.FlipH(img)
	case ot == 3:
		img = imaging.Rotate180(img)
	case ot == 4:
		img = imaging.FlipH(img)
		img = imaging.Rotate180(img)
	case ot == 5:
		img = imaging.FlipV(img)
		img = imaging.Rotate270(img)
	case ot == 6:
		img = imaging.Rotate270(img)
	case ot == 7:
		img = imaging.FlipV(img)
		img = imaging.Rotate90(img)
	case ot == 8:
		img = imaging.Rotate90(img)
	}

	err = imaging.Save(img, src, jpgQuality)
	if err != nil {
		return err
	}

	return nil
}
