import {notifier} from './notify.js';

let _wsocket = null;

export function urls() {
    const slash = document.location.pathname.lastIndexOf("/ui/");
    return {        
        path : document.location.pathname.substring(0, slash),
        wsUrl : 'ws://' + document.location.host + '/restconf/streams',
        apiUrl : '/restconf/data/bartend:'    
    }
}

export function wsocket() {
    if (_wsocket == null) {    
        var driver = new WebSocket(urls().wsUrl);
        _wsocket = new notifier(driver);
        _wsocket.done = function() {
            _wsocket.close();
            _wsocket = null;
        };
    }
    return _wsocket;
}
