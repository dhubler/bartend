'use strict';

import '@vaadin/button';
import '@vaadin/select';

import { LitElement, html, css} from 'lit';
import * as util from './util';

/**
 *  `<bartend-pump>` configures and interacts with liquid pumps
 */
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

    async _update(e) {
        const liquid = e.target.value;
        if (this.pump.liquid == liquid) {
            // nop, also first render, not user interaction
            return
        }
        try {
            await fetch(`${util.url}data/bartend:pump=${this.pump.id}`, {
                method: "PUT",
                body: JSON.stringify({'bartend:liquid': liquid}),
            });
            this.pump.liquid = liquid;
            this.dispatchEvent(new CustomEvent("update"));
        } catch (err) {
            util.err(err);
        }
    }

    async _enable(on) {
        try {
            let cmd = on ? 'on' : 'off';
            await fetch(`${util.url}data/bartend:pump=${this.pump.id}/${cmd}`, {
                method: "POST",
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
                <vaadin-select id="liquid" @change=${this._update} .value="${this.pump.liquid}" .items="${this.liquids.map((item) => {return {label:item,value:item}})}">
                </vaadin-select>
            </div>
            <div class="buttons">
                <vaadin-button @click=${() => this._enable('on')}>On</vaadin-button>
                <vaadin-button @click="${() => this._enable('off')}">Off</vaadin-button>
            </div>
        `;
    }
}

customElements.define('bartend-pump', BartendPump);
