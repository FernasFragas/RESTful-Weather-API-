<div class="itinerary-form-section">
    <h2>Generate Your City Itinerary</h2>
    <form class="itinerary-form" 
          hx-post="/generate-itinerary" 
          hx-target="#itinerary-section" 
          hx-swap="innerHTML"
          hx-indicator="#itinerary-loading">
        {{ template "form_dates" . }}
        {{ template "form_categories" . }}
    </form>
    <div id="itinerary-loading" class="htmx-indicator" style="display: none;">
        <i class="fas fa-spinner fa-spin"></i> Generating itinerary...
    </div>
</div>