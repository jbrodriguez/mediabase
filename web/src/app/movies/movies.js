angular.module( 'mediabase.movies', [
	'ui.state',
	'placeholders',
	'ui.bootstrap'
])

.config(function config( $stateProvider ) {
  $stateProvider
  .state('movies', {
      url: '/movies',
      controller: 'MoviesCtrl',
      templateUrl: 'movies/movies.tpl.html',
      data: {pageTitle: 'Movies'}
  })
  .state('movies.scan', {
      url: '/scan',
      controller: 'MoviesScanCtrl',
      templateUrl: 'movies/movies.scan.tpl.html',
      data: {pageTitle: 'Scan Movies'}
  });
})

.controller( 'MoviesCtrl', ['$scope', 'core', function MoviesCtrl( $scope, core ) {
	$scope.message = 'Initializing ...';
	$scope.items = [];

	core.getMovies()
		.success(function(data, status, headers, config) {
			$scope.message = "I am el mero mero";
			$scope.items = data;
		});

}])

.controller( 'MoviesScanCtrl', ['$scope', 'core', function MoviesScanCtrl($scope, core) {
	$scope.message = "we're the waiting";

	core.scanMovies()
		.then(function() {
			$scope.message = "you broke another mirror";
			// $scope.apply();
		});
}])

;