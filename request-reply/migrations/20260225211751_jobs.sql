CREATE TABLE IF NOT EXISTS jobs(
    id SERIAL PRIMARY KEY,
    job_id VARCHAR(36) NOT NULL UNIQUE,
    status VARCHAR(255) NOT NULL,
    created_at DATE DEFAULT CURRENT_DATE
);

CREATE INDEX idx_jobs_job_id ON jobs(job_id);