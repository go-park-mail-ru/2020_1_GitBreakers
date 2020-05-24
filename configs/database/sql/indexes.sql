-- Indexes

CREATE INDEX IF NOT EXISTS issues_repository_id_idx ON issues (repository_id);
CREATE INDEX IF NOT EXISTS news_repository_id_idx ON news (repository_id);
