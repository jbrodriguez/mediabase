DROP TABLE movie;
DROP TABLE moviefts;

CREATE TABLE movie
(
  title text,
  original_title text,
  year integer,
  runtime integer,
  tmdb_id integer,
  imdb_id text,
  overview text,
  tagline text,
  resolution text,
  filetype text,
  location text,
  cover text,
  backdrop text
);
CREATE INDEX movie_title_idx ON movie (title);

CREATE VIRTUAL TABLE moviefts USING fts4(content="movie", title, original_title);
CREATE TRIGGER movie_bu BEFORE UPDATE ON movie BEGIN
	DELETE FROM moviefts WHERE docid=old.rowid;
END;

CREATE TRIGGER movie_bd BEFORE DELETE ON movie BEGIN
	DELETE FROM moviefts WHERE docid=old.rowid;
END;

CREATE TRIGGER movie_au AFTER UPDATE ON movie BEGIN
	INSERT INTO moviefts(docid, title, original_title) VALUES (new.rowid, new.title, new.original_title);
END;

CREATE TRIGGER movie_ai AFTER INSERT ON movie BEGIN
	INSERT INTO moviefts(docid, title, original_title) VALUES (new.rowid, new.title, new.original_title);
END;
