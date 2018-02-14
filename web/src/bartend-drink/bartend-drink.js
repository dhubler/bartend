'use strict';

import '../../node_modules/@polymer/paper-button/paper-button.js';
import '../../node_modules/@polymer/iron-pages/iron-pages.js';
import '../../node_modules/@polymer/paper-progress/paper-progress.js';
import '../../node_modules/@polymer/paper-dialog-behavior/paper-dialog-shared-styles.js';
import {html, render} from '../../node_modules/lit-html/lib/lit-extended.js';
import {wsocket, urls} from '../../bartend.js';

/*
    `<bartend-drink>` - As a drink is being poured, keep the user informed
        of the progress.
*/
export class BartendDrink extends HTMLElement {

    constructor() {
        super();

        // properties
        this.pouring = null;
        this.recipe = null;

        this.attachShadow({mode: 'open'});
    }

    /**
     * Start making the drink.
     *
     * @param {number} recipe multiplier. e.g. a double is 2, a sample is 0.1
     */
    start(scale) {
        this._manualStepIndex = -1;
        this._manualInstruction = 'Please turn on scale.';
        this._totalWeight = 0;
        this._page = 0;
        this._sub = null;
        fetch(urls().apiUrl + 'recipe=' + this.recipe.id + '/make', {
            method: `POST`,
            body: '{"multiplier":' + scale + '}',
        }).then(() => {
            this._loadCurrent();
        }).catch(console.log);
    }

    /**
     * Stop making the drink in progress.
     */
    stop() {
        fetch(urls().apiUrl + 'current/stop', {
            method: `POST`,
        }).then(() => {
            this._done();
        }).catch(console.log);
    }

    // private

    _invalidate() {
        render(this._render(), this.shadowRoot);
    }

    _loadCurrent() {
        fetch(urls().apiUrl + 'current').then((resp) => {
            resp.json().then((data) => {
                if ('auto' in data) {
                    this.pouring = data.auto;
                    this._subscribe();
                }
                this.manual = data.manual;
                this._invalidate();
            });
        }).catch(console.log);
    }

    _subscribe() {
        this._sub = wsocket().on('bartend-drink', 'current/update', 'bartend', (update, err) => {
            if (err != null) {
                console.log(err);
                return
            }
            var autoComplete = true;
            for (var i = 0; i < update.auto.length; i++) {
                this.pouring[i].percentComplete = update.auto[i].percentComplete;
                if (!update.auto[i].complete) {
                    autoComplete = false;
                }
            }
            this._invalidate();
            if (update.complete) {
                this._done();
            } else if (autoComplete) {
                // pause so user sees the automatic part of
                // the drink is complete
                setTimeout(this._nextAuto.bind(this), 1000);
            }
        });
    }

    _done() {
        if (this._sub != null) {
            wsocket().off(this._sub);
            this._sub = null;
        }
        this.dispatchEvent(new CustomEvent('done', {bubbles:true, composed:true}));
    }

    _nextManual() {
        this._manualStepIndex++;
        var total
        if (this._manualStepIndex >= this.manual.length) {
            this.stop();
        } else {
            this._page = 1;
            var step = this.manual[this._manualStepIndex];
            this._totalWeight += step.ingredient.weight;
            this._manualInstruction = `Pour ${step.ingredient.liquid} until \
                scale reads ${this._totalWeight}`;
            this._invalidate();
        }
    }

    _nextAuto() {
        if (typeof(this.manual) === 'undefined') {
            this._done();
        } else {
            this._page = 1;
            this._invalidate();
        }
    }

    _render() {
        return html`
            <style>
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
                paper-button {
                    background-color: #2196f3;
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
            </style>
            <iron-pages id="steps" selected$="${this._page}">
                <div id="auto">
                    <ul class="content">
                        ${this.pouring.map((item) => html`
                            <li class="ingredient">${item.ingredient.liquid}</li>
                            <li class="progress"><paper-progress value="${item.percentComplete}"></paper-progress></li>
                        `)}
                    </ul>
                    <div class="buttonBar">
                        <paper-button raised on-tap="${() => this.stop()}">Stop</paper-button>
                    </div>
                </div>
                <div id="manual">
                    <div class="content">
                        ${this._manualInstruction}
                    </div>
                    <div class="buttonBar">
                        <paper-button raised on-tap=${() => this.stop()}>Stop</paper-button>
                        <div class="buttonSpacer"></div>
                        <paper-button raised on-tap=${() => this._nextManual()}>Next</paper-button>
                    </div>
                </div>
            </iron-pages>
        `;
    }
}
customElements.define('bartend-drink', BartendDrink);
