'use strict';

/* Directives */


angular.module('myApp.directives', []).
    directive('appVersion', ['version', function(version) {
        return function(scope, elm, attrs) {
            elm.text(version);
        };
    }])
    .directive('ngEnter', function () {
        return function (scope, element, attrs) {
            element.bind("keydown keypress", function (event) {
                if(event.which === 13) {
                    scope.$apply(function (){
                        scope.$eval(attrs.ngEnter);
                    });
                    
                    event.preventDefault();
                }
            });
        };
    })
    .directive('focusMe', function ($timeout) {    
        return {    
            link: function (scope, element, attrs, model) {                
                $timeout(function () {
                    element[0].focus();
                });
            }
        };
    })
    .directive('chart', function () {
	return {
	    restrict: 'E',
	    template: '<div></div>',
	    scope: {
		chartData: "=value"
	    },
	    transclude:true,
	    replace: true,

	    link: function (scope, element, attrs) {
		var chartsDefaults = {
		    
		    chart: {
			renderTo: element[0],
			type: attrs.type || null,
			height: attrs.height || null,
			width: attrs.width || null
		    }
		};
                //Update when charts data changes
		scope.$watch('chartData', function updateChart(newModel, oldModel) {
                    if(!newModel) {
                        return;
                    }
                    // We need deep copy in order to NOT override original chart object.
                    // This allows us to override chart data member and still the keep
                    // our original renderTo will be the same
		    var deepCopy = true;
		    var newSettings = {};
		    $.extend(deepCopy, newSettings, chartsDefaults, scope.chartData);
		    var chart = new Highcharts.Chart(newSettings);
                }, true);
	    }
	}
    });
