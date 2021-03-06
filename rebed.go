// Package rebed brings simple embedded file functionality
// to Go's new embed directive.
//
// It can recreate the directory structure
// from the embed.FS type with or without
// the files it contains. This is useful to
// expose the filesystem to the end user so they
// may see and modify the files.
package rebed

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// folderPerm MkdirAll is called with this permission to prevent restricted folders
// from being created.  0755=rwxr-xr-x
const folderPerm os.FileMode = 0755

// Tree creates the target filesystem folder structure.
func Tree(fsys embed.FS) error {
	return Walk(fsys, ".", func(dirpath string, de fs.DirEntry) error {
		fullpath := filepath.Join(dirpath, de.Name())
		if de.IsDir() {
			return os.MkdirAll(fullpath, folderPerm)
		}
		return nil
	})
}

// Touch creates the target filesystem folder structure in the binary's
// current working directory with empty files. Does not modify
// already existing files.
func Touch(fsys embed.FS) error {
	return Walk(fsys, ".", func(dirpath string, de fs.DirEntry) error {
		fullpath := filepath.Join(dirpath, de.Name())
		if de.IsDir() {
			return os.MkdirAll(fullpath, folderPerm)
		}
		// unsure how IsNotExist works. this could be improved
		_, err := os.Stat(fullpath)
		if os.IsNotExist(err) {
			_, err = os.Create(fullpath)
		}
		return err
	})
}

// Create overwrites files of same path/name
// in binaries current working directory or
// creates new ones if not exist.
func Create(fsys embed.FS) error {
	return Walk(fsys, ".", func(dirpath string, de fs.DirEntry) error {
		fullpath := filepath.Join(dirpath, de.Name())
		if de.IsDir() {
			return os.MkdirAll(fullpath, folderPerm)
		}
		return embedCopyToFile(fsys, fullpath)
	})
}

// Patch creates files which are missing in
// FS filesystem. Does not modify existing files
func Patch(fsys embed.FS) error {
	return Walk(fsys, ".", func(dirpath string, de fs.DirEntry) error {
		fullpath := filepath.Join(dirpath, de.Name())
		if de.IsDir() {
			return os.MkdirAll(fullpath, folderPerm)
		}
		_, err := os.Stat(fullpath)
		if os.IsNotExist(err) {
			_, err = os.Create(fullpath)
		}
		return err
	})
}

// embedCopyToFile copies an embedded file's contents
// to a file machine in same relative path
func embedCopyToFile(fsys embed.FS, path string) error {
	fi, err := fsys.Open(path)
	if err != nil {
		return err
	}
	fo, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(fo, fi)
	return err
}

// Walk expects a path to a directory.
// f called on every file/directory found recursively.
// It is not guaranteed to stay in main package import path.
//
// f's first argument is the relative/absolute path to directory being scanned.
func Walk(fsys embed.FS, startPath string, f func(path string, de fs.DirEntry) error) error {
	folders := make([]string, 0) // buffer of folders to process
	WalkDir(fsys, startPath, func(dirpath string, de fs.DirEntry) error {
		if de.IsDir() {
			folders = append(folders, filepath.Join(dirpath, de.Name()))
		}
		return f(dirpath, de)
	})
	n := len(folders)
	for n != 0 {
		for i := 0; i < n; i++ {
			WalkDir(fsys, folders[i], func(dirpath string, de fs.DirEntry) error {
				if de.IsDir() {
					folders = append(folders, filepath.Join(dirpath, de.Name()))
				}
				return f(dirpath, de)
			})
		}
		// we process n folders at a time, add new folders while
		//processing n folders, then discard those n folders once finished
		// and resume with a new n list of folders
		var newFolders int = len(folders) - n
		folders = folders[n : n+newFolders] // if found 0 new folders, end
		n = len(folders)
	}
	return nil
}

// WalkDir applies f to every file/folder in embedded directory fsys.
// It is not guaranteed to stay in main package import path.
//
// f's first argument is the relative/absolute path to directory being scanned.
func WalkDir(fsys embed.FS, startPath string, f func(path string, de fs.DirEntry) error) error {
	items, err := fsys.ReadDir(startPath)
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := f(startPath, item); err != nil {
			return err
		}
	}
	return nil
}
