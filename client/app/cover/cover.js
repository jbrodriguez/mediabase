(function () {
    'use strict';

    angular
        .module('app.cover')
        .controller('Cover', Cover);

    // Cover.$inject = ['$q', 'api', 'logger'];

    /* @ngInject */
    function Cover($state, $q, api, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.movies = [];
        vm.setWatched = setWatched;
        vm.fixMovie = fixMovie;

        activate();

        function activate() {
            return getCover().then(function() {
                logger.info('activated cover view');
            });
        }

        function getCover() {
            return api.getCover().then(function (data) {
                // logger.info('what is: ', data)
                vm.movies = null;
                vm.movies = data.movies;
                return vm.movies;
            });
        }

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