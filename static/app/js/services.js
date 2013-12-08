'use strict';

var config = {
    chart: {
        type: 'line',
        zoomType: 'x'
    },
    title: {
        text: ''
    },
    subtitle: {
        text: ''
    },
    xAxis: {
        type: 'datetime',
        dateTimeLabelFormats: { // don't display the dummy year
            month: '%e. %b',
            year: '%b'
        }
    },
    yAxis: {
        title: {
            text: ''
        },
        min: 0
    },
    tooltip: {
        formatter: function() {
            return '<b>'+ this.series.name +'</b><br/>'+
                Highcharts.dateFormat('%e. %b', this.x) +': '+ this.y +' m';
        }
    },
    plotOptions: {
        line: {
            marker: {
                enabled: false
            }
        }
    },
    series: []
};


/* Services */
// Demonstrate how to register services
// In this case it is a simple value service.
angular.module('myApp.services', [])
    .factory('socket', ['$rootScope', function($rootScope) {
        var ws;
        var subscribeCallback;
        return {
            connect: function(errorCallback) {
                if(ws) { return; }
                ws = new WebSocket("wss://gadgets.dyndns-ip.com/socket");
                ws.onopen = function() {
                    
                };
                ws.onerror = function() {
                    ws = null;
                    errorCallback();
                }
                ws.onmessage = function(message) {
                    message = JSON.parse(message.data);
                    var event = message[0];
                    var payload = JSON.parse(message[1]);
                    subscribeCallback(event, payload);
                };
            },
            send: function(message) {
                ws.send(message);
            },
            subscribe: function(callback) {
                subscribeCallback = callback;
            },
        }
    }])
    .value('version', '0.1')
    .factory('history', function($rootScope) {
        return {
            getChart: function(series) {
                config.series = series;
                return config;
            }
        }
    });
