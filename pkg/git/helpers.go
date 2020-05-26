package git

import (
	"fmt"
	gogitPlumbing "github.com/go-git/go-git/v5/plumbing"
	"path"
)

func ConvertToProtoURL(protocol ActionProtocol, absPath string) string {
	return fmt.Sprintf("%s://%s", string(protocol), absPath)
}

func CreateRemoteRefName(remoteName, branchName string) string {
	refName := gogitPlumbing.NewRemoteReferenceName(remoteName, branchName).String()
	return path.Clean(refName)
}

func CreateBranchRefName(branchName string) string {
	refName := gogitPlumbing.NewBranchReferenceName(branchName).String()
	return path.Clean(refName)
}
