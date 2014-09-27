(function () {
    'use strict';

    angular
        .module('app.all')
        .controller('All', All);

    // All.$inject = ['$q', 'api', 'logger'];

    /* @ngInject */
    function All($scope, api, logger, options) {

        /*jshint validthis: true */
        var vm = this;

        vm.movies = [];
        vm.current = 0;
        vm.limit = 50;

        vm.setWatched = setWatched;
        vm.fixMovie = fixMovie;


        $scope.$onRootScope('/movies/refresh', doRefresh);

        activate();

        function activate() {
            console.log('data: ', vm.current, vm.limit, options.sortBy, options.sortOrder);
            var args = {current: vm.current, limit: vm.limit, sortBy: options.sortBy, sortOrder: options.sortOrder};
            return getAllMovies(args).then(function() {
                logger.info('activated all view');
            });
        } ;

        function doRefresh() {
            var args = {current: vm.current, limit: vm.limit, sortBy: options.sortBy, sortOrder: options.sortOrder};
            return getAllMovies(args).then(function() {
                logger.info('refreshed list');
            });
        };

        function getAllMovies(args) {
            console.log('args: ', args.current, args.limit, args.sortBy, args.sortOrder);
            return api.getAllMovies(args).then(function (data) {
                vm.movies = data;
                return vm.movies;
            });
        };

        function setWatched(index) {
            console.log("maldecido!!!!: ", index);
            return api.setWatched(vm.movies[index]).then(function(data) {
                logger.success("Movie was updated successfully", "", vm.movies[index].title);
            })
        };

        function fixMovie(index) {
            return api.fixMovie(vm.movies[index]).then(function(data) {
                logger.success("Movie fixed successfully");
            })
        };      
    }
})();