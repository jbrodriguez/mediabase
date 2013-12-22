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
  })
  .state('movies.detail', {
  	  url: '/:index',
  	  controller: 'MoviesDetailCtrl',
  	  templateUrl: 'movies/movies.detail.tpl.html',
  	  data: {pageTitle: 'Movies Detail'}
  });
})

.controller( 'MoviesCtrl', ['$scope', '$state', 'core', function MoviesCtrl( $scope, $state, core ) {
	$scope.items = [];

	$scope.selectedIndex = 0;
	$scope.itemClicked = function($index) {
		$scope.selectedIndex = $index;
		$state.go('movies.detail', {"index": $index})
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

.controller('MoviesDetailCtrl', ['$scope', '$stateParams', 'core', function MoviesDetailCtrl($scope, $stateParams, core) {
	$scope.item = $scope.items[$stateParams.index]

}])

;