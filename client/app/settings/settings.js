(function () {
    'use strict';

    angular
        .module('app.settings')
        .controller('Settings', Settings);

    /* @ngInject */
    function Settings($state, $q, api, options, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.options = options;
        vm.folder = '';
        vm.regex = '';

        vm.addFolder = addFolder;
        vm.addRegex = addRegex;

        activate();

        function activate() {
            console.log("behind petrified eyes", options);
        };

        function addFolder() {
            if (vm.folder === '') {
                return;
            };

            vm.options.config.mediaFolders.push(vm.folder);

            return api.saveConfig(vm.options.config).then(function(data) {
                logger.success('config saved succesfully');
            });
        };

        function addRegex() {
            console.log('adding regex');
        };
    }
})();