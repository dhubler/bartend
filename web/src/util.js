
import {Notification} from '@vaadin/notification';
import "@vaadin/dialog";
import { css} from 'lit';

/**
 *  various code that is shared between components
 */

export const url = "/restconf/";

// NOTE: vaadin-dialogs close on ESC (ok) and clicking outside the dialog (non-ideal).
export function openDialog(title, elem) {
    const dlg = document.createElement('vaadin-dialog');
    dlg.headerTitle = title;
    dlg.renderer = (root, _dialog) => {
        root.appendChild(elem);
    };
    elem.addEventListener("close", () => {
        dlg.opened = false;
        document.body.removeChild(dlg);
    });
    document.body.appendChild(dlg);
    dlg.opened = true;
}

export const commonStyles = css`
    .buttons {
        display: flex;
        gap: 5px;
        flex-direction: row;
    }    
`;

export function info(msg) {
    Notification.show(msg, {position: 'middle'});
}

export function err(msg) {
    Notification.show(msg, {position: 'middle', theme: 'error'});
}

let stream;

// singleton for subscribing to drink updates. saves resources as a singleton 
// but not strictly nec.
export function subscribeDrinkUpdates(listener) {
    if (stream == null) {
        stream = new EventSource(`${url}data/bartend:drink/update`);
    }
    // unwrap the IETF notification message
    const wrapped  = (e) => {
        let msg = JSON.parse(e.data);
        let event = msg['ietf-restconf:notification']['event'];
        listener(event);
    }
    stream.addEventListener("message", wrapped);
    return () => stream.removeEventListener("message", wrapped);
}