package container

import (
	"archive/tar"
	"context"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/containers/image/v5/types"
	"github.com/containers/storage/drivers/copy"
	"github.com/pkg/errors"

	"github.com/hpcng/warewulf/internal/pkg/config"
	"github.com/hpcng/warewulf/internal/pkg/oci"
	"github.com/hpcng/warewulf/internal/pkg/util"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
)

func ImportDocker(uri string, name string, sCtx *types.SystemContext) error {
	OciBlobCacheDir := config.LocalStateDir + "/oci/blobs"

	err := os.MkdirAll(OciBlobCacheDir, 0755)
	if err != nil {
		return err
	}

	if !ValidName(name) {
		return errors.New("VNFS name contains illegal characters: " + name)
	}

	fullPath := RootFsDir(name)

	err = os.MkdirAll(fullPath, 0755)
	if err != nil {
		return err
	}

	p, err := oci.NewPuller(
		oci.OptSetBlobCachePath(OciBlobCacheDir),
		oci.OptSetSystemContext(sCtx),
	)
	if err != nil {
		return err
	}

	if _, err := p.GenerateID(context.Background(), uri); err != nil {
		return err
	}

	if err := p.Pull(context.Background(), uri, fullPath); err != nil {
		return err
	}

	return nil
}

func ImportDirectory(uri string, name string) error {
	fullPath := RootFsDir(name)

	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		return err
	}

	if !util.IsDir(uri) {
		return errors.New("Import directory does not exist: " + uri)
	}

	if !util.IsFile(path.Join(uri, "/bin/sh")) {
		return errors.New("Source directory has no /bin/sh: " + uri)
	}

	err = copy.DirCopy(uri, fullPath, copy.Content, true)
	if err != nil {
		return err
	}

	return nil
}

func ImportTar(uri string, name string) error {

	if !util.IsFile(uri) {
		return errors.New("Tarball was not found: " + uri)
	}

	wwlog.Printf(wwlog.VERBOSE, "Creating temporary directory for files extracted from tarball\n")
	tmpDir, err := ioutil.TempDir(os.TempDir(), ".wwctl-tarfiles-")
	if err != nil {
		return errors.New("Unable to create temporary path for tarball import of file: " + uri)
	}
	defer os.RemoveAll(tmpDir);
	
	wwlog.Printf(wwlog.VERBOSE, "Opening tarball\n")
	tarfd, err := os.Open(uri)
	if err != nil {
		wwlog.Printf(wwlog.VERBOSE, "Unable to open tarball %s\n", uri)
		return err
	}
	defer tarfd.Close()
	
	wwlog.Printf(wwlog.VERBOSE, "Reading contents of tarball\n")
	tarRead := tar.NewReader(tarfd)
	for {

		hdr, err := tarRead.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			wwlog.Printf(wwlog.VERBOSE, "Unable to read tarball %s\n", uri)
			return err
		}

		tmpFileInfo := hdr.FileInfo()
		tmpFilepath := filepath.Join(tmpDir, hdr.Name)

		if tmpFileInfo.IsDir() {
			if err := os.MkdirAll(tmpDir, 0755); err != nil {
				return errors.New("Unable to create directory from tarball")
			}
			continue
		}

		tmpFile, err := os.OpenFile(tmpFilepath, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
		if err != nil {
			return errors.New("Unable to create file from tarball")
		}

		if _, err := io.Copy(tmpFile, tarRead); err != nil {
			wwlog.Printf(wwlog.VERBOSE, "Error extracting contents from tarball %s\n", uri)
			return err
		}
		tmpFile.Close()
	}

	wwlog.Printf(wwlog.VERBOSE, "Importing container extracted from tarball\n", uri)
	if ImportDirectory(tmpDir, name) != nil {
		return errors.New("Unable to import container extracted from tarball file: " + uri)	
	}

	return nil
}
