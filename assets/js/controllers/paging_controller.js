import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
    static targets = ["form", "pageInput"];
    static values = {
        currentPage: Number
    };

    connect() {
    }

    previousPage(ev) {
        ev.preventDefault();
        this._stepByPage(-1);
    }

    nextPage(ev) {
        ev.preventDefault();
        this._stepByPage(1);
    }

    _stepByPage(delta) {
        let currentPage = this.currentPageValue;
        let nextPage = currentPage + delta;

        if (nextPage < 0) {
            return;
        }

        this.pageInputTarget.value = nextPage;
        this.formTarget.requestSubmit();
    }
}