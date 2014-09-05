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
    		getRecentMovies: getRecentMovies,
            startScan: startScan,
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
    	}

        function startScan() {
            return $http.get(ep + '/scan')
                .then(startScanEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for startScan')(message);
                    $location.url('/');
                });

            function startScanEnd(data, status, headers, config) {
                return data.data;
            }
        }

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
        }        
    }

})();