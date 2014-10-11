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
            if (hour > 0) time += (hour + ":");
            if (minute >= 0) {
                if (minute <= 9) time += "0"+minute;
                else time += minute;
            }
            if (hour <= 0) time += "m";

            return time;
        };
    };    

})();