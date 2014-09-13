(function () {
    'use strict';

    var core = angular.module('app.core');

    core.config(toastrConfig);

    /* @ngInject */
    function toastrConfig(toastr){
        toastr.options.timeOut = 4000;
        toastr.options.positionClass = 'toast-bottom-right';
    }

    var config = {
        appTitle: 'Mediabase',
        version: '0.7.0'
    };

    core.value('config', config);
})();