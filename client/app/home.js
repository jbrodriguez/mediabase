(function () {
    'use strict';

    angular
        .module('app')
        .controller('Home', Home)

    /* @ngInject */
    function Home($state, $scope, $rootScope) {

        /*jshint validthis: true */
        var vm = this;

        vm.searchTerm = '';

        $scope.$watch(angular.bind(this, function(searchTerm) {
            return vm.searchTerm;
        }), function(newVal) {
             console.log('searching for either vm.searchTerm: '+vm.searchTerm + ' or newVal: '+newVal);
            $state.go('search').then(function(current) {
                console.log('after search state');
                $rootScope.$emit('/local/search', newVal);
                console.log('emitted event');
            });
        })
    };

})();