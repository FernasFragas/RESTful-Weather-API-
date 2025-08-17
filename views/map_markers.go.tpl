{{ if .ItineraryItems.CoordinatesWithName }}
    {{ range .ItineraryItems.CoordinatesWithName }}
        {{ if and .CoordinatesFloat (ge (len .CoordinatesFloat) 2) }}
            L.marker([{{ index .CoordinatesFloat 0 }}, {{ index .CoordinatesFloat 1 }}])
                .addTo(map)
                .bindPopup('<strong>{{ .Name }}</strong>');
        {{ end }}
    {{ end }}
{{ end }}