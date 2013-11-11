'use strict';

var FLOAT_REGEXP = /^\-?\d+((\.|\,)\d+)?$/;

String.prototype.splice = function(idx, rem, s) {
    return (this.slice(0, idx) + s + this.slice(idx + Math.abs(rem)));
};

/* Directives */
angular.module('vaultee.directives', []).
	directive('appShortversion', ['shortversion', function(version) {
    	return function(scope, elm, attrs) {
      		elm.text(version);
    	};
  	}]).
  	directive('ngBlur', function() {
  		return {
  			restrict: 'A',
  			require: 'ngModel',
  			link: function(scope, elm, attr, ngModelCtrl) {
  				if (attr.type === 'radio' || attr.type === 'checkbox') return;

  				elm.unbind('input').unbind('keydown').unbind('change');

  				elm.bind('keydown keypress', function(event) {
  					if (event.which === 13) {
  						scope.$apply(function() {
  							ngModelCtrl.$setViewValue(elm.val());
  						});
  						scope.$apply(attr.ngBlur);
  					}
  				});

  				elm.bind('blur', function() {
  					scope.$apply(function() {
  						ngModelCtrl.$setViewValue(elm.val());
  					});
					scope.$apply(attr.ngBlur);
  				});
  			}
  		};
  	}).
    directive('smartFloat', function() {
      return {
        require: 'ngModel',
        link: function(scope, elm, attrs, ctrl) {
          ctrl.$parsers.unshift(function(viewValue) {
            if (FLOAT_REGEXP.test(viewValue)) {
              ctrl.$setValidity('float', true);
              return parseFloat(viewValue.replace(',', '.'));
            } else {
              ctrl.$setValidity('float', false);
              return undefined;
            }
          });
        }
      };
    }).
    directive('currencyInput', function() {
        return {
            restrict: 'A',
            scope: {
                field: '='
            },
            replace: true,
            template: '<span style="margin: 0; padding: 0;"><input type="text" ng-model="field" class="span1" required></span>',
            link: function(scope, element, attrs) {

                $(element).bind('keyup', function(e) {
                    var input = element.find('input');
                    var inputVal = input.val();

                    //clearing left side zeros
                    while (scope.field.charAt(0) == '0') {
                        scope.field = scope.field.substr(1);
                    }

                    scope.field = scope.field.replace(/[^\d.\',']/g, '');

                    var point = scope.field.indexOf(".");
                    if (point >= 0) {
                        scope.field = scope.field.slice(0, point + 3);
                    }

                    var decimalSplit = scope.field.split(".");
                    var intPart = decimalSplit[0];
                    var decPart = decimalSplit[1];

                    intPart = intPart.replace(/[^\d]/g, '');
                    if (intPart.length > 3) {
                        var intDiv = Math.floor(intPart.length / 3);
                        while (intDiv > 0) {
                            var lastComma = intPart.indexOf(",");
                            if (lastComma < 0) {
                                lastComma = intPart.length;
                            }

                            if (lastComma - 3 > 0) {
                                intPart = intPart.splice(lastComma - 3, 0, ",");
                            }
                            intDiv--;
                        }
                    }

                    if (decPart === undefined) {
                        decPart = "";
                    }
                    else {
                        decPart = "." + decPart;
                    }
                    var res = intPart + decPart;

                    scope.$apply(function() {scope.field = res});

                });

            }
        };
    });
