(function () {
    'use strict';

    angular
        .module('app.movies')
        .controller('Movies', Movies);

    // Movies.$inject = ['$q', 'api', 'logger'];

    /* @ngInject */
    function Movies($scope, api, logger, options) {

        /*jshint validthis: true */
        var vm = this;

        vm.movies = [];
        vm.current = 0;
        vm.limit = 50;

        vm.busy = false;

        vm.setWatched = setWatched;
        vm.fixMovie = fixMovie;
        vm.scrollPage = scrollPage;

        $scope.$onRootScope('/movies/refresh', doRefresh);

        activate();

        function activate() {
            // // console.log('data: ', vm.current, vm.limit, options.sortBy, options.sortOrder);
            // vm.current = 0;
            // var args = {current: vm.current, limit: vm.limit, sortBy: options.sortBy, sortOrder: options.sortOrder};
            // return getMovies(args).then(function() {
            //     logger.info('activated movies view');
            // });
            doRefresh();
        } ;

        function doRefresh() {
            vm.current = 0;
            var args = {current: vm.current, limit: vm.limit, sortBy: options.sortBy, sortOrder: options.sortOrder};
            return getMovies(args).then(function() {
                logger.info('refreshed list');
            });
        };

        function getMovies(args) {
            // console.log('args: ', args.current, args.limit, args.sortBy, args.sortOrder);
            return api.getMovies(args).then(function (data) {
                vm.movies = null;
                vm.movies = data;
                return vm.movies;
            });
        };

        function setWatched(index) {
            // console.log("maldecido!!!!: ", index);
            return api.setWatched(vm.movies[index]).then(function(data) {
                var title = vm.movies[index].title
                logger.success("Movie was updated successfully", "", title);
            })
        };

        function fixMovie(index) {
            return api.fixMovie(vm.movies[index]).then(function(data) {
                logger.success("Movie fixed successfully");
            })
        };

        function scrollPage() {
            if (vm.busy) return;
            vm.busy = true;
            vm.current += vm.limit;

            var args = {current: vm.current, limit: vm.limit, sortBy: options.sortBy, sortOrder: options.sortOrder};
            return scrollMovies(args).then(function() {
                logger.info('scrolled list');
                vm.busy = false;
            });

            vm.busy = false;
        };

        function scrollMovies(args) {
            return api.getMovies(args).then(function (data) {
                if (vm.current === 0) {
                    vm.movies = null;
                    vm.movies = data;
                } else {
                    for (var i = 0; i < data.length; i++) {
                        vm.movies.push(data[i]);
                    };
                };
                return vm.movies;
            });
        };
    }
})();