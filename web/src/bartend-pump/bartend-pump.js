'use strict';

import '../../node_modules/@polymer/paper-button/paper-button.js';
import '../../node_modules/@polymer/paper-dialog/paper-dialog.js';
import '../../node_modules/@polymer/paper-dropdown-menu/paper-dropdown-menu.js';
import '../../node_modules/@polymer/paper-listbox/paper-listbox.js';
import '../../node_modules/@polymer/paper-item/paper-item.js';
import {html, render} from '../../node_modules/lit-html/lib/lit-extended.js';
import {wsocket, urls} from '../../bartend.js';

export class BartendPump extends HTMLElement {

    constructor() {
        super();
        this.attachShadow({mode: 'open'});
        this._pump = null;
        this._liquids = null;
    }

    // properties

    set pump(pump) {
        this._pump = pump;
        this._invalidate();
    }

    get pump() {
        return this._pump;
    }

    set liquids(liquids) {
        this._liquids = liquids;
        this._invalidate();
    }

    get liquids() {
        return this._liquids;
    }

    // private

    _invalidate() {
        if (this.pump == null || this.liquids == null) {
            return;
        }
        render(this.render(), this.shadowRoot);
    }

    _update(e) {
        if (this.pump.liquid == e.target.selected) {
            // nop, also first render, not user interaction
            return
        }
        this.pump.liquid = e.target.selected;
        fetch(new Request(urls().apiUrl + 'pump=' + this.pump.id, {
            method: `PUT`,
            body: '{"liquid":"' + this.pump.liquid + '"}',
        })).catch(console.log);
    }

    _enable(on) {
        let cmd = on ? 'on' : 'off';
        fetch(new Request(urls().apiUrl + 'pump=' + this.pump.id + '/' + cmd, {
            method: `POST`,
        })).catch(console.log);
    }

    render() {
        // re: paper-dropdown-menu no-animations
        //       disable animations, otherwise error about missing Keyframe
        return html`
            <style>
                :host {
                    display: block;
                    margin: 3px;
                    padding: 10px;
                }

                paper-item {
                    --paper-input-container-input : {
                        font-size: 20px;
                    };
                }
                .content {
                    margin: 10px;
                }
            </style>
            <h1>Pump ${this.pump.id}</h1>
            <div class="content">
                <paper-dropdown-menu id="liquid" no-animations on-iron-select=${(e) => this._update(e)}>
                    <paper-listbox slot="dropdown-content" selected="${this.pump.liquid}" attr-for-selected="value">
                        ${this.liquids.map((liquid) => html`
                            <paper-item value="${liquid}">${liquid}</paper-item>
                        `)}
                    </paper-listbox>
                </paper-dropdown-menu>
            </div>
            <div>
                <paper-button raised on-click=${() => this._enable('on')}>On</paper-button>
                <paper-button raised on-click="${() => this._enable('off')}">Off</paper-button>
            </div>
        `;
    }
}

customElements.define('bartend-pump', BartendPump);
