(function () {
    'use strict';

    angular
        .module('app.import')
        .controller('Import', Import);

    Import.$inject = ['$q', 'api', 'logger'];

    function Import($q, api, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.context = {};

        activate();

        function activate() {
            logger.info('trying to import');
            return startImport().then(function() {
                logger.info('started import function');
            });
        }

        function startImport() {
            return api.startImport().then(function (data) {
                vm.context = data;
                return vm.context;
            });
        }

        function getStatus() {
            return api.getStatus().then(function (data) {
                vm.context = data;
                return vm.context;
            });
        }
    }
})();