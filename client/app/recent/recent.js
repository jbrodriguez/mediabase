(function () {
    'use strict';

    angular
        .module('app.recent')
        .controller('Recent', Recent);

    // Recent.$inject = ['$q', 'api', 'logger'];

    /* @ngInject */
    function Recent($state, $q, api, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.movies = [];
        vm.setWatched = setWatched;
        vm.fixMovie = fixMovie;

        activate();

        function activate() {
            return getRecentMovies().then(function() {
                logger.info('activated recent view');
            });
//             var promises = [getAvengerCount(), getAvengersCast()];
// //            Using a resolver on all routes or dataservice.ready in every controller
// //            return dataservice.ready(promises).then(function(){
//             return $q.all(promises).then(function(){
//                 logger.info('Activated Dashboard View');
//             });
        }

        function getRecentMovies() {
            return api.getRecentMovies().then(function (data) {
                logger.info('what is: ', data)
                vm.movies = data;
                return vm.movies;
            });
        }

        function setWatched(index) {
            console.log("maldecido!!!!: ", index);
            return api.setWatched(vm.movies[index]).then(function(data) {
                logger.success("Movie was updated successfully", "", vm.movies[index].title);
            })
        };

        function fixMovie(index) {
            return api.fixMovie(vm.movies[index]).then(function(data) {
                logger.success("Movie fixed successfully");
                $state.go("recent")
            })
        };
    }
})();