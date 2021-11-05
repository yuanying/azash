package thumbnail

import (
	"os"

	"github.com/disintegration/imaging"
)

func Generate(destPath, srcPath string) (err error) {
	_, err = os.Stat(destPath)
	if !os.IsNotExist(err) {
		if err != nil {
			return
		}
		return nil
	}

	src, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}
	dst := imaging.Fill(src, 364, 514, imaging.Center, imaging.Lanczos)
	err = imaging.Save(dst, destPath)
	return
}
