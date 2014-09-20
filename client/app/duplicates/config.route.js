(function () {
    'use strict';

    angular
        .module('app.duplicates')
        .config(function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('duplicates', {
                    url: '/duplicates',
                    templateUrl: 'app/template/main.html',
                    controller: 'Duplicates as vm',
                })            
        });

})();