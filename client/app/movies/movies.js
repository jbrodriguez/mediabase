(function () {
    'use strict';

    angular
        .module('app.movies')
        .controller('Movies', Movies);

    // Movies.$inject = ['$q', 'api', 'logger'];

    /* @ngInject */
    function Movies($scope, $window, api, logger, options) {

        /*jshint validthis: true */
        var vm = this;

        vm.movies = [];

        vm.itemsPerPage = 50;
        vm.currentPage = 1;
        vm.current = 0;
        vm.total = 0;

        vm.scrollPage = scrollPage;
        vm.setWatched = setWatched;
        vm.fixMovie = fixMovie;

        $scope.$onRootScope('/movies/refresh', refresh);
        $scope.$onRootScope('/movies/search', search);

        activate();

        function activate() {
            update();
        };

        function refresh() {
            options.mode = 'regular';
            update();
        };

        function search() {
            options.mode = 'search';
            update();
        };

        function update() {
            vm.currentPage = 1;
            scrollPage(vm.currentPage);
        };

        // $scope.$watch(angular.bind(this, function() {
        //     return vm.currentPage;
        // }), function(newVal, oldVal) {
        //     scrollPage(newVal);
        // }, true);

        function scrollPage(page) {
            vm.current = (page - 1) * vm.itemsPerPage;

            $window.scrollTo(0, 0);

            var args = {searchTerm: options.searchTerm, current: vm.current, limit: vm.itemsPerPage, sortBy: options.sortBy, sortOrder: options.sortOrder, filterBy: options.filterBy};

            if (options.mode === 'regular' || options.searchTerm === '') {
                return api.getMovies(args).then(function (data) {
                    vm.movies = null;

                    vm.total = data.count;
                    vm.movies = data.movies;

                    return vm.movies;
                });
            } else {
                return api.searchMovies(args).then(function(data) {
                    // console.log("what is?: ", data);
                    vm.movies = null;

                    vm.total = data.count;
                    vm.movies = data.movies;

                    return vm.movies;
                })                
            };
        };

        function setWatched(index) {
            // console.log("maldecido!!!!: ", index);
            return api.setWatched(vm.movies[index]).then(function(movie) {
                vm.movies[index] = null
                vm.movies[index] = movie
                logger.success("Movie was updated successfully", "", movie.title);
            })
        };

        function fixMovie(index) {
            return api.fixMovie(vm.movies[index]).then(function(data) {
                logger.success("Movie fixed successfully");
            })
        };

   

        // function scrollPage() {
        //     if (vm.busy) return;
        //     vm.busy = true;
        //     vm.current += vm.limit;

        //     var args = {current: vm.current, limit: vm.limit, sortBy: options.sortBy, sortOrder: options.sortOrder};
        //     return scrollMovies(args).then(function() {
        //         logger.info('scrolled list');
        //         vm.busy = false;
        //     });

        //     vm.busy = false;
        // };

        // function scrollMovies(args) {
        //     return api.getMovies(args).then(function (data) {
        //         if (vm.current === 0) {
        //             vm.movies = null;
        //             vm.movies = data;
        //         } else {
        //             for (var i = 0; i < data.length; i++) {
        //                 vm.movies.push(data[i]);
        //             };
        //         };
        //         return vm.movies;
        //     });
        // };
    }
})();