DROP DATABASE IF EXISTS forum_db;
CREATE DATABASE forum_db
    WITH OWNER postgres
    LC_COLLATE = 'C'
    LC_CTYPE = 'en_US.utf8'
    TEMPLATE template0;
\connect forum_db;

CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE users (
    nickname CITEXT NOT NULL PRIMARY KEY,
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT NOT NULL,

    CONSTRAINT email_unique UNIQUE (email)
);

CREATE UNIQUE INDEX ON users (nickname, email);

CREATE UNLOGGED TABLE forums (
    slug CITEXT NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    author_nickname CITEXT NOT NULL,
    count_posts INTEGER NOT NULL DEFAULT 0,
    count_threads INTEGER NOT NULL DEFAULT 0,

    FOREIGN KEY (author_nickname) REFERENCES users(nickname)
);

CREATE UNLOGGED TABLE threads (
    id SERIAL NOT NULL PRIMARY KEY,
    slug CITEXT,
    title TEXT NOT NULL,
    author_nickname CITEXT NOT NULL,
    forum_slug CITEXT NOT NULL,
    message TEXT NOT NULL,
    votes INTEGER NOT NULL DEFAULT 0,

    date_created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (author_nickname) REFERENCES users(nickname),
    FOREIGN KEY (forum_slug) REFERENCES forums(slug),

    CONSTRAINT slug_unique UNIQUE (slug)
);

CREATE INDEX ON threads (slug) WHERE slug IS NOT NULL;
CREATE INDEX ON threads(forum_slug, date_created);
CREATE INDEX ON threads(date_created);

CREATE UNLOGGED TABLE authors (
    id SERIAL NOT NULL PRIMARY KEY,
    user_nickname CITEXT NOT NULL,
    forum_slug CITEXT NOT NULL,

    CONSTRAINT author_unique UNIQUE (user_nickname, forum_slug),

    FOREIGN KEY (user_nickname) REFERENCES users(nickname),
    FOREIGN KEY (forum_slug) REFERENCES forums(slug)
);

CREATE INDEX ON authors(user_nickname, forum_slug);

CREATE UNLOGGED TABLE posts (
    id SERIAL NOT NULL PRIMARY KEY,
    parent_message_id INTEGER NOT NULL,
    author_nickname CITEXT NOT NULL,
    message TEXT NOT NULL,
    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    forum_slug CITEXT NOT NULL,
    thread_id INTEGER NOT NULL,
    date_created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    path_of_nesting INTEGER ARRAY DEFAULT '{}' NOT NULL,

    FOREIGN KEY (author_nickname) REFERENCES users(nickname),
    FOREIGN KEY (forum_slug) REFERENCES forums(slug),
    FOREIGN KEY (thread_id) REFERENCES threads(id)
);

CREATE INDEX ON posts((path_of_nesting[1]));
CREATE INDEX ON posts(id, (path_of_nesting[1]));
CREATE UNIQUE INDEX ON posts(id, thread_id);
CREATE UNIQUE INDEX ON posts(id, author_nickname);
CREATE INDEX ON posts(thread_id, path_of_nesting, id);
CREATE INDEX ON posts(thread_id, id);

CREATE UNLOGGED TABLE votes (
    vote INTEGER NOT NULL,
    author_nickname CITEXT NOT NULL,
    thread_id INTEGER NOT NULL,

    CONSTRAINT vote_unique UNIQUE (thread_id, author_nickname),

    FOREIGN KEY (author_nickname) REFERENCES users(nickname),
    FOREIGN KEY (thread_id) REFERENCES threads(id)
);

CREATE UNIQUE INDEX ON votes(author_nickname, thread_id);


CREATE FUNCTION inc_posts_counter() RETURNS TRIGGER AS $$
BEGIN
    UPDATE forums SET
    count_posts = count_posts + 1
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
    count_threads = count_threads + 1
    WHERE slug = NEW.forum_slug;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_inc_threads_counter
    AFTER INSERT
    ON threads
    FOR EACH ROW
    EXECUTE PROCEDURE inc_threads_counter();

CREATE FUNCTION update_threads_votes() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        IF (OLD.vote != NEW.vote) THEN
            UPDATE threads SET
            votes = votes + NEW.vote - OLD.vote
            WHERE id = NEW.thread_id;
        END IF;
    ELSIF TG_OP = 'INSERT' THEN
        UPDATE threads SET
        votes = (votes + new.vote)
        WHERE id = new.thread_id;
        RETURN new;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_threads_votes
    AFTER UPDATE OR INSERT
    ON votes
    FOR EACH ROW
    EXECUTE PROCEDURE update_threads_votes();


CREATE FUNCTION update_path_of_nesting() RETURNS TRIGGER AS $$
DECLARE
    parent_thread INTEGER;
BEGIN
    IF NEW.parent_message_id = 0 THEN
        NEW.path_of_nesting = ARRAY [NEW.id];
    ELSE
        SELECT thread_id
        INTO parent_thread
        FROM posts
        WHERE id = new.parent_message_id;

        IF parent_thread ISNULL THEN
            RAISE EXCEPTION 'Parent post not found %', NEW.parent_message_id;
        ELSIF parent_thread <> NEW.thread_id THEN
            RAISE EXCEPTION 'Thread not found %', NEW.thread_id;
        END IF;

        SELECT path_of_nesting
        INTO NEW.path_of_nesting
        FROM posts
        WHERE id = NEW.parent_message_id;
        NEW.path_of_nesting = array_append(NEW.path_of_nesting, NEW.id);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_path_of_nesting
    BEFORE INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE update_path_of_nesting();


CREATE FUNCTION update_user_author_status() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO authors (user_nickname, forum_slug)
    VALUES (NEW.author_nickname, NEW.forum_slug)
    ON CONFLICT DO NOTHING;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_posts_update_user_author_status
    AFTER INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE update_user_author_status();

CREATE TRIGGER trigger_threads_update_user_author_status
    AFTER INSERT
    ON threads
    FOR EACH ROW
    EXECUTE PROCEDURE update_user_author_status();
