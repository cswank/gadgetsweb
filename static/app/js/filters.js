'use strict';

angular.module('myApp.filters', [])
    .filter('interpolate', ['version', function(version) {
        return function(text) {
            return String(text).replace(/\%VERSION\%/mg, version);
        }
    }])
    .filter('countdown', [function() {
        return function(input) {
            var s = input % 60;
            s = (s < 10) ? '0' + s : s;
            var m = Math.floor(input / 60);
            m = (m < 10) ? '0' + m : m;
            var h = Math.floor(input / 3600);
            return h + ':' + m + ':' + s;
        }
    }]);
