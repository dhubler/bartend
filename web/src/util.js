
import '@vaadin/notification';

export const url = "/restconf/";

export function openDialog(title, elem) {
    let dlg = document.createElement('vaadin-dialog');
    dlg.headerTitle = title;
    dlg.renderer = (root, _dialog) => {
        root.appendChild(elem);
    };
    elem.attachListener("close", (e) => {
        dlg.opened = false;
        document.body.removeChild(dlg);
    });
    dlg.opened = true;
}

export function info(msg) {
    Notification.show(msg, {position: 'middle'});
}

export function err(msg) {
    Notification.show(msg, {position: 'middle', theme: 'error'});
}