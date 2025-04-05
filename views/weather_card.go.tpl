{{ if .GeneralInfo }}
<div class="weather-card">
    <div class="weather-content">
        <h2>
            {{ .GeneralInfo.City }}
            <span class="fi fi-{{ .GeneralInfo.Country }}"></span>
        </h2>
        <div class="weather-info">
            <p><i class="bi bi-thermometer-half"></i> <strong>{{ .GeneralInfo.Weather.Temperature }}°C</strong></p>
            <p><i class="bi bi-moisture"></i> Humidity: {{ .GeneralInfo.Weather.Humidity }}%</p>
            <p><i class="bi bi-cloud"></i> {{ .GeneralInfo.Weather.Condition }}</p>
            <p><i class="fas fa-water"></i> Waves Height: {{ .GeneralInfo.Waves.Height }}m</p>
        </div>
    </div>
</div>

<style>
    .weather-card {
        width: 100%;
        height: 100%;
        background-size: cover;
        background-position: center;
        border-radius: 12px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        display: flex;
        align-items: center;
        justify-content: center;
        position: relative;
        overflow: hidden;
        margin: 0;
        padding: 0;
    }

    .weather-content {
        background: rgba(0, 0, 0, 0.1);
        backdrop-filter: blur(4px);
        -webkit-backdrop-filter: blur(4px);
        padding: 1rem;
        border-radius: 8px;
        text-align: center;
        width: auto;
        min-width: 200px;
        max-width: 90%;
        border: 1px solid rgba(255, 255, 255, 0.08);
    }

    .weather-content h2 {
        font-family: 'Inter', sans-serif;
        font-weight: 500;
        color: #ffffff;
        margin-bottom: 1rem;
        margin-top: 0.5rem;
        font-size: clamp(1.2rem, 2.5vw, 1.8rem);
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
    }

    .weather-info {
        display: flex;
        flex-direction: column;
        gap: 0.4rem;
    }

    .weather-info p {
        font-family: 'Inter', sans-serif;
        color: #ffffff;
        font-size: clamp(0.9rem, 1.3vw, 1.1rem);
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.4rem;
        margin: 0;
        padding: 0.2rem 0;
    }

    .weather-info i {
        opacity: 0.8;
        font-size: 1.1em;
    }

    /* Responsive adjustments */
    @media (max-width: 768px) {
        .weather-content {
            padding: 0.75rem;
            min-width: 180px;
        }
    }

    @media (max-width: 480px) {
        .weather-content {
            padding: 0.5rem;
            min-width: 160px;
        }
        
        .weather-info {
            gap: 0.3rem;
        }
    }

    /* Add smooth transitions */
    .weather-content {
        transition: all 0.3s ease;
    }

    /* Ensure flag icon is properly sized */
    .fi {
        font-size: 1.1em;
        vertical-align: middle;
    }
</style>
{{ end }}