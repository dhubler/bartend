'use strict';

import { LitElement, html, css} from 'lit';
import { until } from 'lit-html/directives/until';
import '@vaadin/button';
import '@vaadin/progress-bar';
import * as util from './util';

/*
    `<bartend-drink>` - As a drink is being poured, keep the user informed
        of the progress.
*/
export class BartendDrink extends LitElement {

    constructor() {
        super();

        // properties
        this.pouring = null;
        this.recipe = null;
    }

    static get properties() {
        return {
            recipe: {
                type: Object
            },
            pouring: {
                type: Object
            }
        }
    }

    static get styles() {
        return css`
            :host {
                display:block;
                padding: 3px;
                min-width: 300px;
                max-width: 300px;
                min-height: 150px;
            }
            .content {
                padding: 20px;
            }
            .buttonBar {
                display: flex;
                flex-direction: row;
            }
            .ingredient, .progress {
                list-style-type: none;
            }
            .buttonSpacer {
                flex: 1;
            }
        `
    }

    connectedCallback() {
        super.connectedCallback();
        this.loading = this._loadCurrent();
    }

    _subscribe() {
        this._subscription = new EventSource(`${util.url}restconf/data/bartend:current/update`);
        this._subscription.addEventListener("message", e => {
            let msg = JSON.parse(e.data);
            let allComplete = true;
            msg.event.auto.map((update, i) => {
                this.pouring[i].percentComplete = update.percentComplete;
                if (!update.complete) {
                    allComplete = false;
                }
            });
            if (allComplete) {
                this._done();
            }
        });        
    }

    _unsubscribe() {
        if (this._subscription != null) {
            this._subscription.close();
            this._subscription = null;
        }
    }

    disconnectedCallback() {
        this._unsubscribe();
    }

    /**
     * Start making the drink.
     *
     * @param {number} recipe multiplier. e.g. a double is 2, a sample is 0.1
     */
    async start(scale) {
        this._unsubscribe();
        try {
            await fetch(`${util.url}/data/bartend:recipe=${this.recipe.id}/make`, {
                method: "POST",
                body: JSON.stringify({multiplier:scale}),
            });
            await this._loadCurrent();
        } catch (err) {
            util.err(err);
        }
    }

    /**
     * Stop making the drink in progress.
     */
    async stop() {
        try {
            await fetch(`${util.url}/data/bartend:current/stop`);
            this._done();
        } catch (err) {
            util.err(err);
        }
    }

    async _loadCurrent() {
        try {
            let resp = await fetch(`${util.url}/data/bartend:current`);
            let data = await resp.json();
            this.pouring = data.auto;
            this._subscribe();
            this.requestUpdate();
        } catch (err) {
            util.err(err);
        }
    }

    _done() {
        this.dispatchEvent(new CustomEvent('close'));
    }

    render() {
        let content = this.loading.then(() => html`
            <ul class="content">
                ${this.pouring.map((item) => html`
                    <li class="ingredient">${item.ingredient.liquid}</li>
                    <li class="progress"><vaadin-progress-bar value="${item.percentComplete}"></vaadin-progress-bar></li>
                `)}
            </ul>
            <div class="buttonBar">
                <vaadin-button @click="${this.stop}">Stop</vaadin-button>
            </div>
        `);
        return html`${until(content, html`<span>Loading...</span>`)}`;
    }
}
customElements.define('bartend-drink', BartendDrink);
