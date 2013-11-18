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
      templateUrl: 'movies/movies.scan.tpl.html'
  });
})

.controller( 'MoviesCtrl', function MoviesCtrl( $scope ) {
})

.controller( 'MoviesScanCtrl', function MoviesScanCtrl($scope) {
})

;