(function () {
    'use strict';

    angular
        .module('app.about')
        .controller('About', About);

    // About.$inject = ['$q', 'api', 'logger'];

    /* @ngInject */
    function About(logger) {

        /*jshint validthis: true */
        var vm = this;

        activate();

        function activate() {
        }
    }
})();