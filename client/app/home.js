(function () {
    'use strict';

    angular
        .module('app')
        .controller('Home', Home)

    /* @ngInject */
    function Home($state, $scope, $rootScope, api, options, logger, storage) {

        /*jshint validthis: true */
        var vm = this;

        vm.options = options;
        vm.pruneMovies = pruneMovies;

        activate();

        function activate() {
            vm.options.filterBy = storage.get('filterBy') || 'title';
            vm.options.sortBy = storage.get('sortBy') || 'added';

            return getConfig().then(function() {
                logger.info('initialized state');
            })
        };

        function getConfig() {
            return api.getConfig().then(function(data) {
                vm.options.config = data;

                if (vm.options.config.mediaPath.length === 0) {
                    $state.go('settings');
                } else {
                    $state.go('recent');
                };
            });
        };         

        $scope.$watch(angular.bind(this, function() {
            return vm.options.searchTerm;
        }), function(newVal) {
             console.log('searching for either vm.searchTerm: '+vm.options.searchTerm + ' or newVal: '+newVal);
            $state.go('search').then(function(current) {
                console.log('after search state');
                $rootScope.$emit('/local/search', newVal);
                console.log('emitted event');
            });
        });

        $scope.$watch(angular.bind(this, function() {
            return vm.options
        }), function(newVal, oldVal) {
            console.log('current: ', $state.$current.name);
            storage.set('filterBy', vm.options.filterBy);
            storage.set('sortBy', vm.options.sortBy);
        }, true);        

        function pruneMovies() {
            return api.pruneMovies().then(function(data) {
                console.log("are you ready for the fallout?");
                $state.go("recent");
            })
        }
    };

})();