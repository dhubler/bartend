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
        this.loading = this._load();
        this._tabIndex = 0;
    }

    connectedCallback() {
        super.connectedCallback();
        let _complete = false;

        // feature: notify all users when anyone makes a drink.
        util.subscribeDrinkUpdates((e) => {
            if (e['bartend:complete'] != _complete) {
                _complete = e['bartend:complete'];
                util.info(_complete ? 'Drink complete' : 'Drink in progress...');
            }
        });
    }

    async _load() {
        try {
            const resp = await fetch(`${util.url}data/bartend:`);
            const data = await resp.json();
            this._recipes = data['bartend:available'] || [];
            this._pumps = data['bartend:pump'];
            this._liquids = data['bartend:liquids'];
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
            vaadin-tabs {
                text-transform: uppercase;
            }
        `;
    }

    _nav(e) {
        this._tabIndex = e.detail.value;
        let items = this.shadowRoot.getElementById("tab-sheet").children;
        for (let i = 0; i < items.length; i++ ) {
            items[i].style.display = (i == this._tabIndex ? 'block' : 'none');
        }
    }

    render() {
        let content = this.loading.then(() =>
            html`
                <vaadin-tabs slot="tabs" selected="${this._tabIndex}" @selected-changed=${this._nav} id="menu">
                    <vaadin-tab id="recipes-tab">Recipes</vaadin-tab>
                    <vaadin-tab id="setup-tab">Setup</vaadin-tab>
                </vaadin-tabs>
                <div id="tab-sheet">
                    <div id="recipes-tab">
                        ${this._recipes.length == 0 ? html`No available recipes. Select different liquids at pumps.` :
                            this._recipes.map((recipe) => html`
                                <bartend-recipe .recipe=${recipe}></bartend-recipe>
                            `)}
                    </div>
                    <div id="setup-tab">
                        ${this._pumps.map((pump) => html`
                            <bartend-pump @update="${this._load}" .pump=${pump} .liquids=${this._liquids}></bartend-pump>
                        `)}
                    </div>
                </div>
            `
        );

        return html`${until(content, html`<span>Loading...</span>`)}`;
    }
}
customElements.define('bartend-app', BartendApp);
