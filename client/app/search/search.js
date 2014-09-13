(function () {
    'use strict';

    angular
        .module('app.search')
        .controller('Search', Search);

    // Search.$inject = ['$q', 'api', 'logger'];

    /* @ngInject */
    function Search($state, $q, $scope, api, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.movies = [];
        vm.setWatched = setWatched;
        vm.fixMovie = fixMovie;

        console.log('activated search view');
        $scope.$onRootScope('/local/search', doSearch);

        function doSearch(me, term) {
            console.log('searching for me: '+me+'term: '+term);
            return api.searchMovies(term).then(function(data) {
                vm.movies = data;
                return vm.movies;
            })
        };

        function setWatched(index) {
            console.log("maldecido!!!!: ", index);
            return api.setWatched(vm.movies[index]).then(function(data) {
                logger.success("Movie was updated successfully", "", vm.movies[index].title);
            })
        };

        function fixMovie(index) {
            var movie = vm.movies[index];
            if (!movie.tmdbid_new) {
                return
            }

            movie.tmdb_id = movie.tmdbid_new

            return api.fixMovie(movie).then(function(data) {
                logger.success("Movie fixed successfully");
                // $state.go("recent")
            })
        };

        // function changeDate(movie, watchedDate) {
        //     console.log("yist: ", watchedDate);

        // };

//         function activate() {
//             return getRecentMovies().then(function() {
//                 logger.info('activated recent view');
//             });
// //             var promises = [getAvengerCount(), getAvengersCast()];
// // //            Using a resolver on all routes or dataservice.ready in every controller
// // //            return dataservice.ready(promises).then(function(){
// //             return $q.all(promises).then(function(){
// //                 logger.info('Activated Dashboard View');
// //             });
//         }

//         function getRecentMovies() {
//             return api.getRecentMovies().then(function (data) {
//                 logger.info('what is: ', data)
//                 vm.movies = data;
//                 return vm.movies;
//             });
//         }
    }
})();