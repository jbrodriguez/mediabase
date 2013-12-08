angular.module( 'mediabase.services', [
])

.factory('core', ['$http', function($http) {
	var api = "/api/v1/";
	var rest = {};

	rest.scanMovies = function() {
		return $http.get(api + "movies/scan");
	};

	rest.getMovies = function() {
		return $http.get(api + "movies");
	};

	rest.searchMovies = function(term) {
		return $http.get(api + "movies/search&q=" + term)
	}

	return rest;
}]);