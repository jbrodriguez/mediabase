DROP TABLE IF EXISTS movie;
DROP TABLE IF EXISTS movietitle;
DROP INDEX IF EXISTS movie_filetype_idx;
DROP INDEX IF EXISTS movie_location_idx;
DROP INDEX IF EXISTS movie_title_idx;

DROP TRIGGER IF EXISTS movie_ai;
DROP TRIGGER IF EXISTS movie_au;
DROP TRIGGER IF EXISTS movie_bd;
DROP TRIGGER IF EXISTS movie_bu;

CREATE TABLE movie
(
  title text,
  original_title text,
  file_title text,
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
  backdrop text,
  genres text,
  vote_average integer,
  vote_count integer,
  countries text,
  added text,
  modified text,
  last_watched text,
  all_watched text,
  count_watched integer,
  score integer
);
CREATE INDEX movie_title_idx ON movie (title);
CREATE INDEX movie_location_idx ON movie (location);
CREATE INDEX movie_filetype_idx ON movie (filetype);

CREATE VIRTUAL TABLE movietitle USING fts4(content="movie", title, original_title, file_title);
CREATE TRIGGER movie_bu BEFORE UPDATE ON movie BEGIN
	DELETE FROM movietitle WHERE docid=old.rowid;
END;

CREATE TRIGGER movie_bd BEFORE DELETE ON movie BEGIN
	DELETE FROM movietitle WHERE docid=old.rowid;
END;

CREATE TRIGGER movie_au AFTER UPDATE ON movie BEGIN
	INSERT INTO movietitle(docid, title, original_title, file_title) VALUES (new.rowid, new.title, new.original_title, new.file_title);
END;

CREATE TRIGGER movie_ai AFTER INSERT ON movie BEGIN
	INSERT INTO movietitle(docid, title, original_title, file_title) VALUES (new.rowid, new.title, new.original_title, new.file_title);
END;
