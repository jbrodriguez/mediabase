DROP TABLE movie;
DROP TABLE moviename;

CREATE TABLE movie
(
  name text,
  year integer,
  resolution text,
  filetype text,
  location text,
  cover text,
  backdrop text
);
CREATE INDEX movie_name_idx ON movie (name);
CREATE UNIQUE INDEX movie_location_idx ON movie (location);

CREATE VIRTUAL TABLE moviename USING fts4(content="movie", name);
CREATE TRIGGER movie_bu BEFORE UPDATE ON movie BEGIN
	DELETE FROM moviename WHERE docid=old.rowid;
END;

CREATE TRIGGER movie_bd BEFORE DELETE ON movie BEGIN
	DELETE FROM moviename WHERE docid=old.rowid;
END;

CREATE TRIGGER movie_au AFTER UPDATE ON movie BEGIN
	INSERT INTO moviename(docid, name) VALUES (new.rowid, new.name);
END;

CREATE TRIGGER movie_ai AFTER INSERT ON movie BEGIN
	INSERT INTO moviename(docid, name) VALUES (new.rowid, new.name);
END;
