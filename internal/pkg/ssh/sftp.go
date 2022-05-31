package ssh

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func InitSftpClient(sshClient *ssh.Client) (*sftp.Client, error) {
	return sftp.NewClient(sshClient)
}

func UploadFile(localPath, remotePath string, sftpClient *sftp.Client) error {
	fileName := path.Base(localPath)
	srcFile, err := os.Open(localPath)
	if err != nil {
		return err // open file fail
	}
	defer srcFile.Close()

	// creates the named file mode 0666 (before umask), truncating it if it already exists
	sftpClient.MkdirAll(remotePath) // If path is already a directory, MkdirAll does nothing and returns nil
	dstFile, err := sftpClient.Create(path.Join(remotePath, fileName))
	if err != nil {
		return err // create file fail
	}
	defer dstFile.Close()

	// chmod executable
	dstFile.Chmod(0777)

	// better?
	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		return err
	}
	// bytes, err := ioutil.ReadAll(srcFile)
	// if err != nil {
	// 	return err // readall src file fail
	// }
	// if _, err := dstFile.Write(bytes); err != nil {
	// 	return err // write dst file fail
	// }
	return nil
}

func UploadDir(localPath, remotePath string, sftpClient *sftp.Client) error {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		return err // read dir fail
	}

	for _, file := range localFiles {
		localFilePath := path.Join(localPath, file.Name())
		remoteFilePath := path.Join(remotePath, file.Name())
		if file.IsDir() {
			sftpClient.MkdirAll(remoteFilePath) // If path contains a regular file, an error is returned, discard it
			if err := UploadDir(localFilePath, remoteFilePath, sftpClient); err != nil {
				return err
			}
		} else {
			if err := UploadFile(path.Join(localPath, file.Name()), remotePath, sftpClient); err != nil {
				return err
			}
		}
	}
	return nil
}

func DownLoadFile(localPath, remotePath string, sftpClient *sftp.Client) error {
	srcFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return err // read file fail
	}
	defer srcFile.Close()

	localFile := path.Base(remotePath)
	dstFile, err := os.Create(path.Join(localPath, localFile))
	if err != nil {
		return err // create file fail
	}
	defer dstFile.Close()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		return err // write file fail
	}
	return nil
}
