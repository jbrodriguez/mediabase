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
        vm.sortOrder = sortOrder;

        activate();

        function activate() {
            vm.options.filterBy = storage.get('filterBy') || 'title';
            vm.options.sortBy = storage.get('sortBy') || 'added';

            return getConfig().then(function() {
                logger.info('initialized state');
            })
        };

        function pruneMovies() {
            return api.pruneMovies().then(function(data) {
                console.log("are you ready for the fallout?");
                $state.go("recent");
            })
        };

        function sortOrder() {
            console.log("is there anybody out there");
            vm.options.sortOrder = vm.options.sortOrder === 'desc' ? 'asc' : 'desc';
            if ($state.$current.name === 'all') {
                $rootScope.$emit('/movies/refresh');
            } else {
                $state.go('all');
            };
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
            return vm.options.filterBy;
        }), function(newVal, oldVal) {
            console.log('current: ', $state.$current.name);
            storage.set('filterBy', vm.options.filterBy);
            $state.go('all');
        }, true);

        $scope.$watch(angular.bind(this, function() {
            return vm.options.sortBy;
        }), function(newVal, oldVal) {
            console.log('current: ', $state.$current.name);
            storage.set('sortBy', vm.options.sortBy);
            if ($state.$current.name === 'all') {
                $rootScope.$emit('/movies/refresh');
            } else {
                $state.go('all');
            };
        }, true);        

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

    };

})();