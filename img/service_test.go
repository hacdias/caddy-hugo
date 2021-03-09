package img

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func TestService_Resize(t *testing.T) {
	testCases := map[string]struct {
		options []Option
		width   int
		height  int
		source  func(t *testing.T) afero.File
		matcher func(t *testing.T, reader io.Reader)
		wantErr bool
	}{
		"fill upscale": {
			options: []Option{WithMode(ResizeModeFill)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 50, 20)
			},
			matcher: sizeMatcher(100, 100),
		},
		"fill downscale": {
			options: []Option{WithMode(ResizeModeFill)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 200, 150)
			},
			matcher: sizeMatcher(100, 100),
		},
		"fit upscale": {
			options: []Option{WithMode(ResizeModeFit)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 50, 20)
			},
			matcher: sizeMatcher(50, 20),
		},
		"fit downscale": {
			options: []Option{WithMode(ResizeModeFit)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 200, 150)
			},
			matcher: sizeMatcher(100, 75),
		},
		"keep original format": {
			options: []Option{},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayPng(t, 200, 150)
			},
			matcher: formatMatcher(FormatPng),
		},
		"convert to jpeg": {
			options: []Option{WithFormat(FormatJpeg)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 200, 150)
			},
			matcher: formatMatcher(FormatJpeg),
		},
		"convert to png": {
			options: []Option{WithFormat(FormatPng)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 200, 150)
			},
			matcher: formatMatcher(FormatPng),
		},
		"convert to gif": {
			options: []Option{WithFormat(FormatGif)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 200, 150)
			},
			matcher: formatMatcher(FormatGif),
		},
		"convert to tiff": {
			options: []Option{WithFormat(FormatTiff)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 200, 150)
			},
			matcher: formatMatcher(FormatTiff),
		},
		"convert to bmp": {
			options: []Option{WithFormat(FormatBmp)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 200, 150)
			},
			matcher: formatMatcher(FormatBmp),
		},
		"convert to unknown": {
			options: []Option{WithFormat(Format(-1))},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 200, 150)
			},
			matcher: formatMatcher(FormatJpeg),
		},
		"resize png": {
			options: []Option{WithMode(ResizeModeFill)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayPng(t, 200, 150)
			},
			matcher: sizeMatcher(100, 100),
		},
		"resize gif": {
			options: []Option{WithMode(ResizeModeFill)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayGif(t, 200, 150)
			},
			matcher: sizeMatcher(100, 100),
		},
		"resize tiff": {
			options: []Option{WithMode(ResizeModeFill)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayTiff(t, 200, 150)
			},
			matcher: sizeMatcher(100, 100),
		},
		"resize bmp": {
			options: []Option{WithMode(ResizeModeFill)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayBmp(t, 200, 150)
			},
			matcher: sizeMatcher(100, 100),
		},
		"resize heic": {
			options: []Option{WithMode(ResizeModeFill)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return openFile(t, "testdata/image3.heic")
			},
			matcher: sizeMatcher(100, 100),
		},
		"resize with high quality": {
			options: []Option{WithMode(ResizeModeFill), WithQuality(QualityHigh)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 200, 150)
			},
			matcher: sizeMatcher(100, 100),
		},
		"resize with medium quality": {
			options: []Option{WithMode(ResizeModeFill), WithQuality(QualityMedium)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 200, 150)
			},
			matcher: sizeMatcher(100, 100),
		},
		"resize with low quality": {
			options: []Option{WithMode(ResizeModeFill), WithQuality(QualityLow)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 200, 150)
			},
			matcher: sizeMatcher(100, 100),
		},
		"resize with unknown quality": {
			options: []Option{WithMode(ResizeModeFill), WithQuality(Quality(-1))},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return newGrayJpeg(t, 200, 150)
			},
			matcher: sizeMatcher(100, 100),
		},
		"get thumbnail from file with APP0 JFIF": {
			options: []Option{WithMode(ResizeModeFill), WithQuality(QualityLow)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return openFile(t, "testdata/gray-sample.jpg")
			},
			matcher: sizeMatcher(125, 128),
		},
		"get thumbnail from file without APP0 JFIF": {
			options: []Option{WithMode(ResizeModeFill), WithQuality(QualityLow)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return openFile(t, "testdata/20130612_142406.jpg")
			},
			matcher: sizeMatcher(320, 240),
		},
		"resize from file without IFD1 thumbnail": {
			options: []Option{WithMode(ResizeModeFill), WithQuality(QualityLow)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return openFile(t, "testdata/IMG_2578.JPG")
			},
			matcher: sizeMatcher(100, 100),
		},
		"resize for higher quality levels": {
			options: []Option{WithMode(ResizeModeFill), WithQuality(QualityMedium)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				return openFile(t, "testdata/gray-sample.jpg")
			},
			matcher: sizeMatcher(100, 100),
		},
		"broken file": {
			options: []Option{WithMode(ResizeModeFit)},
			width:   100,
			height:  100,
			source: func(t *testing.T) afero.File {
				t.Helper()
				fs := afero.NewMemMapFs()
				file, err := fs.Create("image.jpg")
				require.NoError(t, err)

				_, err = file.WriteString("this is not an image")
				require.NoError(t, err)

				return file
			},
			wantErr: true,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			svc := New(1)
			source := test.source(t)
			defer source.Close()

			buf := &bytes.Buffer{}
			err := svc.Resize(context.Background(), source, test.width, test.height, buf, test.options...)
			if (err != nil) != test.wantErr {
				t.Fatalf("GetMarketSpecs() error = %v, wantErr %v", err, test.wantErr)
			}
			if err != nil {
				return
			}
			test.matcher(t, buf)
		})
	}
}

func sizeMatcher(width, height int) func(t *testing.T, reader io.Reader) {
	return func(t *testing.T, reader io.Reader) {
		resizedImg, _, err := image.Decode(reader)
		require.NoError(t, err)

		require.Equal(t, width, resizedImg.Bounds().Dx())
		require.Equal(t, height, resizedImg.Bounds().Dy())
	}
}

func formatMatcher(format Format) func(t *testing.T, reader io.Reader) {
	return func(t *testing.T, reader io.Reader) {
		_, decodedFormat, err := image.DecodeConfig(reader)
		require.NoError(t, err)

		require.Equal(t, format.String(), decodedFormat)
	}
}

func newGrayJpeg(t *testing.T, width, height int) afero.File {
	fs := afero.NewMemMapFs()
	file, err := fs.Create("image.jpg")
	require.NoError(t, err)

	img := image.NewGray(image.Rect(0, 0, width, height))
	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	require.NoError(t, err)

	_, err = file.Seek(0, io.SeekStart)
	require.NoError(t, err)

	return file
}

func newGrayPng(t *testing.T, width, height int) afero.File {
	fs := afero.NewMemMapFs()
	file, err := fs.Create("image.png")
	require.NoError(t, err)

	img := image.NewGray(image.Rect(0, 0, width, height))
	err = png.Encode(file, img)
	require.NoError(t, err)

	_, err = file.Seek(0, io.SeekStart)
	require.NoError(t, err)

	return file
}

func newGrayGif(t *testing.T, width, height int) afero.File {
	fs := afero.NewMemMapFs()
	file, err := fs.Create("image.gif")
	require.NoError(t, err)

	img := image.NewGray(image.Rect(0, 0, width, height))
	err = gif.Encode(file, img, nil)
	require.NoError(t, err)

	_, err = file.Seek(0, io.SeekStart)
	require.NoError(t, err)

	return file
}

func newGrayTiff(t *testing.T, width, height int) afero.File {
	fs := afero.NewMemMapFs()
	file, err := fs.Create("image.tiff")
	require.NoError(t, err)

	img := image.NewGray(image.Rect(0, 0, width, height))
	err = tiff.Encode(file, img, nil)
	require.NoError(t, err)

	_, err = file.Seek(0, io.SeekStart)
	require.NoError(t, err)

	return file
}

func newGrayBmp(t *testing.T, width, height int) afero.File {
	fs := afero.NewMemMapFs()
	file, err := fs.Create("image.bmp")
	require.NoError(t, err)

	img := image.NewGray(image.Rect(0, 0, width, height))
	err = bmp.Encode(file, img)
	require.NoError(t, err)

	_, err = file.Seek(0, io.SeekStart)
	require.NoError(t, err)

	return file
}

func openFile(t *testing.T, name string) afero.File {
	appfs := afero.NewOsFs()
	file, err := appfs.Open(name)

	require.NoError(t, err)

	return file
}

func TestService_FormatFromExtension(t *testing.T) {
	testCases := map[string]struct {
		ext     string
		want    Format
		wantErr error
	}{
		"jpg": {
			ext:  ".jpg",
			want: FormatJpeg,
		},
		"jpeg": {
			ext:  ".jpeg",
			want: FormatJpeg,
		},
		"png": {
			ext:  ".png",
			want: FormatPng,
		},
		"gif": {
			ext:  ".gif",
			want: FormatGif,
		},
		"tiff": {
			ext:  ".tiff",
			want: FormatTiff,
		},
		"bmp": {
			ext:  ".bmp",
			want: FormatBmp,
		},
		"heic": {
			ext:  ".heic",
			want: FormatHeic,
		},
		"heif": {
			ext:  ".heif",
			want: FormatHeic,
		},
		"unknown": {
			ext:     ".mov",
			wantErr: ErrUnsupportedFormat,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			svc := New(1)
			got, err := svc.FormatFromExtension(test.ext)
			require.Truef(t, errors.Is(err, test.wantErr), "error = %v, wantErr %v", err, test.wantErr)
			if err != nil {
				return
			}
			require.Equal(t, test.want, got)
		})
	}
}
