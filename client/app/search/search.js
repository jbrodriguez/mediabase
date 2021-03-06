(function () {
    'use strict';

    angular
        .module('app.search')
        .controller('Search', Search);

    // Search.$inject = ['$q', 'api', 'logger'];

    /* @ngInject */
    function Search($state, $q, $scope, api, logger, options) {

        /*jshint validthis: true */
        var vm = this;

        vm.movies = [];
        vm.options = options;
        vm.options.current = 0;
        vm.options.limit = 50;

        vm.setWatched = setWatched;
        vm.fixMovie = fixMovie;

        // console.log('activated search view');
        $scope.$onRootScope('/movies/search', doSearch);

        function doSearch(me, term) {
            // console.log('searching for me: '+me+'term: '+term+'options.searchTerm: '+options.searchTerm);
            // var args = {current: vm.current, limit: vm.limit, filterBy: options.filterBy, searchTerm: options.searchTerm};
            return api.searchMovies(options).then(function(data) {
                // console.log("what is?: ", data);
                vm.movies = null;
                vm.movies = data;
                return vm.movies;
            })
        };

        function setWatched(index) {
            // console.log("maldecido!!!!: ", index);
            return api.setWatched(vm.movies[index]).then(function(data) {
                var title = vm.movies[index].title;
                logger.success("Movie was updated successfully", "", title);
            })
        };

        function fixMovie(index) {
            return api.fixMovie(vm.movies[index]).then(function(data) {
                logger.success("Movie fixed successfully");
            })
        };
    }
})();