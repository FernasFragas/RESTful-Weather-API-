{{ if .ItineraryItems }}
    {{ if .ItineraryItems.CoordinatesWithName }}
        {{ range .ItineraryItems.CoordinatesWithName }}
            {{ if .CoordinatesFloat }}
                {{ if ge (len .CoordinatesFloat) 2 }}
                    L.marker([{{ index .CoordinatesFloat 0 }}, {{ index .CoordinatesFloat 1 }}])
                        .addTo(window.itineraryMap)
                        .bindPopup("<b>{{ .Name }}</b>");
                {{ end }}
            {{ end }}
        {{ end }}
    {{ end }}
{{ end }}