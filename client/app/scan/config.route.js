(function () {
    'use strict';

    angular
        .module('app.scan')
        .config(function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('scan', {
                    url: '/scan',
                    templateUrl: 'app/scan/scan.html',
                    controller: 'Scan as vm',
                })            
        });

})();