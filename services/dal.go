package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path/filepath"
)

type Dal struct {
	Bus    *bus.Bus
	Config *helper.Config
	db     *sql.DB
	err    error

	storeMovie *sql.Stmt
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
		log.Fatal(err)
	}
	return stmt
}

func (self *Dal) Start() {
	log.Printf("starting dal service ...")

	self.db, self.err = sql.Open("sqlite3", filepath.Join(self.Config.AppDir, "/db/mediabase.db"))
	if self.err != nil {
		log.Fatal(self.err)
	}

	// self.exists = self.prepare("select id from item where name = ?")
	self.storeMovie = self.prepare("insert or ignore into movie (name, year, resolution, filetype, location, picture) values (?, ?, ?, ?, ?, ?)")

	// self.authenticate = self.prepare("select id, password from account where email = $1")
	// self.getUserDataById = self.prepare("select name, email from account where id = $1")
	// getAssets, err := self.db.Prepare("select asset.id, asset.name, asset.category, asset.created, asset.modified, assetCategory.name as categoryName from asset, assetCategory where asset.account_id = ? and assetCategory.id = asset.category order by asset.created desc")
	// getRevisions, err := self.db.Prepare("select id, asset_id, index, created from revision where asset_id = ? order by index desc")
	// getItems, err := self.db.Prepare("select prod.id, prod.name, prod.asin, prod.upc, it.id as typeId, it.name as typeName,  itm.quantity, itm.price, itm.reference from itemtype it, product prod, item itm, asset ast, revision rev where ast.id = ? and rev.id = ? and itm.revision_id = rev.id and prod.id = itm.product_id and it.id = prod.itemtype_id order by prod.itemtype_id asc")
	// getCategories, err := self.db.Prepare("select id, name from itemtype")
	// getLastRevision, err := self.db.Prepare("select max(index) as index from revision where asset_id = ?")
	// getAsset, err := self.db.Prepare("select asset.id, asset.name, asset.category, asset.created, asset.modified, assetCategory.name as categoryName from asset, assetCategory where asset.id = ? and asset.account_id = ? and assetCategory.id = asset.category order by asset.created desc")

	// putAsset, err := self.db.Prepare("insert into asset (account_id, name, category, created, modified) values (?, ?, ?, ?, ?)")
	// putProduct, err := self.db.Prepare("insert into product (itemtype_id, name, asin, sku, upc, ean) VALUES (?, ?, ?, ?, ?, ?)")
	// putRevision, err := self.db.Prepare("insert into revision (asset_id, index, created) values (?, ?, ?)")
	// putItem, err := self.db.Prepare("insert into item (revision_id, product_id, reference, quantity, price) values (?, ?, ?, ?, ?)")

	log.Printf("connected to database")

	go self.react()
}

func (self *Dal) Stop() {
	self.storeMovie.Close()
	self.db.Close()
}

func (self *Dal) react() {
	for {
		select {
		case msg := <-self.Bus.StoreMovie:
			go self.doStoreMovie(msg)
		}
	}
}

func (self *Dal) doStoreMovie(movie *message.Movie) {
	_, self.err = self.storeMovie.Exec(movie.Name, movie.Year, movie.Resolution, movie.Type, movie.Path, movie.Picture)
	if self.err != nil {
		log.Fatal(self.err)
	}
}

// func (self *Dal) doAuthenticate(user *model.UserAuthReq, reply chan *model.UserAuthRep) {
// 	var id int8
// 	var pwd string

// 	err := self.authenticate.QueryRow(user.Email).Scan(&id, &pwd)
// 	if err == sql.ErrNoRows {
// 		reply <- nil
// 		return
// 	} else if err != nil {
// 		panic(err.Error())
// 	}

// 	reply <- &model.UserAuthRep{id, pwd}
// }

// func (self *Dal) doGetUserDataById(user *model.UserDataReq, reply chan *model.UserDataRep) {
// 	var name string
// 	var email string

// 	err := self.getUserDataById.QueryRow(user.Id).Scan(&name, &email)
// 	if err == sql.ErrNoRows {
// 		reply <- nil
// 		return
// 	} else if err != nil {
// 		panic(err.Error())
// 	}

// 	reply <- &model.UserDataRep{name, email}
// }

// func (self *Dal) search() {
// 	rows, err := self.db.Query("SELECT id, heroku_id FROM resources WHERE destroyed_at IS NULL")
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()
// }
