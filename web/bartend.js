var bartend = {};

var pathRx = /(.*)\/ui\/.*/;
console.log(document.location.pathname);
var pathMatch = pathRx.exec(document.location.pathname);
var path = pathMatch[1];
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
