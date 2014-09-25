(function () {
    'use strict';

    angular.module('app', [
        /*
         * Order is not important. Angular makes a
         * pass to register all of the modules listed
         * and then when app.dashboard tries to use app.data,
         * it's components are available.
         */

        /*
         * Everybody has access to these.
         * We could place these under every feature area,
         * but this is easier to maintain.
         */
        'app.core',

        /*
         * Feature areas
         */
        'app.recent',
        'app.import',
        'app.search',
        'app.all',
        'app.duplicates',
        'app.settings',
    ]);

    angular
        .module('app')
        // .config(function($stateProvider, $urlRouterProvider, $locationProvider) {
        //     $locationProvider.html5Mode(true);
        //     $urlRouterProvider.otherwise('/recent');
        // })
        .config(['$provide', function($provide) {
          $provide.decorator('$rootScope', ['$delegate', function($delegate) {
            $delegate.constructor.prototype.$onRootScope = function(name, listener) {
              var unsubscribe = $delegate.$on(name, listener);
              this.$on('$destroy', unsubscribe);
            };
            return $delegate;
          }]);
        }])
})();