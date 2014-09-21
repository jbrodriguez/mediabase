(function () {
    'use strict';

    angular
        .module('app')
        .controller('Home', Home)

    /* @ngInject */
    function Home($state, $scope, $rootScope, api, options) {

        /*jshint validthis: true */
        var vm = this;

        vm.options = options;
        vm.pruneMovies = pruneMovies;

        $scope.$watch(angular.bind(this, function(searchTerm) {
            return vm.options.searchTerm;
        }), function(newVal) {
             console.log('searching for either vm.searchTerm: '+vm.options.searchTerm + ' or newVal: '+newVal);
            $state.go('search').then(function(current) {
                console.log('after search state');
                $rootScope.$emit('/local/search', newVal);
                console.log('emitted event');
            });
        })

        function pruneMovies() {
            return api.pruneMovies().then(function(data) {
                console.log("are you ready for the fallout?");
                $state.go("recent");
            })
        }
    };

})();