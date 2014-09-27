(function () {
    'use strict';

    angular
        .module('app.all')
        .controller('All', All);

    // All.$inject = ['$q', 'api', 'logger'];

    /* @ngInject */
    function All($q, api, logger, options) {

        /*jshint validthis: true */
        var vm = this;

        vm.movies = [];
        vm.current = 0;
        vm.limit = 50;

        activate();

        function activate() {
            console.log('data: ', vm.current, vm.limit, options.sortBy, options.sortOrder);
            var args = {current: vm.current, limit: vm.limit, sortBy: options.sortBy, sortOrder: options.sortOrder};
            return getAllMovies(args).then(function() {
                logger.info('activated all view');
            });
        } ;

        function getAllMovies(args) {
            console.log('args: ', args.current, args.limit, args.sortBy, args.sortOrder);
            return api.getAllMovies(args).then(function (data) {
                vm.movies = data;
                return vm.movies;
            });
        };
    }
})();