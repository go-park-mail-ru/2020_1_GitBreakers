package git

import gogitConfig "github.com/go-git/go-git/v5/config"

func (repo Repository) AddRemote(protocol ActionProtocol, remoteName, absRemotePath string) (Remote, error) {
	gogitRemote, err := repo.CreateRemote(&gogitConfig.RemoteConfig{
		Name: remoteName,
		URLs: []string{ConvertToProtoURL(protocol, absRemotePath)},
	})
	if err != nil {
		return Remote{}, err
	}
	return Remote{gogitRemote}, nil
}
