<!-- Flex Container for Weather Display and Iframe -->
<div class="d-flex" style="width:100%; align-items: stretch;">
    <!-- Weather Display -->
    {{ template "weather_card" . }}

<!-- Iframe for Windy -->
    <iframe 
        src="{{ .GeneralInfo.EmbedURL }}" 
        style="width:75%; background-color: #f5f5f5; border:none; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
    </iframe>
</div>