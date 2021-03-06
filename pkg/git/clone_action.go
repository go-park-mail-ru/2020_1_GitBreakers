package git

import (
	gogit "github.com/go-git/go-git/v5"
	gogitPlumbing "github.com/go-git/go-git/v5/plumbing"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/pkg/errors"
)

type ActionProtocol string

func CloneBranchOnly(proto ActionProtocol, absSrcPath, absDstPath, srcBranch,
	remoteName string, depth int) (Repository, error) {

	if proto == ActionNoneProtocol {
		return Repository{}, errors.Wrap(entityerrors.Invalid(), "invalid protocol")
	}

	gogitRepo, err := gogit.PlainClone(absDstPath, false, &gogit.CloneOptions{
		URL:           ConvertToProtoURL(proto, absSrcPath),
		RemoteName:    remoteName,
		ReferenceName: gogitPlumbing.NewBranchReferenceName(srcBranch),
		SingleBranch:  true,
		NoCheckout:    false,
		Depth:         depth,
		Progress:      nil,
		Tags:          gogit.NoTags,
	})
	if err != nil {
		return Repository{}, errors.WithStack(err)
	}

	return Repository{gogitRepo}, nil
}
