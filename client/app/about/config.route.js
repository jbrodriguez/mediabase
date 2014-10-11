(function () {
    'use strict';

    angular
        .module('app.about')
        .config(function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('about', {
                    url: '/about',
                    templateUrl: 'app/about/about.html',
                    controller: 'About as vm',
                })            
        });

})();