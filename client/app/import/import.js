(function () {
    'use strict';

    angular
        .module('app.import')
        .controller('Import', Import);

    /* @ngInject */
    function Import($scope, $state, $timeout, api, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.context = {};

        activate();

        function activate() {
            logger.info('trying to import');
            return startImport().then(function() {
                logger.info('started import function');
                update();
            });
        };

        function startImport() {
            return api.startImport().then(function (data) {
                vm.context = null;
                vm.context = data;
                return vm.context;
            });
        };

        function getStatus() {
            return api.getStatus().then(function (data) {
                vm.context = null;
                vm.context = data;
                return vm.context;
            });
        };

        function update() {
            getStatus();
            if (!vm.context.completed) {
                schedule(update, 1000);
            } else {
                $state.go('cover');
            };
        };

        function schedule(fn, delay) {
            var promise = $timeout(fn, delay);
            var deregister = $scope.$on('$destroy', function() {
                $timeout.cancel(promise);
            });
            promise.then(deregister);
        };
    }
})();