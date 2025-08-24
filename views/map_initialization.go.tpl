// Check if map already exists and destroy it
if (typeof window.itineraryMap !== 'undefined') {
    window.itineraryMap.remove();
}

// Create new map instance
window.itineraryMap = L.map('map').setView([{{ .GeneralInfo.Lat }}, {{ .GeneralInfo.Lon }}], 12);

L.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>',
    subdomains: 'abcd',
    maxZoom: 19
}).addTo(window.itineraryMap);

{{ template "map_markers" . }}