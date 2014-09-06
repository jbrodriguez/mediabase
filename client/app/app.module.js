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
    ]);

    angular
        .module('app')
        .config(function($stateProvider, $urlRouterProvider) {
            $urlRouterProvider.otherwise('/recent');
        })
        .config(['$provide', function($provide) {
          $provide.decorator('$rootScope', ['$delegate', function($delegate) {
            $delegate.constructor.prototype.$onRootScope = function(name, listener) {
              var unsubscribe = $delegate.$on(name, listener);
              this.$on('$destroy', unsubscribe);
            };

            return $delegate;
          }]);
        }])        
        .controller('Home', Home)

    /* @ngInject */
    function Home($state, $scope, $rootScope) {

        /*jshint validthis: true */
        var vm = this;

        vm.searchTerm = '';

        $scope.$watch(angular.bind(this, function(searchTerm) {
            return this.searchTerm;
        }), function(newVal) {
             console.log('mothefucker');
            $state.go('search');
            $rootScope.$emit('/local/search', vm.searchTerm)           
        })
    }

})();