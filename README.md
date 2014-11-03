Mediabase
=========

*tl;dr* **mediabase** is a proof-of-concept application to catalogue a media library consisting of movies. It scans the folders you choose looking for movies, then fetch metadata from [themoviedb.org](www.themoviedb.org) and [The OMDB API](www.omdbapi.com) and present the information in a nice web page.

Check the [this blog post](http://www.apertoire.net/introducing-mediabase) for a general description of the app.

## Install Guide (End Users)

Please take the following steps, which assume that your home folder is named "MyUser" (/Users/MyUser in MAC OS X or /home/MyUser in Linux)

- Download the zip file containing the binary release
<pre><code>https://github.com/apertoire/mediabase/releases/latest
</code></pre>
- Create a folder (.mediabase) in your home folder.
<pre><code>mkdir /Users/MyUser/.mediabase
</code></pre>
- Cd into this folder and unzip the binary file
<pre><code>cd /Users/MyUser/.mediabase
unzip /Users/MyUser/Downloads/mediabase.zip
</code></pre>
- Run the server
<pre><code>./mediabase
</code></pre>

The server will listen on port 3267 by default, so you can now open a web browser and point it to the app url
<pre><code>http://localhost:3267/
</code></pre>
> **Note:**
> - Both OSX and Linux 64-bit binaries are available.

## Contributing (Developers)

Fork and clone the repo to your drive, y download dependencies
<pre><code>go get ./...
</code></pre>
and
<pre><code>make build
</code></pre>
to create an executable at ./dist. It will also copy the client code and assets to this folder too.

To run the app do
<pre><code>make run
</code></pre>

Please note that go-sqlite3 needs to be built with FTS4 support. So you need to issue
<pre><code>CGO_CFLAGS=-DSQLITE_ENABLE_FTS4 go get -u github.com/mattn/go-sqlite3
</code></pre>


## Credits

 - [Go](https://golang.org/)
 - [AngularJS](https://angularjs.org/)
 - [Foundation](http://foundation.zurb.com/)
 - [Sqlite](http://www.sqlite.org/)
 - [workpool (Ardan Studios)](https://github.com/goinggo/workpool)
 - [resize (Jan Schlicht)](https://github.com/nfnt/resize)
 - [fsm (looplab)](https://github.com/looplab/fsm)
 - [go-sqlite3 (Yasuhiro Matsumoto)](https://github.com/mattn/go-sqlite3)
 - [gin (gin-gonic)](https://github.com/gin-gonic/gin)
 - [semver (fsaintjacques)](https://github.com/fsaintjacques/semver-tool)
 - [go-tmdb (rharter)](https://github.com/rharter/go-tmdb)
 - [go-log (siddontang)](https://github.com/siddontang/go-log)
 - [napping (jmcvetta)](https://github.com/jmcvetta/napping)

## License
[MIT license](http://jbrodriguez.mit-license.org)