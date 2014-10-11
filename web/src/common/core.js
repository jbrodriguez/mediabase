angular.module( 'mediabase.services', [
])

.factory('core', ['$http', function($http) {
	var api = "/api/v1/";
	var rest = {};

	rest.scanMovies = function() {
		return $http.get(api + "movies/scan");
	};

	rest.pruneMovies = function() {
		return $http.get(api + "movies/prune");
	};

	rest.fixMovies = function() {
		return $http.get(api + "movies/fix");
	};

	rest.getMovies = function() {
		return $http.get(api + "movies");
	};

	rest.listMovies = function() {
		return $http.get(api + "movies/all");
	}

	rest.showDuplicates = function() {
		return $http.get(api + "movies/duplicates");
	}

	rest.listByRuntime = function() {
		return $http.get(api + "movies/runtime");
	}

	rest.searchMovies = function(term) {
		return $http.get(api + "movies/search&q=" + term)
	}

	return rest;
}]);