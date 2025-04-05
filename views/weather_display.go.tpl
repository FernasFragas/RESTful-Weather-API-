<!-- weather_display.go.tpl -->
<div class="weather-display-container">
    <!-- Weather Display -->
    <div class="weather-section">
        {{ template "weather_card" . }}
    </div>

    <!-- Iframe for Windy -->
    <div class="map-section">
        <iframe 
            src="{{ .GeneralInfo.EmbedURL }}" 
            class="windy-map">
        </iframe>
    </div>
</div>

<style>
.weather-display-container {
    width: 100%;
    margin: 0 auto 0.5rem;
    display: flex;
    flex-direction: row;
    gap: 0.25rem;
    padding: 0;
}

.weather-section {
    width: 100%;
    display: flex;
    justify-content: center;
    aspect-ratio: 16/9;
}

.weather-section > div {
    width: 100%;
    height: 100%;
    border-radius: 16px;
    padding: 1.5rem;
    margin: 0;
    display: flex;
    flex-direction: column;
    justify-content: center;
}

.map-section {
    aspect-ratio: 16/9;
    width: 100%;
    border-radius: 16px;
    overflow: hidden;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
}

.windy-map {
    width: 100%;
    height: 100%;
    border: none;
    background-color: transparent;
}

/* Responsive adjustments */
@media (min-width: 1024px) {
    .weather-display-container {
        padding: 0 2rem;
        flex-direction: row; /* Ensure side by side on large screens */
    }

    .weather-section,
    .map-section {
        aspect-ratio: 21/9;
    }
}

@media (max-width: 1023px) {
    .weather-display-container {
        flex-direction: column; /* Stack vertically on smaller screens */
        gap: 0.5rem;
    }

    .weather-section,
    .map-section {
        width: 100%;
        aspect-ratio: 16/9;
    }

    .weather-section > div {
        padding: 1.25rem;
    }
}

@media (max-width: 480px) {
    .weather-display-container {
        padding: 0 0.5rem;
        gap: 0.5rem;
    }

    .weather-section > div {
        padding: 1rem;
        border-radius: 12px;
    }

    .weather-section,
    .map-section {
        aspect-ratio: 4/3;
        border-radius: 12px;
    }
}

/* Fallback for browsers that don't support backdrop-filter */
@supports not (backdrop-filter: blur(10px)) {
    .weather-section > div,
    .map-section {
        background: rgba(255, 255, 255, 0.95);
    }
}

/* Ensure content inside weather card is responsive */
.weather-section h2 {
    font-size: clamp(1.2rem, 2vw, 1.8rem);
}

.weather-section p {
    font-size: clamp(0.9rem, 1.5vw, 1.2rem);
}

/* Add smooth transitions */
.weather-section,
.map-section {
    transition: all 0.3s ease;
}
</style>