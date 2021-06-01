CREATE UNLOGGED TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    nickname TEXT NOT NULL,
    fullname TEXT NOT NULL,
    about TEXT,
    email TEXT NOT NULL,

    CONSTRAINT nickname_unique UNIQUE (nickname),
    CONSTRAINT email_unique UNIQUE (email)
);

CREATE UNLOGGED TABLE forums (
    id SERIAL NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    slug TEXT NOT NULL,

    FOREIGN KEY (author_id) REFERENCES users(id),

    CONSTRAINT slug_unique UNIQUE (slug)
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

