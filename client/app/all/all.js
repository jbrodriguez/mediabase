(function () {
    'use strict';

    angular
        .module('app.all')
        .controller('All', All);

    // All.$inject = ['$q', 'api', 'logger'];

    /* @ngInject */
    function All($q, api, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.movies = [];

        activate();

        function activate() {
            return getAllMovies().then(function() {
                logger.info('activated all view');
            });
        } ;

        function getAllMovies() {
            return api.getAllMovies().then(function (data) {
                vm.movies = data;
                return vm.movies;
            });
        };
    }
})();