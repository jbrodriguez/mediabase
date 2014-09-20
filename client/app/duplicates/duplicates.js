(function () {
    'use strict';

    angular
        .module('app.duplicates')
        .controller('Duplicates', Duplicates);

    // Duplicates.$inject = ['$q', 'api', 'logger'];

    /* @ngInject */
    function Duplicates($state, $q, api, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.movies = [];

        activate();

        function activate() {
            return getDuplicateMovies().then(function() {
                logger.info('activated duplicates view');
            });
        };

        function getDuplicateMovies() {
            return api.getDuplicateMovies().then(function (data) {
                vm.movies = data;
                return vm.movies;
            });
        };
    }
})();