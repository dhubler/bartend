var bartend = {};

bartend.apiUrl = '/restconf/';
bartend.wsUrl = 'ws://' + document.location.host + '/restsock/';

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
