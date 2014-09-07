(function () {
    'use strict';

    angular.module('app', [
        /*
         * Order is not important. Angular makes a
         * pass to register all of the modules listed
         * and then when app.dashboard tries to use app.data,
         * it's components are available.
         */

        /*
         * Everybody has access to these.
         * We could place these under every feature area,
         * but this is easier to maintain.
         */
        'app.core',

        /*
         * Feature areas
         */
        'app.recent',
        'app.import',
        'app.search',
    ]);

    angular
        .module('app')
        .config(function($stateProvider, $urlRouterProvider) {
            $urlRouterProvider.otherwise('/recent');
        })
        .config(['$provide', function($provide) {
          $provide.decorator('$rootScope', ['$delegate', function($delegate) {
            $delegate.constructor.prototype.$onRootScope = function(name, listener) {
              var unsubscribe = $delegate.$on(name, listener);
              this.$on('$destroy', unsubscribe);
            };

            return $delegate;
          }]);
        }])        
        .controller('Home', Home)
        .directive('dtpicker', DtPicker)
        .directive('starRating', StarRating);

    /* @ngInject */
    function Home($state, $scope, $rootScope) {

        /*jshint validthis: true */
        var vm = this;

        vm.searchTerm = '';

        $scope.$watch(angular.bind(this, function(searchTerm) {
            return vm.searchTerm;
        }), function(newVal) {
             console.log('searching for either vm.searchTerm: '+vm.searchTerm + ' or newVal: '+newVal);
            $state.go('search').then(function(current) {
                console.log('after search state');
                $rootScope.$emit('/local/search', newVal);
                console.log('emitted event');
            });
        })
    };

    function DtPicker($parse) {
        return function(scope, element, attrs, controller) {
            var ngModel = $parse(attrs.ngModel);
            $(function() {
                element.datetimepicker({
                    timepicker: false,
                    format: 'm/d/Y',
                    dayOfWeekStart: 1,
                    onChangeDateTime: function(dp, $input) {
                        scope.$apply(function(scope) {
                            ngModel.assign(scope, $input.val());
                        });
                    }
                });
            });
        }
    };

    function StarRating() {
        var directive = {};

        directive.restrict = 'AE';
        directive.scope = {
            score: '=score',
            max: '=max'
        };

        directive.templateUrl = 'app/template/rating.html';

        directive.link = function(scope, elements, attr) {
            scope.updateStars = function() {
                var idx = 0;
                scope.stars = [];
                for (idx = 0; idx < scope.max; idx += 1) {
                    scope.stars.push({
                        full. scope.score > idx
                    })
                }
            };

            scope.hover = function(/** Integer **/ idx) {
                scope.hoverIdx = idx;
            };

            scope.stopHover = function() {
                scope.hoverIdx = -1;
            };

            scope.starColor = function(/** Integer **/ idx) {
                var starClass = 'rating-normal';
                if (idx <= scope.hoverIdx) {
                    starClass = 'rating-highlight';
                };
                return starClass;
            };

            scope.starClass = function(/** Star **/ star, /** Integer **/ idx) {
                var starClass = 'fa-star-o';
                if (star.full || idx <= scope.hoverIdx) {
                    starClass = 'fa-star';
                };
                return starClass;
            };

            scope.setRating = function(idx) {
                scope.score = idx + 1;
                scope.stopHover();
            };

            scope.$watch('score', function(newValue, oldValue) {
                if (newValue !== null && newValue !== undefined) {
                    scope.updateStars();
                }
            });
        };

        return directive;
    };
})();