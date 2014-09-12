(function () {
    'use strict';

    angular
        .module('app.recent')
        .config(function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('recent', {
                    url: '/recent',
                    templateUrl: 'app/template/main.html',
                    controller: 'Recent as vm',
                })            
        });

})();