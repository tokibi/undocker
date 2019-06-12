package untar

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	OverwriteSymlinkRefs bool
}

func Untar(r io.Reader, dir string, opts Options) error {
	return untar(r, dir, opts)
}

func untar(r io.Reader, dir string, opts Options) error {
	tr := tar.NewReader(r)
	for {
		f, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if f == nil {
			continue
		}
		rel := filepath.FromSlash(f.Name)
		abs := filepath.Join(dir, rel)

		switch {
		case f.Typeflag == tar.TypeDir:
			if _, err := os.Stat(abs); err != nil {
				if err := os.MkdirAll(abs, 0755); err != nil {
					return err
				}
			}

			if err = os.Chmod(abs, f.FileInfo().Mode()); err != nil {
				return err
			}
			if err = os.Lchown(abs, f.Uid, f.Gid); err != nil {
				return err
			}

		case f.Typeflag == tar.TypeReg:
			// whiteout file
			if strings.Contains(abs, ".wh.") {
				rm := strings.Replace(abs, ".wh.", "", 1)
				os.RemoveAll(rm)
				continue
			}
			wf, err := os.OpenFile(abs, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(f.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(wf, tr); err != nil {
				return err
			}
			wf.Close()

			if err = os.Chmod(abs, f.FileInfo().Mode()); err != nil {
				return err
			}
			if err = os.Lchown(abs, f.Uid, f.Gid); err != nil {
				return err
			}

		case f.Typeflag == tar.TypeSymlink:
			if opts.OverwriteSymlinkRefs {
				os.Symlink(filepath.Join(dir, f.Linkname), abs)
			} else {
				os.Symlink(f.Linkname, abs)
			}

			if err = os.Lchown(abs, f.Uid, f.Gid); err != nil {
				return err
			}

		case f.Typeflag == tar.TypeLink:
			os.Link(filepath.Join(dir, f.Linkname), abs)

			if err = os.Chmod(abs, f.FileInfo().Mode()); err != nil {
				return err
			}
			if err = os.Lchown(abs, f.Uid, f.Gid); err != nil {
				return err
			}
		}
	}
	return nil
}
