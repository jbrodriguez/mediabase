(function () {
    'use strict';

    angular
        .module('app.scan')
        .controller('Scan', Scan);

    // Recent.$inject = ['$q', 'api', 'logger'];

    /* @ngInject */
    function Scan($q, api, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.context = {};

        activate();

        function activate() {
            return startScan().then(function() {
                logger.info('started scanning function');
            });
        }

        function startScan() {
            return api.startScan().then(function (data) {
                vm.context = data.data;
                return vm.context;
            });
        }

        function getStatus() {
            return api.getStatus().then(function (data) {
                vm.context = data.data;
                return vm.context;
            });
        }
    }
})();