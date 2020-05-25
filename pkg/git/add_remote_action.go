package git

import gogitConfig "github.com/go-git/go-git/v5/config"

func (repo Repository) AddRemote(protocol ActionProtocol, absRemote string) (Remote, error) {
	gogitRemote, err := repo.CreateRemote(&gogitConfig.RemoteConfig{
		Name: UpstreamRemoteName,
		URLs: []string{convertToProtoURL(protocol, absRemote)},
	})
	if err != nil {
		return Remote{}, err
	}
	return Remote{gogitRemote}, nil
}
