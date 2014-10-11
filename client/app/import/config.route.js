(function () {
    'use strict';

    angular
        .module('app.import')
        .config(function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('import', {
                    url: '/import',
                    templateUrl: 'app/import/import.html',
                    controller: 'Import as vm',
                })            
        });

})();