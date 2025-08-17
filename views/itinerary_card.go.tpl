<div class="itinerary-container" id="itinerary-section">
    {{ template "itinerary_header" . }}
    
    <div class="itinerary-content">
        {{ if .ItineraryItems }}
            {{ if .ItineraryItems.CoordinatesWithName }}
                {{ template "itinerary_map" . }}
            {{ end }}
        {{ else }}
            {{ template "itinerary_form" . }}
        {{ end }}
    </div>
</div>