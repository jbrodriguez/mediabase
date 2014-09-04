(function () {
    'use strict';

    angular
        .module('app.core')
        .factory('api', api);

    // api.$inject = ['$http', '$location', exception, logger];

    /* @ngInject */
    function api($http, $location, exception, logger) {
    	var ep = "/api/v1/";

    	var service = {
    		getRecentMovies: getRecentMovies
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
    			return data;
    		}
    	}
    }

})();