{{ if .Hotels }}
<div class="weather-card p-4 shadow-sm" style="width:25%; margin-right: 6px; background-size: cover; background-position: center; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
    <h2 style="font-family: 'Inter', sans-serif; font-weight: 500; color: #dddddd;">
        {{ .Hotels.HotelURL }}
    </h2>
    <p style="font-family: 'Inter', sans-serif; color: #dddddd;"> <strong>{{ .Hotels.HotelName }}</strong></p>
</div>
{{ end }}