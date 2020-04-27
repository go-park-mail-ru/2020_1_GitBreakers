package postgres

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	"github.com/jmoiron/sqlx"
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
		INSERT INTO issues (author_id, repo_id, title, message, label) VALUES ($1, $2, $3, $4, $5)
		`,
		issue.AuthorID,
		issue.RepoID,
		issue.Title,
		issue.Message,
		issue.Label,
	)
	if err != nil {
		return errors.Wrapf(err, "error occurs in IssuesRepository in CreateIssue "+
			"with issue=%+v", issue)
	}

	return nil
}

func (repo *IssueRepository) UpdateIssue(issue models.Issue) error {
	_, err := repo.DB.Exec(`
			UPDATE issues 
			SET author_id = $2,
				repo_id = $3,
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
			WHERE id = $1 AND is_closed = FALSE RETURNING TRUE
		`,
		issueID).Scan(&isAffected)

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
				   	repo_id,
				   	title,
				   	message,
				   	label,
				   	is_closed,
				   	created_at
			FROM issues
			WHERE repo_id = $1
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
				   	repo_id,
				   	title,
				   	message,
				   	label,
				   	is_closed,
				   	created_at
			FROM issues
			WHERE repo_id = $1 AND is_closed = FALSE
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
				   	repo_id,
				   	title,
				   	message,
				   	label,
				   	is_closed,
				   	created_at
			FROM issues
			WHERE repo_id = $1 AND is_closed = TRUE
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
		)
		if err != nil {
			return nil, errors.Wrapf(err, "error occurs in IssuesRepository in GetClosedIssuesList "+
				"while scanning issues with repoID=%v, limit=%v, offset=%v", repoID, limit, offset)
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

func (repo *IssueRepository) CheckAccessIssue(userID, issueID int64) (perm.Permission, error) {
	panic("implement me") // TODO
}

func (repo *IssueRepository) CheckAccessRepo(userID, repoID int64) (perm.Permission, error) {
	panic("implement me") // TODO use repository for GitRepository
}

func (repo *IssueRepository) GetIssue(issueID int64) (issue models.Issue, err error) {
	err = repo.DB.QueryRow(`
			SELECT 	id,
					author_id,
					repo_id,
					title,
					message,
					label,
					is_closed,
					created_at
			FROM issues
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
