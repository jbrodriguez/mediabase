(function () {
    'use strict';

    angular
        .module('app.recent')
        .controller('Recent', Recent);

    Recent.$inject = ['$q', 'api', 'logger'];

    function Recent($q, api, logger) {

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

        function setWatched(idx) {
            console.log("maldecido!!!!: ", idx);
            var index = idx;
            return api.setWatched(vm.movies[idx]).then(function(data) {
                console.log('renacuajo!!!: ', vm.movies[index]);
            })
        };        

        function fixMovie(movie) {
            console.log("this is the movie: ", movie.title)
        }
    }
})();