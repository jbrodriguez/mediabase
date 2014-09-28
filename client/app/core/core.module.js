(function () {
    'use strict';

    angular.module('app.core', [
        /*
         * Angular modules
         */
        // 'ngAnimate', 'ui.router', 'ngSanitize',
        'ui.router', 'infinite-scroll',
        /*
         * Our reusable cross app code modules
         */
        'blocks.exception', 'blocks.logger', 'blocks.storage', /*'blocks.router',*/
        /*
         * 3rd Party modules
         */
        // 'ngplus'
    ]);

})();