package services

import (
	"apertoire.net/mediabase/server/bus"
	"apertoire.net/mediabase/server/message"
	"apertoire.net/mediabase/server/model"
	"database/sql"
	"fmt"
	"github.com/apertoire/mlog"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	"strings"
	"time"
)

type Dal struct {
	Bus         *bus.Bus
	Config      *model.Config
	db          *sql.DB
	dbase       string
	err         error
	count       uint64
	searchCount uint64
	searchArgs  string

	countRows  *sql.Stmt
	storeMovie *sql.Stmt
	// searchMovies    *sql.Stmt
	// searchGenre     *sql.Stmt
	// searchCountry   *sql.Stmt
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

	self.dbase = filepath.Join(".", "db", "mediabase.db")
	self.db, self.err = sql.Open("sqlite3", self.dbase)
	if self.err != nil {
		mlog.Fatalf("open database: %s (%s)", self.err, self.dbase)
	}

	stmtExist := self.prepare(`select name from sqlite_master where type='table' and name='movie'`)
	defer stmtExist.Close()

	var name string
	err := stmtExist.QueryRow().Scan(&name)
	if err != nil {
		mlog.Fatalf("unable to check for existence of movie database: %s (%s)", self.err, self.dbase)
	}

	if name != "movie" {
		mlog.Info("Initializing database schema ...")
		self.initSchema()
	}

	self.count = 0
	self.searchCount = 0
	self.searchArgs = ""

	self.countRows = self.prepare("select count(*) from movie;")
	self.listMovies = self.prepare("select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched, score, director, writer, actors, awards, imdb_rating, imdb_votes from movie order by ? desc limit ? offset ?")
	self.listByRuntime = self.prepare("select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched, score, director, writer, actors, awards, imdb_rating, imdb_votes from movie order by runtime")
	self.listMoviesToFix = self.prepare("select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched, score, director, writer, actors, awards, imdb_rating, imdb_votes from movie where original_title = 'FIXMOV23'")

	var abs string
	if abs, err = filepath.Abs(self.dbase); err != nil {
		mlog.Info("unable to get absolute path: %s, ", err)
		return
	}

	mlog.Info("connected to database (%s)", abs)

	// self.initSchema()

	go self.react()
}

func (self *Dal) Stop() {
	self.listMoviesToFix.Close()
	self.listByRuntime.Close()
	self.listMovies.Close()
	// self.searchMovies.Close()
	// self.storeMovie.Close()
	self.db.Close()

	mlog.Info("dal service stopped")
}

func (self *Dal) react() {
	for {
		select {
		case msg := <-self.Bus.GetCover:
			go self.doGetCover(msg)
		case msg := <-self.Bus.GetMovies:
			go self.doGetMovies(msg)

		case msg := <-self.Bus.StoreMovie:
			self.doStoreMovie(msg)
		case msg := <-self.Bus.DeleteMovie:
			self.doDeleteMovie(msg)
		case msg := <-self.Bus.UpdateMovie:
			self.doUpdateMovie(msg)
		case msg := <-self.Bus.ShowDuplicates:
			go self.doShowDuplicates(msg)
		case msg := <-self.Bus.SearchMovies:
			go self.doSearchMovies(msg)
		case msg := <-self.Bus.CheckMovie:
			go self.doCheckMovie(msg)
		// case msg := <-self.Bus.GetMoviesToFix:
		// 	go self.doGetMoviesToFix(msg)
		case msg := <-self.Bus.WatchedMovie:
			go self.doWatchedMovie(msg)
		}
	}
}

func (self *Dal) ConfigChanged(conf *model.Config) {
	self.Config = conf
}

func (self *Dal) doCheckMovie(msg *message.CheckMovie) {
	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare("select rowid from movie where upper(location) = ?")
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at prepare: %s", err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(strings.ToUpper(msg.Movie.Location)).Scan(&id)

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
	self.count = 0

	mlog.Info("STARTED SAVING %s [%d]", movie.Title)

	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	// stmt, err := tx.Prepare("insert into movie(title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, director, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	stmt, err := tx.Prepare("insert into movie(title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched, score, director, writer, actors, awards, imdb_rating, imdb_votes) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at prepare: %s", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(movie.Title, movie.Original_Title, movie.File_Title, movie.Year, movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline, movie.Resolution, movie.FileType, movie.Location, movie.Cover, movie.Backdrop,
		movie.Genres, movie.Vote_Average, movie.Vote_Count, movie.Production_Countries, movie.Added, movie.Modified, movie.Last_Watched, movie.All_Watched, movie.Count_Watched, movie.Score, movie.Director, movie.Writer, movie.Actors, movie.Awards,
		movie.Imdb_Rating, movie.Imdb_Votes)
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
	mlog.Info("FINISHED SAVING %s", movie.Title)

	// _, self.err = self.storeMovie.Exec(movie.Name, movie.Year, movie.Resolution, movie.Type, movie.Path, movie.Picture)
	// if self.err != nil {
	// 	mlog.Fatalf(self.err)
	// }
}

func (self *Dal) doDeleteMovie(movie *message.Movie) {
	self.count = 0

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
	mlog.Info("STARTED UPDATING %s", movie.Title)

	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare("update movie set title = ?, original_title = ?, year = ?, runtime = ?, tmdb_id = ?, imdb_id = ?, overview = ?, tagline = ?, cover = ?, backdrop = ?, genres = ?, vote_average = ?, vote_count = ?, countries = ?, modified = ?, director = ?, writer = ?, actors = ?, awards = ?, imdb_rating = ?, imdb_votes = ? where rowid = ?")
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at prepare: %s", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(movie.Title, movie.Original_Title, movie.Year, movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline, movie.Cover, movie.Backdrop, movie.Genres, movie.Vote_Average, movie.Vote_Count, movie.Production_Countries, movie.Modified, movie.Director, movie.Writer, movie.Actors, movie.Awards, movie.Imdb_Rating, movie.Imdb_Votes, movie.Id)
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at exec: %s", err)
	}

	tx.Commit()
	mlog.Info("FINISHED UPDATING %s", movie.Title)
}

func (self *Dal) doGetCover(msg *message.Movies) {
	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("unable to begin transaction: %s", err)
	}

	stmt, err := tx.Prepare("select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched, score, director, writer, actors, awards, imdb_rating, imdb_votes from movie order by added desc limit ?")
	if err != nil {
		mlog.Fatalf("unable to prepare transaction: %s", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(60)
	if err != nil {
		mlog.Fatalf("unable to run transaction: %s", err)
	}

	items := make([]*message.Movie, 0)

	for rows.Next() {
		movie := message.Movie{}
		err := rows.Scan(&movie.Id, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors, &movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes)
		if err != nil {
			mlog.Info("errored: %s", err)
		}
		// mlog.Info("%+v\n", movie)

		items = append(items, &movie)
	}
	rows.Close()

	tx.Commit()

	// mlog.Info("got back %+v", items)

	msg.Reply <- &message.MoviesDTO{Movies: items}
}

func (self *Dal) doGetMovies(msg *message.Movies) {
	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("unable to begin transaction: %s", err)
	}

	options := msg.Options
	mlog.Info("what is: %+v", options)

	stmt, err := tx.Prepare(fmt.Sprintf("select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched, score, director, writer, actors, awards, imdb_rating, imdb_votes from movie order by %s %s limit ? offset ?", options.SortBy, options.SortOrder))
	if err != nil {
		mlog.Fatalf("unable to prepare transaction: %s", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(options.Limit, options.Current)
	if err != nil {
		mlog.Fatalf("unable to prepare transaction: %s", self.err)
	}

	items := make([]*message.Movie, 0)

	if self.count == 0 {
		err = self.countRows.QueryRow().Scan(&self.count)
		if err != nil {
			mlog.Fatalf("unable to count rows: %s", err)
		}
	}

	var count = 0
	for rows.Next() {
		movie := message.Movie{}
		rows.Scan(&movie.Id, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors, &movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes)
		items = append(items, &movie)
		count++
	}
	rows.Close()

	tx.Commit()

	mlog.Info("Listed %d movies", count)
	mlog.Info("Representing %d movies", self.count)

	msg.Reply <- &message.MoviesDTO{Count: self.count, Movies: items}
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
	rows, err := self.db.Query("select a.rowid, a.title, a.original_title, a.file_title, a.year, a.runtime, a.tmdb_id, a.imdb_id, a.overview, a.tagline, a.resolution, a.filetype, a.location, a.cover, a.backdrop, a.genres, a.vote_average, a.vote_count, a.countries, a.added, a.modified, a.last_watched, a.all_watched, a.count_watched, a.score, a.director, a.writer, a.actors, a.awards, a.imdb_rating, a.imdb_votes from movie a join (select title, year from movie group by title, year having count(*) > 1) b on a.title = b.title and a.year = b.year;")
	if err != nil {
		mlog.Fatalf("unable to prepare transaction: %s", self.err)
	}

	items := make([]*message.Movie, 0)

	var count uint64 = 0

	for rows.Next() {
		movie := message.Movie{}
		rows.Scan(&movie.Id, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors, &movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes)
		items = append(items, &movie)
		count++
	}
	rows.Close()

	tx.Commit()

	mlog.Info("Found %d duplicate movies", count)

	msg.Reply <- &message.MoviesDTO{Count: count, Movies: items}
}

func (self *Dal) doSearchMovies(msg *message.Movies) {
	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("unable to begin transaction: %s", err)
	}

	term := msg.Options.SearchTerm + "*"
	args := msg.Options.FilterBy

	mlog.Info("this is: %s %s %s", term, self.searchArgs, args)

	// if self.searchArgs != args {
	// self.searchArgs = args

	search := fmt.Sprintf(`select count(*) from movie dt, %s vt where vt.%s match ? and dt.rowid = vt.docid;`, "movie"+msg.Options.FilterBy, "movie"+msg.Options.FilterBy)

	stmt, err := tx.Prepare(search)
	if err != nil {
		mlog.Fatalf("unable to prepare transaction: %s", err)
	}
	defer stmt.Close()

	mlog.Info("sup dude %s", search)

	err = stmt.QueryRow(term).Scan(&self.searchCount)
	if err != nil {
		mlog.Fatalf("unable to count rows: %s", err)
	}
	// }

	sql := fmt.Sprintf(`select dt.rowid, dt.title, dt.original_title, dt.year, dt.runtime, dt.tmdb_id, dt.imdb_id, dt.overview, dt.tagline, dt.resolution,
				dt.filetype, dt.location, dt.cover, dt.backdrop, dt.genres, dt.vote_average, dt.vote_count, dt.countries, dt.added, dt.modified, 
				dt.last_watched, dt.all_watched, dt.count_watched, dt.score, dt.director, dt.writer, dt.actors, dt.awards, dt.imdb_rating, dt.imdb_votes
				from movie dt, %s vt where vt.%s match ? and dt.rowid = vt.docid order by dt.%s %s limit ? offset ?`,
		"movie"+msg.Options.FilterBy, "movie"+msg.Options.FilterBy, msg.Options.SortBy, msg.Options.SortOrder)

	// mlog.Info("my main man: %s", sql)

	stmt, err = tx.Prepare(sql)
	if err != nil {
		mlog.Fatalf("unable to prepare transaction: %s", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(term, msg.Options.Limit, msg.Options.Current)
	if err != nil {
		mlog.Fatalf("unable to begin transaction: %s", self.err)
	}

	items := make([]*message.Movie, 0)

	for rows.Next() {
		movie := message.Movie{}
		rows.Scan(&movie.Id, &movie.Title, &movie.Original_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors, &movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes)
		// movie := &message.Movie{}
		// rows.Scan(movie.Id, movie.Title, movie.Original_Title, movie.Year, movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline, movie.Resolution, movie.FileType, movie.Location, movie.Cover, movie.Backdrop)
		// mlog.Info("title: (%s)", movie.Title)
		items = append(items, &movie)
	}
	rows.Close()

	tx.Commit()

	mlog.Info("Representing %d movies", self.searchCount)

	msg.Reply <- &message.MoviesDTO{Count: self.searchCount, Movies: items}
}

func (self *Dal) doWatchedMovie(msg *message.SingleMovie) {
	mlog.Info("STARTED UPDATING WATCHED MOVIE %s (%s)", msg.Movie.Title, msg.Movie.Last_Watched)

	tx, err := self.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare("update movie set last_watched = ?, all_watched = ?, count_watched = ?, score = ?, modified = ? where rowid = ?")
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at prepare: %s", err)
	}
	defer stmt.Close()

	now := time.Now().UTC().Format(time.RFC3339)

	var all_watched string
	count_watched := msg.Movie.Count_Watched
	if !strings.Contains(msg.Movie.All_Watched, msg.Movie.Last_Watched) {
		count_watched++
		if msg.Movie.All_Watched == "" {
			all_watched = msg.Movie.Last_Watched
		} else {
			all_watched += "|" + msg.Movie.Last_Watched
		}
	}

	_, err = stmt.Exec(msg.Movie.Last_Watched, all_watched, count_watched, msg.Movie.Score, now, msg.Movie.Id)
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at exec: %s", err)
	}

	tx.Commit()
	mlog.Info("FINISHED UPDATING WATCHED MOVIE %s", msg.Movie.Title)

	msg.Movie.All_Watched = all_watched
	msg.Movie.Count_Watched = count_watched
	msg.Movie.Modified = now

	msg.Reply <- msg.Movie
}

func (self *Dal) initSchema() {
	sql := `
DROP TABLE IF EXISTS movie;
DROP TABLE IF EXISTS movietitle;
DROP TABLE IF EXISTS moviegenre;
DROP TABLE IF EXISTS moviecountry;
DROP TABLE IF EXISTS moviedirector;
DROP TABLE IF EXISTS movieactor;

DROP INDEX IF EXISTS movie_filetype_idx;
DROP INDEX IF EXISTS movie_location_idx;
DROP INDEX IF EXISTS movie_title_idx;

DROP TRIGGER IF EXISTS movie_ai;
DROP TRIGGER IF EXISTS movie_au;
DROP TRIGGER IF EXISTS movie_bd;
DROP TRIGGER IF EXISTS movie_bu;

DROP TRIGGER IF EXISTS genre_ai;
DROP TRIGGER IF EXISTS genre_au;
DROP TRIGGER IF EXISTS genre_bd;
DROP TRIGGER IF EXISTS genre_bu;

DROP TRIGGER IF EXISTS country_ai;
DROP TRIGGER IF EXISTS country_au;
DROP TRIGGER IF EXISTS country_bd;
DROP TRIGGER IF EXISTS country_bu;

DROP TRIGGER IF EXISTS director_ai;
DROP TRIGGER IF EXISTS director_au;
DROP TRIGGER IF EXISTS director_bd;
DROP TRIGGER IF EXISTS director_bu;

DROP TRIGGER IF EXISTS actor_ai;
DROP TRIGGER IF EXISTS actor_au;
DROP TRIGGER IF EXISTS actor_bd;
DROP TRIGGER IF EXISTS actor_bu;

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
  score integer,
  director text,
  writer text,
  actors text,
  awards text,
  imdb_rating integer,
  imdb_votes integer
);
CREATE INDEX movie_title_idx ON movie (title);
CREATE INDEX movie_location_idx ON movie (location);
CREATE INDEX movie_filetype_idx ON movie (filetype);

/* titles */
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

/* genres */
CREATE VIRTUAL TABLE moviegenre USING fts4(content="movie", genres);
CREATE TRIGGER genre_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM moviegenre WHERE docid=old.rowid;
END;

CREATE TRIGGER genre_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM moviegenre WHERE docid=old.rowid;
END;

CREATE TRIGGER genre_au AFTER UPDATE ON movie BEGIN
  INSERT INTO moviegenre(docid, genres) VALUES (new.rowid, new.genres);
END;

CREATE TRIGGER genre_ai AFTER INSERT ON movie BEGIN
  INSERT INTO moviegenre(docid, genres) VALUES (new.rowid, new.genres);
END;

/* country */
CREATE VIRTUAL TABLE moviecountry USING fts4(content="movie", countries);
CREATE TRIGGER country_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM moviecountry WHERE docid=old.rowid;
END;

CREATE TRIGGER country_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM moviecountry WHERE docid=old.rowid;
END;

CREATE TRIGGER country_au AFTER UPDATE ON movie BEGIN
  INSERT INTO moviecountry(docid, countries) VALUES (new.rowid, new.countries);
END;

CREATE TRIGGER country_ai AFTER INSERT ON movie BEGIN
  INSERT INTO moviecountry(docid, countries) VALUES (new.rowid, new.countries);
END;

/* director */
CREATE VIRTUAL TABLE moviedirector USING fts4(content="movie", director);
CREATE TRIGGER director_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM moviedirector WHERE docid=old.rowid;
END;

CREATE TRIGGER director_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM moviedirector WHERE docid=old.rowid;
END;

CREATE TRIGGER director_au AFTER UPDATE ON movie BEGIN
  INSERT INTO moviedirector(docid, director) VALUES (new.rowid, new.director);
END;

CREATE TRIGGER director_ai AFTER INSERT ON movie BEGIN
  INSERT INTO moviedirector(docid, director) VALUES (new.rowid, new.director);
END;

/* actor */
CREATE VIRTUAL TABLE movieactor USING fts4(content="movie", actors);
CREATE TRIGGER actor_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM movieactor WHERE docid=old.rowid;
END;

CREATE TRIGGER actor_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM movieactor WHERE docid=old.rowid;
END;

CREATE TRIGGER actor_au AFTER UPDATE ON movie BEGIN
  INSERT INTO movieactor(docid, actors) VALUES (new.rowid, new.actors);
END;

CREATE TRIGGER actor_ai AFTER INSERT ON movie BEGIN
  INSERT INTO movieactor(docid, actors) VALUES (new.rowid, new.actors);
END;
	`

	_, err := self.db.Exec(sql)
	if err != nil {
		mlog.Info("%q: %s", err, sql)
		return
	}

	mlog.Info("inited schema")
}
