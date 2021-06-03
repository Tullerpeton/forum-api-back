CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE users (
    nickname CITEXT NOT NULL PRIMARY KEY,
    fullname TEXT NOT NULL,
    about TEXT,
    email TEXT NOT NULL,

    CONSTRAINT email_unique UNIQUE (email)
);

CREATE UNLOGGED TABLE forums (
    slug TEXT NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    author_nickname CITEXT NOT NULL,
    count_posts INTEGER NOT NULL DEFAULT 0,
    count_threads INTEGER NOT NULL DEFAULT 0,

    FOREIGN KEY (author_nickname) REFERENCES users(nickname)
);

CREATE UNLOGGED TABLE threads (
    id SERIAL NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    author_id INTEGER  NOT NULL,
    forum_id INTEGER NOT NULL,
    message TEXT NOT NULL,
    slug TEXT NOT NULL,
    date_created IMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (author_id) REFERENCES users(id),
    FOREIGN KEY (forum_id) REFERENCES forums(id),

    CONSTRAINT slug_unique UNIQUE (slug)
);

CREATE UNLOGGED TABLE posts (
    id SERIAL NOT NULL PRIMARY KEY,
    parent_message_id INTEGER NOT NULL,
    author_id INTEGER  NOT NULL,
    message TEXT NOT NULL,
    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    forum_id INTEGER NOT NULL,
    thread_id INTEGER NOT NULL,
    date_created IMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (author_id) REFERENCES users(id),
    FOREIGN KEY (forum_id) REFERENCES forums(id),
    FOREIGN KEY (thread_id) REFERENCES threads(id)
);

CREATE UNLOGGED TABLE votes (
    vote INTEGER NOT NULL,
    author_id INTEGER  NOT NULL,
    thread_id INTEGER NOT NULL,

    CONSTRAINT vote_unique UNIQUE (thread_id, author_id)
);



CREATE FUNCTION inc_posts_counter() RETURNS TRIGGER AS $$
BEGIN
UPDATE forums SET
    count_posts = count_posts + 1;
WHERE slug = NEW.forum_slug;
RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_inc_posts_counter
    AFTER INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE inc_posts_counter();

CREATE FUNCTION inc_threads_counter() RETURNS TRIGGER AS $$
BEGIN
UPDATE forums SET
    count_threads = count_threads + 1;
WHERE slug = NEW.forum_slug;
RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_inc_threads_counter
    AFTER INSERT
    ON threads
    FOR EACH ROW
    EXECUTE PROCEDURE inc_threads_counter();

