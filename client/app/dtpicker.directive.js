(function () {
    'use strict';

    angular
        .module('app')
        .directive('dtpicker', DtPicker);

    function DtPicker($parse) {
        return function(scope, element, attrs, controller) {
            var ngModel = $parse(attrs.ngModel);
            $(function() {
                element.datetimepicker({
                    timepicker: false,
                    format: 'm/d/Y',
                    dayOfWeekStart: 1,
                    closeOnDateSelect: true,
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