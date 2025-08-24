<div class="itinerary-form-section">
    <div class="container-fluid">
        <div class="row justify-content-center">
            <div class="col-12 col-lg-10 col-xl-8">
                <h2 class="itinerary-title">Generate Your City Itinerary</h2>
                <form class="itinerary-form" 
                      hx-post="/generate-itinerary" 
                      hx-target="#itinerary-section" 
                      hx-swap="innerHTML"
                      hx-indicator="#itinerary-loading">
                    {{ template "form_dates" . }}
                    {{ template "form_categories" . }}
                </form>
                <div id="itinerary-loading" class="htmx-indicator text-center" style="display: none;">
                    <i class="fas fa-spinner fa-spin"></i> Generating itinerary...
                </div>
            </div>
        </div>
    </div>
</div>