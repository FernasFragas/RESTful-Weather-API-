{{ if .GeneralInfo }}
<div class="weather-card p-4 shadow-sm" style="width:25%; margin-right: 6px; background-size: cover; background-position: center; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
    <h2 style="font-family: 'Inter', sans-serif; font-weight: 500; color: #dddddd;">
        {{ .GeneralInfo.City }}
        <span class="fi fi-{{ .GeneralInfo.Country }}"></span>
    </h2>
    <p style="font-family: 'Inter', sans-serif; color: #dddddd;"><i class="bi bi-thermometer-half"></i> <strong>{{ .GeneralInfo.Weather.Temperature }}°C</strong></p>
    <p style="font-family: 'Inter', sans-serif; color: #dddddd;"><i class="bi bi-moisture"></i> Humidity: {{ .GeneralInfo.Weather.Humidity }}%</p>
    <p style="font-family: 'Inter', sans-serif; color: #dddddd;"><i class="bi bi-cloud"></i> {{ .GeneralInfo.Weather.Condition }}</p>
    <p style="font-family: 'Inter', sans-serif; color: #dddddd;"><i class="fas fa-water"></i> Waves Height: {{ .GeneralInfo.Waves.Height }}m</p>
</div>
{{ end }}