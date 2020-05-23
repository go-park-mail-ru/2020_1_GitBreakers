package git

import "fmt"

func convertToProtoURL(protocol ActionProtocol, absPath string) string {
	return fmt.Sprintf("%s://%s", string(protocol), absPath)
}
