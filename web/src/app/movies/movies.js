angular.module( 'mediabase.movies', [
	'ui.router'
])

.config(function config( $stateProvider ) {
  $stateProvider
  .state('movies', {
      url: '/movies',
      controller: 'MoviesCtrl',
      templateUrl: 'movies/movies.tpl.html',
      data: {pageTitle: 'Movies'}
  });
})

.controller( 'MoviesCtrl', ['$scope', 'core', function MoviesCtrl( $scope, core ) {
	$scope.items = [];

	$scope.selectedIndex = 0;
	$scope.itemClicked = function($index) {
		$scope.selectedIndex = $index;
	}

	$scope.scan = function() {
		core.scanMovies()
	}

	$scope.$onRootScope('app.search', function(selfie, term) {
		core.searchMovies(term)
			.success(function(data, status, headers, config) {
				$scope.items = data;
			});
	})

	core.getMovies()
		.success(function(data, status, headers, config) {
			$scope.items = data;
		});

}])

;