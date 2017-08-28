'use strict';
var bartend = {};

var slash = document.location.pathname.lastIndexOf("/ui/");
var path = document.location.pathname.substring(0, slash);
bartend.wsUrl = 'ws://' + document.location.host + '/restconf/streams';
bartend.apiUrl = '/restconf/data/bartend:';

bartend.wsocket = function() {
    if (bartend._wsocket == null) {
        var driver = new WebSocket(bartend.wsUrl);
        bartend._wsocket = new notify.handler(driver);
        bartend._wsocket.done = function() {
            bartend._wsocket.close();
            bartend._wsocket = null;
        };
    }
    return bartend._wsocket
}
