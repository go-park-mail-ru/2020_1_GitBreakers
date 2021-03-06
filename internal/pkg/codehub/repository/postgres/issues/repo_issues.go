package issues

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type IssueRepository struct {
	DB *sqlx.DB
}

func NewIssueRepository(db *sqlx.DB) IssueRepository {
	return IssueRepository{
		DB: db,
	}
}

func (repo *IssueRepository) CreateIssue(issue models.Issue) error {
	_, err := repo.DB.Exec(`
		INSERT INTO issues (author_id, repository_id, title, message, label) VALUES ($1, $2, $3, $4, $5)
		`,
		issue.AuthorID,
		issue.RepoID,
		issue.Title,
		issue.Message,
		issue.Label,
	)
	if err == nil {
		return nil
	}
	if err, ok := err.(*pq.Error); ok {
		switch err.Code {
		case "23505":
			return entityerrors.AlreadyExist()
		default:
			return errors.Wrapf(err, "error occurs in IssuesRepository in CreateIssue "+
				"with issue=%+v", issue)

		}
	}
	return err

}

func (repo *IssueRepository) UpdateIssue(issue models.Issue) error {
	_, err := repo.DB.Exec(`
			UPDATE issues 
			SET author_id = $2,
				repository_id = $3,
				title = $4,
				message = $5,
				label = $6
			WHERE id = $1
		`, issue.ID,
		issue.AuthorID,
		issue.RepoID,
		issue.Title,
		issue.Message,
		issue.Label,
	)
	if err != nil {
		return errors.Wrapf(err, "error occurs in IssuesRepository in UpdateIssue "+
			"with issue=%+v", issue)
	}

	return nil
}

// CloseIssue return entityerrors.Invalid() if issueID is not valid ot this issue already closed
func (repo *IssueRepository) CloseIssue(issueID int64) error {
	var isAffected bool
	err := repo.DB.QueryRow(`
			UPDATE issues 
			SET is_closed = TRUE
			WHERE id = $1 AND is_closed = FALSE RETURNING TRUE AS result`,
		issueID,
	).Scan(
		&isAffected,
	)

	switch {
	case err == sql.ErrNoRows:
		return entityerrors.Invalid()
	case err != nil:
		return errors.Wrapf(err, "error occurs in IssuesRepository in CloseIssue "+
			"with issueID=%+v", issueID)
	}

	return nil
}

func (repo *IssueRepository) GetIssuesList(repoID int64, limit int64, offset int64) (issues []models.Issue, err error) {
	rows, err := repo.DB.Query(`
			SELECT 	id,
				   	author_id,
				   	repository_id,
				   	title,
				   	message,
				   	label,
				   	is_closed,
				   	created_at,
			       	user_login,
			       	user_avatar_path
			FROM issues_users_view
			WHERE repository_id = $1
			LIMIT $2 OFFSET $3`,

		repoID, limit, offset,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "error occurs in IssuesRepository in GetIssuesList "+
			"with repoID=%v, limit=%v, offset=%v", repoID, limit, offset)
	}

	defer func() {
		if errRows := rows.Close(); errRows != nil {
			err = errors.Wrap(err, errRows.Error())
		}
	}()

	for rows.Next() {
		issue := models.Issue{}
		err = rows.Scan(
			&issue.ID,
			&issue.AuthorID,
			&issue.RepoID,
			&issue.Title,
			&issue.Message,
			&issue.Label,
			&issue.IsClosed,
			&issue.CreatedAt,
			&issue.AuthorLogin,
			&issue.AuthorImage,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "error occurs in IssuesRepository in GetIssuesList "+
				"while scanning issues with repoID=%v, limit=%v, offset=%v", repoID, limit, offset)
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

func (repo *IssueRepository) GetOpenedIssuesList(repoID int64, limit int64, offset int64) (issues []models.Issue, err error) {
	rows, err := repo.DB.Query(`
			SELECT 	id,
				   	author_id,
				   	repository_id,
				   	title,
				   	message,
				   	label,
				   	is_closed,
				   	created_at,
			       	user_login,
			       	user_avatar_path
			FROM issues_users_view
			WHERE repository_id = $1 AND is_closed = FALSE
			LIMIT $2 OFFSET $3`,

		repoID, limit, offset,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "error occurs in IssuesRepository in GetOpenedIssuesList "+
			"with repoID=%v, limit=%v, offset=%v", repoID, limit, offset)
	}

	defer func() {
		if errRows := rows.Close(); errRows != nil {
			err = errors.Wrap(err, errRows.Error())
		}
	}()

	for rows.Next() {
		issue := models.Issue{}
		err = rows.Scan(
			&issue.ID,
			&issue.AuthorID,
			&issue.RepoID,
			&issue.Title,
			&issue.Message,
			&issue.Label,
			&issue.IsClosed,
			&issue.CreatedAt,
			&issue.AuthorLogin,
			&issue.AuthorImage,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "error occurs in IssuesRepository in GetOpenedIssuesList "+
				"while scanning issues with repoID=%v, limit=%v, offset=%v", repoID, limit, offset)
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

func (repo *IssueRepository) GetClosedIssuesList(repoID int64, limit int64, offset int64) (issues []models.Issue, err error) {
	rows, err := repo.DB.Query(`
			SELECT 	id,
				   	author_id,
				   	repository_id,
				   	title,
				   	message,
				   	label,
				   	is_closed,
				   	created_at,
			       	user_login,
			       	user_avatar_path
			FROM issues_users_view
			WHERE repository_id = $1 AND is_closed = TRUE
			LIMIT $2 OFFSET $3`,

		repoID, limit, offset,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "error occurs in IssuesRepository in GetClosedIssuesList "+
			"with repoID=%v, limit=%v, offset=%v", repoID, limit, offset)
	}

	defer func() {
		if errRows := rows.Close(); errRows != nil {
			err = errors.Wrap(err, errRows.Error())
		}
	}()

	for rows.Next() {
		issue := models.Issue{}
		err = rows.Scan(
			&issue.ID,
			&issue.AuthorID,
			&issue.RepoID,
			&issue.Title,
			&issue.Message,
			&issue.Label,
			&issue.IsClosed,
			&issue.CreatedAt,
			&issue.AuthorLogin,
			&issue.AuthorImage,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "error occurs in IssuesRepository in GetClosedIssuesList "+
				"while scanning issues with repoID=%v, limit=%v, offset=%v", repoID, limit, offset)
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

func (repo *IssueRepository) CheckEditAccessIssue(userID, issueID int64) (perm.Permission, error) {
	var issueAuthorId int64
	var issueRepoId int64

	err := repo.DB.QueryRow(
		`SELECT author_id, repository_id FROM issues WHERE id = $1`,
		issueID,
	).Scan(&issueAuthorId, &issueRepoId)

	switch {
	case err == sql.ErrNoRows:
		return perm.NoAccess(), entityerrors.DoesNotExist()
	case err != nil:
		return perm.NoAccess(), errors.Wrapf(err, "error occurs in IssuesRepository"+
			" while getting issueAuthorId and issueRepoId in CheckEditAccessIssue with userID=%v, issueID=%v",
			userID, issueID)
	}

	if issueAuthorId == userID {
		return perm.OwnerAccess(), nil
	}

	var gitRepoRole string
	err = repo.DB.QueryRow(
		`SELECT role FROM users_git_repositories WHERE repository_id = $1 AND user_id = $2`,
		issueRepoId, userID).Scan(&gitRepoRole)

	switch {
	case err == sql.ErrNoRows:
		return perm.NoAccess(), nil
	case err != nil:
		return perm.NoAccess(), errors.Wrapf(err, "error occurs in IssuesRepository"+
			" in CheckEditAccessIssue while checking repo access with userID=%v, issueID=%v", userID, issueID)
	}

	return perm.Permission(gitRepoRole), nil
}

func (repo *IssueRepository) CheckAccessRepo(userID, repoID int64) (perm.Permission, error) {
	panic("implement panic")
}

func (repo *IssueRepository) GetIssue(issueID int64) (issue models.Issue, err error) {
	err = repo.DB.QueryRow(`
			SELECT 	id,
					author_id,
					repository_id,
					title,
					message,
					label,
					is_closed,
					created_at,
			       	user_login,
			       	user_avatar_path
			FROM issues_users_view
			WHERE id = $1`, issueID,
	).Scan(
		&issue.ID,
		&issue.AuthorID,
		&issue.RepoID,
		&issue.Title,
		&issue.Message,
		&issue.Label,
		&issue.IsClosed,
		&issue.CreatedAt,
		&issue.AuthorLogin,
		&issue.AuthorImage,
	)

	switch {
	case err == sql.ErrNoRows:
		return issue, entityerrors.DoesNotExist()
	case err != nil:
		return issue, errors.Wrapf(err, "error occurs in IssuesRepository in GetIssue "+
			"with issueID=%v", issueID)
	}

	return issue, nil
}
