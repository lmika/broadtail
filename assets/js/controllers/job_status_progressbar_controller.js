import $ from 'jquery'
import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
    static targets = [
        'percentage',
        'summary'
    ];

    static values = {
        jobId: String
    }

    initialize() {
    }

    handleNewMessage(ev) {
        let message = ev.detail.message;
        if (this.jobIdValue !== message.id) {
            return;
        }

        switch (message.type) {
        case "update":
            this._updateProgress(message.percent, message.summary);
            break;
        case "newstate":
            switch (ev.state) {
            case "Done":
                this._updateProgress(100.0, "Done");
                break;
            case "Error":
                this._updateProgress(0.0, "Error");
                break;
            case "Cancelled":
                this._updateProgress(0.0, "Cancelled");
                break;
            }
        }
    }

    _updateProgress(percentage, summary) {
        $(this.percentageTarget).css({width: percentage * 2.0});
        this.summaryTargets.forEach(e => $(e).text(summary));
    }
}