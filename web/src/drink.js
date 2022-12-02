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
        this._drinkUpdate = null;
        this.recipe = null;
    }

    static get properties() {
        return {
            recipe: {
                type: Object
            }
        }
    }

    static get styles() {
        return [
            util.commonStyles,
            css`
                :host {
                    display:block;
                    padding: 3px;
                    min-width: 300px;
                    max-width: 300px;
                    min-height: 150px;
                }
                .content {
                    padding: 20px;
                    list-style: none;
                }
            `
        ];
    }

    connectedCallback() {
        super.connectedCallback();
        this._unsubscribe = util.subscribeDrinkUpdates((event) => {
            this._drinkUpdate = event;
            if (this._drinkUpdate['bartend:complete']) {
                this._done();
            } else {
                this.requestUpdate();
            }
        });
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
        try {
            await fetch(`${util.url}data/bartend:available=${encodeURIComponent(this.recipe.name)}/make`, {
                method: "POST",
                body: JSON.stringify({'bartend:input':{multiplier:scale}}),
            });
        } catch (err) {
            util.err(err);
        }
    }

    /**
     * Stop making the drink in progress.
     */
    async stop() {
        try {
            await fetch(`${util.url}data/bartend:drink/stop`,{method:'POST'});
            this._done();
        } catch (err) {
            util.err(err);
        }
    }

    _done() {
        this.dispatchEvent(new CustomEvent('close'));
    }

    render() {
        if (this._drinkUpdate == null) {
            return html`<span>Loading...</span>`;
        }
        return html`
            <ul class="content">
                ${this._drinkUpdate['bartend:pour'].map((item) => html`
                    <li class="ingredient">${item.liquid}</li>
                    <li class="progress"><vaadin-progress-bar min="0" max="100" value="${item.percentComplete}"></vaadin-progress-bar></li>
                `)}
            </ul>
            <div class="buttons">
                <vaadin-button @click="${this.stop}">Stop</vaadin-button>
            </div>
        `;
    }
}
customElements.define('bartend-drink', BartendDrink);
