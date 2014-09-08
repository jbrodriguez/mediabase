(function () {
    'use strict';

    angular
        .module('app.all')
        .config(function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('all', {
                    url: '/all',
                    templateUrl: 'app/template/main.html',
                    controller: 'All as vm',
                })            
        });

})();