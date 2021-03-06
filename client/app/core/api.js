(function () {
    'use strict';

    angular
        .module('app.core')
        .factory('api', api);

    // api.$inject = ['$http', '$location', exception, logger];

    /* @ngInject */
    function api($http, $location, exception, logger) {
    	var ep = "/api/v1";

    	var service = {
            getConfig: getConfig,
            saveConfig: saveConfig,
    		getCover: getCover,
            getMovies: getMovies,
            startImport: startImport,
            searchMovies: searchMovies,
            getStatus: getStatus,
            setWatched: setWatched,
            fixMovie: fixMovie,
            getDuplicateMovies: getDuplicateMovies,
            pruneMovies: pruneMovies,
    	};

    	return service;

        function getConfig() {
            return $http.get(ep + '/config')
                .then(getConfigEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for getConfig')(message);
                    $location.url('/');
                });

            function getConfigEnd(data, status, header, config) {
                return data.data
            };
        };

        function saveConfig(arg) {
            return $http.put(ep + '/config', arg)
                .then(saveConfigEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for saveConfig')(message);
                    $location.url('/');
                });

            function saveConfigEnd(data, status, header, config) {
                return data.data
            };
        };

    	function getCover() {
    		return $http.get(ep + '/movies/cover')
    			.then(getCoverEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for getCover')(message);
                    $location.url('/');
                });

    		function getCoverEnd(data, status, headers, config) {
                // logger.info('this is what i got: ', data);
    			return data.data;
    		}
    	};

        function getMovies(args) {
            // console.log('api: ', args.current, args.limit, args.sortBy, args.sortOrder);

            return $http.post(ep + '/movies', args)
                .then(getMoviesEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for getMovies')(message);
                    $location.url('/');
                });

            function getMoviesEnd(data, status, headers, config) {
                return data.data;
            }
        };        

        function startImport() {
            return $http.get(ep + '/import')
                .then(startImportEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for startImport')(message);
                    $location.url('/');
                });

            function startImportEnd(data, status, headers, config) {
                return data.data;
            }
        };

        function searchMovies(args) {
            return $http.post(ep + '/movies/search', args)
                .then(searchMoviesEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for searchMovies')(message);
                    $location.url('/');
                });

            function searchMoviesEnd(data, status, headers, config) {
                return data.data;
            }
        };        

        function getStatus() {
            return $http.get(ep + '/import/status')
                .then(getStatusEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for getStatus')(message);
                    $location.url('/');
                });

            function getStatusEnd(data, status, headers, config) {
                return data.data;
            }
        };

        function setWatched(movie) {
            // convert movie.watched to UTC and save it to last_watched
            if (movie.watched) {
                movie.last_watched = movie.watched.toISOString();
            }
                        
            return $http.put(ep + '/movies/watched', movie)
                .then(setWatchedEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for saveMovie')(message);
                    $location.url('/');
                });

            function setWatchedEnd(data, status, headers, config) {
                return data.data;
            }          
        };

        function fixMovie(movie) {
            return $http.post(ep + '/movies/fix', movie)
                .then(fixMovieEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for fixMovie')(message);
                    $location.url('/');
                });

            function fixMovieEnd(data, status, headers, config) {
                return data.data;
            }          
        };

        function getDuplicateMovies() {
            return $http.get(ep + '/movies/duplicates')
                .then(getDuplicateMoviesEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for getDuplicateMovies')(message);
                    $location.url('/');
                });

            function getDuplicateMoviesEnd(data, status, headers, config) {
                // logger.info('this is what i got: ', data);
                return data.data;
            }
        };

        function pruneMovies() {
            return $http.post(ep + '/movies/prune')
                .then(pruneMoviesEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for pruneMovies')(message);
                    $location.url('/');
                });

            function pruneMoviesEnd(data, status, headers, config) {
                return data.data;
            }          
        };        

    }

})();