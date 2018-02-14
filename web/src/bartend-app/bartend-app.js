'use strict';

import '../../node_modules/@polymer/paper-button/paper-button.js';
import '../../node_modules/@polymer/paper-tabs/paper-tabs.js';
import '../../node_modules/@polymer/paper-tabs/paper-tab.js';
import '../../node_modules/@polymer/paper-card/paper-card.js';
import '../../node_modules/@polymer/paper-toast/paper-toast.js';
import '../../node_modules/@polymer/iron-pages/iron-pages.js';
import '../../node_modules/@polymer/iron-flex-layout/iron-flex-layout.js';
import '../bartend-recipe/bartend-recipe.js';
import '../bartend-pump/bartend-pump.js';
import {html, render} from '../../node_modules/lit-html/lib/lit-extended.js';
import {wsocket, urls} from '../../bartend.js';

/**
 * `<bartend-app>` show all recipes and pumps from bartender machine.
 *  You only need one of these per application and it accepts no
 *  parameters.
*/
export class BartendApp extends HTMLElement {

    constructor() {
        super();
        this.attachShadow({mode: 'open'});
        this._recipes = [];
        this._pumps = [];
        this._page = 0;
    }

    connectedCallback() {
        this._invalidate();
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
        const resp = await fetch(urls().apiUrl);
        const data = await resp.json();
        this._recipes = data.recipe;
        this._pumps = data.pump;
        this._liquids = data.liquids;
        this._invalidate();
    }

    _changeTab(page) {
        this._page = page;
        this._invalidate();
    }

    _render() {
        return html`
            <style>
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
                paper-tabs {
                    background-color: #2196f3;
                    text-transform: uppercase;
                }
            </style>
            <paper-tabs id="menu" selected$="${this._page}" on-selected-changed=${(e) => {
                    this._changeTab(e.target.selected);
                }}>
                <paper-tab link><a href="#recipes" class="link" tabindex="-1">Recipes</a></paper-tab>
                <paper-tab link><a href="#setup" class="link" tabindex="-1">Setup</a></paper-tab>
            </paper-tabs>
            <paper-toast id="msg"></paper-toast>
            <iron-pages selected$="${this._page}" id="pages">
                <div class="layout vertical" id="recipes">
                    ${this._recipes.map((recipe) => html`
                        <paper-card class="flex card">
                            <bartend-recipe recipe=${recipe}></bartend-recipe>
                        </paper-card>
                    `)}
                </div>
                <div class="layout vertical" id="setup">
                    ${this._pumps.map((pump) => html`
                        <paper-card class="flex card">
                            <bartend-pump pump=${pump} liquids=${this._liquids}></bartend-pump>
                        </paper-card>
                    `)}
                </div>
            </iron-pages>`;
    }

    // watch bartender for drink updates, potentially from other users logged into
    // bartender app
    _subscribe() {
        wsocket().on('bartend-app', 'current/update', 'bartend', (update, err) => {
            if (err != null) {
                console.log(err);
                return
            }
            if (update == null) {
                return
            }
            if (update.complete != this.complete) {
                let msg = this._elem('msg');
                if (update.complete) {
                    msg.text = 'Drink complete';
                } else {
                    msg.text = 'Drink in progress...';
                }
                this.complete = update.complete;
                msg.open();
            }
        });
    }
}
customElements.define('bartend-app', BartendApp);
