(function () {
    'use strict';

    angular
        .module('app.movies')
        .config(function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('movies', {
                    url: '/movies',
                    templateUrl: 'app/movies/movies.html',
                    controller: 'Movies as vm',
                })            
        });

})();