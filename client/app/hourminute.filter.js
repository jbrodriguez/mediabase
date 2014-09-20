(function () {
    'use strict';

    angular
        .module('app')
        .filter('hourMinute', HourMinute);

    function HourMinute() {
        return function (minutes) {
            var hour = Math.floor(minutes / 60);
            var minute = Math.floor(minutes % 60);

            var time = '';
            if (hour > 0) time += (hour + "h ");
            if (minute > 0) time += (minute + "m");

            return time;
        };
    };    

})();