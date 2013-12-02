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

	core.getMovies()
		.success(function(data, status, headers, config) {
			$scope.items = data;
		});

}])

;