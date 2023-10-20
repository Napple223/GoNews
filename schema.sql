/* postgres */

DROP TABLE IF EXISTS posts;

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    pub_time INTEGER DEFAULT 0,
    link TEXT NOT NULL UNIQUE 
)