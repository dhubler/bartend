'use strict';

import '../../node_modules/@polymer/paper-button/paper-button.js';
import '../../node_modules/@polymer/paper-dialog/paper-dialog.js';
import '../bartend-drink/bartend-drink.js';
import {html, render} from '../../node_modules/lit-html/lib/lit-extended.js';

/*
    `<bartend-recipe>` - Show an individual recipe. Launch point for actually
      making the drink.
*/
export class BartendRecipe extends HTMLElement {

    constructor() {
        super();
        this.attachShadow({mode: 'open'});
    }

    // properties
    set recipe(r) {
        this._recipe = r;
        this._invalidate();
    }

    get recipe() {
        return this._recipe;
    }

    // pour just a sample of this recipe into the glass
    sample() {
        this.scale(0.1);
    }

    // make this recipe
    make() {
        this.scale(1);
    }

    /**
     * make a scaled up or down version of this recipe
     *
     * @param {number} e.g. 0.1 for sample, 2 for a double
     */
    scale(multiplier) {
        this._elem('drinkDialog').toggle();
        this._elem('drink').start(multiplier);
    }

    // private

    _invalidate() {
        render(this._render(), this.shadowRoot);
    }

    _elem(id) {
        return this.shadowRoot.getElementById(id);
    }

    _render() {
        const closeDrinkDialog = () => {
            this._elem('drinkDialog').close();
        };
        return html`
            <style>
                :host {
                    display: block;
                    margin: 1px;
                    padding: 5px;
                }
                .ingredient {
                    list-style-type: none;
                }
                .content {
                    display: flex;
                    flex-direction: row;
                }
                .main {
                    flex: 1;
                }
                .body {
                    flex: 2;
                    font-size: large;
                }
                paper-button.major {
                    background-color: #2196f3;
                    color: #ffffff;
                }
                .buttonBar {
                    display: flex;
                    flex-direction: row;
                }
                .buttonSpacer {
                    flex: 1;
                }
            </style>
            <div class="content">
                <div class="main">
                <h1>${this.recipe.name}</h1>
                <p>${this.recipe.description}</p>
                </div>
                <div class="body">
                <ul>
                    ${this.recipe.ingredient.map((ingredient) => html`
                        <li class="ingredient">${ingredient.amount} oz ${ingredient.liquid}</li>
                    `)}
                </ul>
                <div class="buttonBar">
                    <paper-button raised class="major" on-tap=${() => this.make()}>Make</paper-button>
                    <paper-button raised on-tap=${() => this.sample()}>Sample</paper-button>
                </div>
                </div>
            </div>
            <paper-dialog modal id="drinkDialog">
                <bartend-drink id="drink" recipe=${this.recipe} on-done=${closeDrinkDialog}></bartend-drink>
            </paper-dialog>
          `;
    }
}
customElements.define('bartend-recipe', BartendRecipe);
