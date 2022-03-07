import { Application } from "@hotwired/stimulus"
import $ from 'jquery'
import jqueryUjsInit from 'jquery-ujs'

import PagingController from "./controllers/paging_controller"
import FavouriteLinkController from "./controllers/favourite_link_controller"
import JobStatusController from "./controllers/job_status_controller"
import JobStatusProgressbarController from "./controllers/job_status_progressbar_controller"

window.Stimulus = Application.start();
Stimulus.register("paging", PagingController);
Stimulus.register("favourite-link", FavouriteLinkController);
Stimulus.register("job-status", JobStatusController);
Stimulus.register("job-status-progressbar", JobStatusProgressbarController);

jqueryUjsInit($, undefined);