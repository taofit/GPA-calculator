BEGIN;

CREATE TABLE IF NOT EXISTS grade (
    id SERIAL PRIMARY KEY,
    school_id INT,
    student_id INT,
    course_id INT,
    grade VARCHAR(2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    UNIQUE(school_id, student_id, course_id)
);

CREATE TABLE IF NOT EXISTS grade_scale (
    id SERIAL PRIMARY KEY,
    school_id INT,
    grade VARCHAR(2),
    scale NUMERIC(2,1),
    percent INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    UNIQUE(school_id, grade)
);

COMMIT;