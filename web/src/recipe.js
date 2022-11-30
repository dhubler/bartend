'use strict';

import '@vaadin/button';
import '@vaadin/dialog';
import './drink';
import { LitElement, html, css } from 'lit';
import * as util from './util';

/*
    `<bartend-recipe>` - Show an individual recipe. Launch point for actually
      making the drink.
*/
export class BartendRecipe extends LitElement {

    static get styles() {
        return css`
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
        `;
    }

    static get properties() {
        return {
            recipe: {
                type: Object
            }
        }
    }

    /**
     * make a scaled up or down version of this recipe
     *
     * @param {number} e.g. 0.1 for sample, 2 for a double
     */
    scale(multiplier) {
        let drink = document.createElement('bartend-drink');
        drink.recipe = this.recipe;
        drink.start(multiplier);
        util.openDialog(drink);
    }

    render() {
        return html`
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
                    <vaadin-button raised theme="primary" @click=${() => this.scale(1)}>Make</vaadin-button>
                    <vaadin-button raised on-tap=${() => this.scale(0.1)}>Sample</vaadin-button>
                </div>
                </div>
            </div>
          `;
    }
}
customElements.define('bartend-recipe', BartendRecipe);
