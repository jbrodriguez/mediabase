package services

import (
	"apertoire.net/mediabase/server/bus"
	"apertoire.net/mediabase/server/helper"
	"apertoire.net/mediabase/server/message"
	"database/sql"
	"github.com/apertoire/mlog"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
)

type Dal struct {
	Bus    *bus.Bus
	Config *helper.Config
	db     *sql.DB
	dbase  string
	err    error
	cnt    int

	storeMovie      *sql.Stmt
	searchMovies    *sql.Stmt
	listMovies      *sql.Stmt
	listByRuntime   *sql.Stmt
	listMoviesToFix *sql.Stmt
	// getAssets       *sql.Stmt
	// getRevisions    *sql.Stmt
	// getItems        *sql.Stmt
	// getCategories   *sql.Stmt
	// getLastRevision *sql.Stmt
	// getAsset        *sql.Stmt

	// putAsset    *sql.Stmt
	// putProduct  *sql.Stmt
	// putRevision *sql.Stmt
	// putItem     *sql.Stmt
}

func (self *Dal) prepare(sql string) *sql.Stmt {
	stmt, err := self.db.Prepare(sql)
	if err != nil {
		mlog.Fatalf("prepare sql: %s (%s)", err, sql)
	}
	return stmt
}

func (self *Dal) Start() {
	mlog.Info("starting dal service ...")

	self.dbase = filepath.Join(self.Config.AppDir, "/db/mediabase.db")
	self.db, self.err = sql.Open("sqlite3", self.dbase)
	if self.err != nil {
		mlog.Fatalf("open database: %s (%s)", self.err, self.dbase)
	}

	self.cnt = 0

	self.searchMovies = self.prepare("select dt.rowid, dt.title, dt.original_title, dt.year, dt.runtime, dt.tmdb_id, dt.imdb_id, dt.overview, dt.tagline, dt.resolution, dt.filetype, dt.location, dt.cover, dt.backdrop, dt.genres, dt.vote_average, dt.vote_count, dt.countries, dt.added, dt.modified, dt.last_watched, dt.all_watched, dt.count_watched, dt.score from movie dt, movietitle vt where vt.movietitle match ? and dt.rowid = vt.docid order by dt.title;")
	self.listMovies = self.prepare("select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched, score from movie order by title")
	self.listByRuntime = self.prepare("select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched, score from movie order by runtime")
	self.listMoviesToFix = self.prepare("select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched, score from movie where original_title = 'FIXMOV23'")

	mlog.Info("connected to database")

	// self.initSchema()

	go self.react()
}

func (self *Dal) Stop() {
	self.listMoviesToFix.Close()
	self.listByRuntime.Close()
	self.listMovies.Close()
	self.searchMovies.Close()
	// self.storeMovie.Close()
	self.db.Close()

	mlog.Info("dal service stopped")
}

func (self *Dal) react() {
	for {
		select {
		case msg := <-self.Bus.StoreMovie:
			self.doStoreMovie(msg)
		case msg := <-self.Bus.DeleteMovie:
			self.doDeleteMovie(msg)
		case msg := <-self.Bus.UpdateMovie:
			self.doUpdateMovie(msg)
		case msg := <-self.Bus.GetMovies:
			go self.doGetMovies(msg)
		case msg := <-self.Bus.ListMovies:
			go self.doListMovies(msg)
		case msg := <-self.Bus.ShowDuplicates:
			go self.doShowDuplicates(msg)
		case msg := <-self.Bus.ListByRuntime:
			go self.doListByRuntime(msg)
		case msg := <-self.Bus.SearchMovies:
			go self.doSearchMovies(msg)
		case msg := <-self.Bus.CheckMovie:
			go self.doCheckMovie(msg)
		case msg := <-self.Bus.GetMoviesToFix:
			go self.doGetMoviesToFix(msg)
		}
	}
}

func (self *Dal) initSchema() {
	sql := `
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

	`

	_, err := self.db.Exec(sql)
	if err != nil {
		mlog.Info("%q: %s", err, sql)
		return
	}

	mlog.Info("inited schema")
}

func (self *Dal) doCheckMovie(msg *message.CheckMovie) {
	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare("select rowid from movie where location = ?")
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at prepare: %s", err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(msg.Movie.Location).Scan(&id)

	// if err == sql.ErrNoRows {
	// 	mlog.Fatalf("id = %d, err = %d", id, err)
	// }

	// mlog.Fatalf("gone and done")
	if err != sql.ErrNoRows && err != nil {
		tx.Rollback()
		mlog.Fatalf("at queryrow: %s", err)
	}

	tx.Commit()

	msg.Result <- (id != 0)
}

func (self *Dal) doStoreMovie(movie *message.Movie) {
	self.cnt++

	mlog.Info("STARTED SAVING %s [%d]", movie.Title, self.cnt)

	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	// stmt, err := tx.Prepare("insert into movie(title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, director, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	stmt, err := tx.Prepare("insert into movie(title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched, score) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at prepare: %s", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(movie.Title, movie.Original_Title, movie.File_Title, movie.Year, movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline, movie.Resolution, movie.FileType, movie.Location, movie.Cover, movie.Backdrop,
		movie.Genres, movie.Vote_Average, movie.Vote_Count, movie.Production_Countries, movie.Added, movie.Modified, movie.Last_Watched, movie.All_Watched, movie.Count_Watched, movie.Score)
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at exec: %s", err)
	}

	// mlog.Info("Movie is %v", movie)

	// _, self.err = self.storeMovie.Exec(movie.Title, movie.Year, movie.Resolution, movie.FileType, movie.Location)
	// if self.err != nil {
	// 	mlog.Fatalf("at storemovie: %s", self.err)
	// }

	tx.Commit()
	mlog.Info("FINISHED SAVING %s [%d]", movie.Title, self.cnt)

	// _, self.err = self.storeMovie.Exec(movie.Name, movie.Year, movie.Resolution, movie.Type, movie.Path, movie.Picture)
	// if self.err != nil {
	// 	mlog.Fatalf(self.err)
	// }
}

func (self *Dal) doDeleteMovie(movie *message.Movie) {
	mlog.Info("STARTED DELETING [%d] %s", movie.Id, movie.Title)

	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	// stmt, err := tx.Prepare("insert into movie(title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, director, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	stmt, err := tx.Prepare("delete from movie where rowid = ?")
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at prepare: %s", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(movie.Id)
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at exec: %s", err)
	}

	// mlog.Info("Movie is %v", movie)

	// _, self.err = self.storeMovie.Exec(movie.Title, movie.Year, movie.Resolution, movie.FileType, movie.Location)
	// if self.err != nil {
	// 	mlog.Fatalf("at storemovie: %s", self.err)
	// }

	tx.Commit()
	mlog.Info("FINISHED DELETING [%d] %s", movie.Id, movie.Title)

	// _, self.err = self.storeMovie.Exec(movie.Name, movie.Year, movie.Resolution, movie.Type, movie.Path, movie.Picture)
	// if self.err != nil {
	// 	mlog.Fatalf(self.err)
	// }
}

func (self *Dal) doUpdateMovie(movie *message.Movie) {
	mlog.Info("STARTED UPDATING %s [%d]", movie.Title, self.cnt)

	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare("update movie set title = ?, original_title = ?, file_title = ?, year = ?, runtime = ?, imdb_id = ?, overview = ?, tagline = ?, cover = ?, backdrop = ?, genres = ?, vote_average = ?, vote_count = ?, countries = ?, modified = ? where tmdb_id = ?")
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at prepare: %s", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(movie.Title, movie.Original_Title, movie.File_Title, movie.Year, movie.Runtime, movie.Imdb_Id, movie.Overview, movie.Tagline, movie.Cover, movie.Backdrop, movie.Genres, movie.Vote_Average, movie.Vote_Count, movie.Production_Countries, movie.Modified, movie.Tmdb_Id)
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at exec: %s", err)
	}

	tx.Commit()
	mlog.Info("FINISHED UPDATING %s", movie.Title)
}

func (self *Dal) doGetMovies(msg *message.GetMovies) {
	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("unable to begin transaction: %s", err)
	}

	stmt, err := tx.Prepare("select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched from movie order by added desc limit ?")
	if err != nil {
		mlog.Fatalf("unable to prepare transaction: %s", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(30)
	if err != nil {
		mlog.Fatalf("unable to run transaction: %s", err)
	}

	items := make([]*message.Movie, 0)

	for rows.Next() {
		movie := message.Movie{}
		rows.Scan(&movie.Id, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, movie.All_Watched, &movie.Count_Watched)
		items = append(items, &movie)
	}
	rows.Close()

	tx.Commit()

	msg.Reply <- items
}

func (self *Dal) doListMovies(msg *message.ListMovies) {
	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("unable to begin transaction: %s", err)
	}

	rows, err := self.listMovies.Query()
	if err != nil {
		mlog.Fatalf("unable to prepare transaction: %s", self.err)
	}

	items := make([]*message.Movie, 0)

	self.cnt = 0

	for rows.Next() {
		movie := message.Movie{}
		rows.Scan(&movie.Id, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score)
		items = append(items, &movie)
		self.cnt++
	}
	rows.Close()

	tx.Commit()

	mlog.Info("Listed %d movies", self.cnt)

	msg.Reply <- items
}

func (self *Dal) doListByRuntime(msg *message.Movies) {
	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("unable to begin transaction: %s", err)
	}

	rows, err := self.listByRuntime.Query()
	if err != nil {
		mlog.Fatalf("unable to prepare transaction: %s", self.err)
	}

	items := make([]*message.Movie, 0)

	self.cnt = 0

	for rows.Next() {
		movie := message.Movie{}
		rows.Scan(&movie.Id, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched)
		items = append(items, &movie)
		self.cnt++
	}
	rows.Close()

	tx.Commit()

	mlog.Info("Listed (runtime) %d movies", self.cnt)

	msg.Reply <- items
}

func (self *Dal) doShowDuplicates(msg *message.Movies) {
	mlog.Info("started from the bottom now we're here")

	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("unable to begin transaction: %s", err)
	}

	// rows, err := self.listMovies.Query()
	// if err != nil {
	// 	mlog.Fatalf("unable to prepare transaction: %s", err)
	// }

	// rows, err := self.db.Query("select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched from movie where title in (select title from movie group by title having count(*) > 1);")
	rows, err := self.db.Query("select a.rowid, a.title, a.original_title, a.file_title, a.year, a.runtime, a.tmdb_id, a.imdb_id, a.overview, a.tagline, a.resolution, a.filetype, a.location, a.cover, a.backdrop, a.genres, a.vote_average, a.vote_count, a.countries, a.added, a.modified, a.last_watched, a.all_watched, a.count_watched from movie a join (select title, year from movie group by title, year having count(*) > 1) b on a.title = b.title and a.year = b.year;")
	if err != nil {
		mlog.Fatalf("unable to prepare transaction: %s", self.err)
	}

	items := make([]*message.Movie, 0)

	self.cnt = 0

	for rows.Next() {
		movie := message.Movie{}
		rows.Scan(&movie.Id, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched)
		items = append(items, &movie)
		self.cnt++
	}
	rows.Close()

	tx.Commit()

	mlog.Info("Found %d duplicate movies", self.cnt)

	msg.Reply <- items
}

func (self *Dal) doSearchMovies(msg *message.SearchMovies) {
	term := msg.Term + "*"
	mlog.Info("this is: %s", term)

	rows, err := self.searchMovies.Query(term)
	if err != nil {
		mlog.Fatalf("unable to begin transaction: %s", self.err)
	}

	items := make([]*message.Movie, 0)

	for rows.Next() {
		movie := message.Movie{}
		rows.Scan(&movie.Id, &movie.Title, &movie.Original_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score)
		// movie := &message.Movie{}
		// rows.Scan(movie.Id, movie.Title, movie.Original_Title, movie.Year, movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline, movie.Resolution, movie.FileType, movie.Location, movie.Cover, movie.Backdrop)
		// mlog.Info("title: (%s)", movie.Title)
		items = append(items, &movie)
	}
	rows.Close()

	msg.Reply <- items
}

func (self *Dal) doGetMoviesToFix(msg *message.Movies) {
	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("unable to begin transaction: %s", err)
	}

	rows, err := self.listMoviesToFix.Query()
	if err != nil {
		mlog.Fatalf("unable to begin transaction: %s", self.err)
	}

	items := make([]*message.Movie, 0)

	self.cnt = 0

	for rows.Next() {
		movie := message.Movie{}
		rows.Scan(&movie.Id, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched)
		items = append(items, &movie)
		self.cnt++
	}
	rows.Close()

	tx.Commit()

	mlog.Info("Listed %d movies to fix", self.cnt)

	msg.Reply <- items
}
