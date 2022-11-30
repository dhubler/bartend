'use strict';

import '@vaadin/button';
import '@vaadin/select';

import { LitElement, html, css} from 'lit';
import * as util from './util';

export class BartendPump extends LitElement {

    static get properties() {
        return {
            pump: {
                type: Object
            },
            liquids: {
                type: Array
            }
        }
    }

    _update(e) {
        if (this.pump.liquid == e.target.selected) {
            // nop, also first render, not user interaction
            return
        }
        this.pump.liquid = e.target.selected;
        fetch(`${urls().apiUrl}pump=${this.pump.id}`, {
            method: `PUT`,
            body: JSON.stringify({liquid: this.pump.liquid}),
        }).catch(console.log);
    }

    async _enable(on) {
        try {
            let cmd = on ? 'on' : 'off';
            await fetch(`${util.url}data/bartend:pump=${this.pump.id}/${cmd}`, {
                method: `POST`,
            });
        } catch (err) {
            util.err(err);
        }
    }

    static get styles() {
        return css`
            :host {
                display: block;
                margin: 3px;
                padding: 10px;
            }
            .content {
                margin: 10px;
            }
        `;
    }    

    render() {
        return html`
            <h1>Pump ${this.pump.id}</h1>
            <div class="content">
                <vaadin-select id="liquid" .items="${this.liquids}" @select=${this._update}>
                </vaadin-select>
            </div>
            <div>
                <vaadin-button @click=${() => this._enable('on')}>On</vaadin-button>
                <vaadin-button @click="${() => this._enable('off')}">Off</vaadin-button>
            </div>
        `;
    }
}

customElements.define('bartend-pump', BartendPump);
