DROP TABLE movie;
CREATE TABLE movie
(
  id integer primary key,
  name varchar(255),
  year integer,
  resolution varchar(255),
  filetype varchar(255),
  location varchar(255),
  picture varchar(255)
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
