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
                vm.movies = data.data;
                return vm.movies;
            });
        }
    }
})();