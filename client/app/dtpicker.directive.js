(function () {
    'use strict';

    angular
        .module('app')
        .directive('dtpicker', DtPicker);

    function DtPicker($parse) {
        return function(scope, element, attrs, controller) {
            var ngModel = $parse(attrs.ngModel);
            $(function() {
                var now = new Date();
                alert(now);
                dtNow = now.getMonth() + '/' + now.getDate() + '/' + now.getFullYear();
                alert(dtNow);
                element.datetimepicker({
                    timepicker: false,
                    format: 'm/d/Y',
                    dayOfWeekStart: 1,
                    closeOnDateSelect: true,
                    value: dtNow,
                    onChangeDateTime: function(dp, $input) {
                        scope.$apply(function(scope) {
                            ngModel.assign(scope, $input.val());
                        });
                    }
                });
            });
        }
    };    

})();