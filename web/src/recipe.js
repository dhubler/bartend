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
        return [
            util.commonStyles,
            css`
                :host {
                    display: flex;
                    margin: 1px;
                    padding: 5px;
                    flex-direction: row;
                }
                .ingredient {
                    list-style-type: none;
                }
                .top {
                    flex: 1;
                }
                .bottom {
                    flex: 2;
                    font-size: large;
                }
                .content {
                    list-style: none;
                }
            `
        ];
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
        util.openDialog(this.recipe.name, drink);
    }

    render() {
        return html`
            <div class="top">
                <h1>${this.recipe.name}</h1>
                <p>${this.recipe.description}</p>
            </div>
            <div class="bottom">
                <ul class="content">
                    ${this.recipe.ingredient.map((ingredient) => html`
                        <li class="ingredient">${ingredient.amount} oz ${ingredient.liquid}</li>
                    `)}
                </ul>
                <div class="buttons">
                    <vaadin-button theme="primary" @click=${() => this.scale(1)}>Make</vaadin-button>
                    <vaadin-button @click=${() => this.scale(0.1)}>Sample</vaadin-button>
                </div>
            </div>
        `;
    }
}
customElements.define('bartend-recipe', BartendRecipe);
