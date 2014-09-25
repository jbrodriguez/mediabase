(function () {
    'use strict';

    angular
        .module('app.settings')
        .controller('Settings', Settings);

    /* @ngInject */
    function Settings($state, $q, api, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.today = "mero";

        activate();

        function activate() {
            console.log("behing petrified eyes");
        };
    }
})();