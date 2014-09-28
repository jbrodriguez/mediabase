(function () {
    'use strict';

    angular
        .module('app.cover')
        .config(function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('cover', {
                    url: '/cover',
                    templateUrl: 'app/cover/cover.html',
                    controller: 'Cover as vm',
                })            
        });

})();