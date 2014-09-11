(function () {
    'use strict';

    angular
        .module('app')
        .filter('truncate', Truncate);

    function Truncate($parse) {
        return function (text, length, end) {
            length = length || 10;
            end = end || '...';
 
            if (text.length <= length || text.length - end.length <= length) {
                return text;
            }
            else {
                return String(text).substring(0, length-end.length) + end;
            }
 
        };
    };    

})();