'use strict';

import { LitElement, html, css} from 'lit';
import { until } from 'lit-html/directives/until';
import '@vaadin/button';
import '@vaadin/tabs';
import * as util from './util';
import './recipe';
import './pump';

/**
 * `<bartend-app>` show all recipes and pumps from bartender machine.
 *  You only need one of these per application and it accepts no
 *  parameters.
*/
export class BartendApp extends LitElement {

    constructor() {
        super();
        this._recipes = [];
        this._pumps = [];
        this._page = 0;
    }

    connectedCallback() {
        this._subscribe();
        this._load().catch(console.log);
    }

    // private

    _invalidate() {
        render(this._render(), this.shadowRoot);
    }

    _elem(id) {
        return this.shadowRoot.getElementById(id);
    }

    async _load() {
        try {
            const resp = await fetch(urls().apiUrl);
            this.loading = await resp.json();
            this._recipes = this.loading.recipe;
            this._pumps = this.loading.pump;
            this._liquids = this.loading.liquids;
            this.requestUpdate();
        } catch (err) {
            util.err(err);
        }
    }

    static get styles() {
        return css`
            :host {
                display: block;
            }
            .link {
                color: #ffffff;
                text-align: center;
                font-size: large;
                margin-top: 15px;
            }
            .card {
                display: block;
                margin: 5px;
            }
            vaadin-tabs {
                background-color: #2196f3;
                text-transform: uppercase;
            }
        `;
    }

    render() {
        let content = this.loading.then(() => {
            html`
                <vaadin-tabsheet>
                    <vaadin-tabs slot="tabs" id="menu">
                        <vaadin-tab id="recipes-tab">Recipes</vaadin-tab>
                        <vaadin-tab id="setup-tab">Setup</vaadin-tab>
                    </vaadin-tabs>
                    <div class="layout vertical" id="recipes-tab">
                        ${this._recipes.map((recipe) => html`
                            <bartend-recipe recipe=${recipe}></bartend-recipe>
                        `)}
                    </div>
                    <div class="layout vertical" id="setup-tab">
                        ${this._pumps.map((pump) => html`
                            <bartend-pump pump=${pump} liquids=${this._liquids}></bartend-pump>
                        `)}
                    </div>
                </vaadin-tabsheet>
            `;
        });

        return html`${until(content, html`<span>Loading...</span>`)}`;
    }

    // feature: notify all users when anyone makes a drink.
    _subscribe() {
        let _complete = false;
        this._subscription = new EventSource(`${util.url}restconf/data/bartend:current/update`);
        this._subscription.addEventListener("message", e => {
            let msg = JSON.parse(e.data);
            if (msg.event.complete != _complete) {
                _complete = msg.event.complete;
                util.info(_complete ? 'Drink complete' : 'Drink in progress...');
            }
        });
    }
}
customElements.define('bartend-app', BartendApp);
