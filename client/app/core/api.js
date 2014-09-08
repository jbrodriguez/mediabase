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
    		getRecentMovies: getRecentMovies,
            getAllMovies: getAllMovies,
            startImport: startImport,
            searchMovies: searchMovies,
            getStatus: getStatus
    	};

    	return service;

    	function getRecentMovies() {
    		return $http.get(ep + '/movies')
    			.then(getRecentMoviesEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for getRecentMovies')(message);
                    $location.url('/');
                });

    		function getRecentMoviesEnd(data, status, headers, config) {
                logger.info('this is what i got: ', data);
    			return data.data;
    		}
    	};

        function getAllMovies() {
            return $http.get(ep + '/all')
                .then(getAllMoviesEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for getAllMovies')(message);
                    $location.url('/');
                });

            function getAllMoviesEnd(data, status, headers, config) {
                logger.info('allmovies > this is what i got: ', data);
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

        function searchMovies(term) {
            return $http.get(ep + '/search/' + term)
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
            return $http.get(ep + '/status')
                .then(getStatusEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for getStatus')(message);
                    $location.url('/');
                });

            function getStatusEnd(data, status, headers, config) {
                return data.data;
            }
        };   
    }

})();