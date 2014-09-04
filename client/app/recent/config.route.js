(function () {
    'use strict';

    angular
        .module('app.recent')
        .config(function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('recent', {
                    url: '/recent',
                    templateUrl: 'app/recent/recent.html',
                    controller: 'Recent as vm',
                })            
        });

})();