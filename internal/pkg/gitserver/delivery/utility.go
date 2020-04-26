package delivery

import (
	"github.com/sosedoff/gitkit"
	"net/http"
	"regexp"
	"strings"
)

var reSlashDedup = regexp.MustCompile(`/{2,}`)

func getCredential(request *http.Request) (cred gitkit.Credential, ok bool) {
	cred.Username, cred.Password, ok = request.BasicAuth()
	return cred, ok
}

func getRepoInfo(serviceSuffix string, request *http.Request) repoBasicInfo {
	path := strings.Replace(request.URL.Path, serviceSuffix, "", 1)

	ownerLogin, repoName := getNamespaceAndRepo(path)

	return repoBasicInfo{
		ownerLogin: ownerLogin,
		repoName:   repoName,
	}
}

func getNamespaceAndRepo(input string) (string, string) {
	if input == "" || input == "/" {
		return "", ""
	}

	// Remove duplicate slashes
	input = reSlashDedup.ReplaceAllString(input, "/")

	// Remove leading slash
	if input[0] == '/' && input != "/" {
		input = input[1:]
	}

	blocks := strings.Split(input, "/")
	num := len(blocks)

	if num < 2 {
		return "", blocks[0]
	}

	return strings.Join(blocks[0:num-1], "/"), blocks[num-1]
}
