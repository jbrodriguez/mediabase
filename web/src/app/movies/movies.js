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
  .state('movies.all', {
      url: '/movies/all',
      controller: 'MoviesCtrl',
      templateUrl: 'movies/movies.all.tpl.html',
      data: {pageTitle: 'All Movies'}
  })
  .state('movies.duplicates', {
      url: '/movies/duplicates',
      controller: 'MoviesCtrl',
      templateUrl: 'movies/movies.duplicates.tpl.html',
      data: {pageTitle: 'Duplicate Movies'}
  })
  .state('movies.runtime', {
      url: '/movies/runtime',
      controller: 'MoviesCtrl',
      templateUrl: 'movies/movies.runtime.tpl.html',
      data: {pageTitle: 'Duplicate Movies'}
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

	$scope.prune = function() {
		core.pruneMovies()
	}

	$scope.list = function() {
		core.listMovies()
			.success(function(data, status, headers, config) {
				$scope.items = data;
			})
	}

	$scope.duplicates = function() {
		core.showDuplicates()
			.success(function(data, status, headers, config) {
				$scope.items = data;
			})
	}	

	$scope.runtime = function() {
		core.listByRuntime()
			.success(function(data, status, headers, config) {
				$scope.items = data;
			})
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