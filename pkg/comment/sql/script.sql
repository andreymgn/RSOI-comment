CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE comments (
    uid UUID PRIMARY KEY,
    user_uid UUID NOT NULL,
    post_uid UUID NOT NULL,
    body TEXT NOT NULL,
    parent_uid UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    modified_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);
