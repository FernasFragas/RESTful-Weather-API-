<div class="itinerary-container">
    <div class="section-header">
        <h2><i class="fas fa-route"></i> Your City Itinerary</h2>
        <p class="section-subtitle">Discover the best places to visit in {{ .GeneralInfo.City }}</p>
    </div>
    
    <!-- Itinerary Cards Container -->
    <div class="hotels-container">
        {{ if .ItineraryItems }}
            {{ range .ItineraryItems }}
            <div class="hotel-card">
                <div class="hotel-header">
                    <h2>{{ .Title }}</h2>
                    {{ if .Category }}
                    <div class="rating">{{ .Category }}</div>
                    {{ end }}
                </div>

                {{ if .Images }}
                <div class="hotel-images image-scroll">
                    {{ range .Images }}
                    <img src="{{ . }}" alt="Itinerary Location Image">
                    {{ end }}
                </div>
                {{ end }}

                <div class="hotel-info">
                    {{ if .Description }}
                    <div class="info-item">
                        <i class="fas fa-info-circle"></i>
                        <span>{{ .Description }}</span>
                    </div>
                    {{ end }}
                    
                    {{ if .Address }}
                    <div class="info-item">
                        <i class="fas fa-map-marker-alt"></i>
                        <span>{{ .Address }}</span>
                    </div>
                    {{ end }}

                    {{ if .Duration }}
                    <div class="info-item">
                        <i class="fas fa-clock"></i>
                        <span>{{ .Duration }}</span>
                    </div>
                    {{ end }}

                    {{ if .Price }}
                    <div class="info-item">
                        <i class="fas fa-dollar-sign"></i>
                        <span>{{ .Price }}</span>
                    </div>
                    {{ end }}
                </div>

                {{ if .Tips }}
                <div class="hotel-reviews">
                    <h3><i class="fas fa-lightbulb"></i> Travel Tips</h3>
                    <div class="reviews-container">
                        {{ range .Tips }}
                        <div class="review">
                            <p class="review-text">{{ . }}</p>
                        </div>
                        {{ end }}
                    </div>
                </div>
                {{ end }}
            </div>
            {{ end }}
        {{ else }}
            <div class="hotel-card">
                <div class="hotel-header">
                    <h2><i class="fas fa-route"></i> Itinerary</h2>
                </div>
                <div style="text-align: center; padding: 2rem; color: #f0f0f0;">
                    <div class="itinerary-form-section">
                        <h2>Generate Your City Itinerary</h2>
                        <form class="itinerary-form">
                            <div class="date-row">
                                <div>
                                    <label>Start Date</label>
                                    <input type="date" name="start_date">
                                </div>
                                <div>
                                    <label>End Date</label>
                                    <input type="date" name="end_date">
                                </div>
                                <button type="submit">Generate</button>
                            </div>
                            <div class="categories-row">
                                <label>Select Your Categories</label>
                                <div class="categories-checkboxes">
                                    <label><input type="checkbox" checked> Culture</label>
                                    <label><input type="checkbox"> Art</label>
                                    <label><input type="checkbox"> Sports</label>
                                    <label><input type="checkbox"> Museums</label>
                                    <label><input type="checkbox" checked> Parks</label>
                                    <label><input type="checkbox" checked> Restaurants</label>
                                </div>
                            </div>
                        </form>
                        <div class="options-grid">
                            <div class="options-column">
                                <label><input type="radio" name="place1"> Some Museum</label>
                                <label><input type="radio" name="place1"> Restaurant 1</label>
                                <label><input type="radio" name="place1"> Museum 2</label>
                                <label><input type="radio" name="place1"> stadium 2</label>
                                <label><input type="radio" name="place1"> Park 1</label>
                            </div>
                            <!-- Repeat .options-column as needed -->
                        </div>
                        <div class="map-preview">
                            <!-- Replace with your map integration or a placeholder image -->
                            <img src="https://via.placeholder.com/500x250?text=Map+Preview" alt="Map Preview" style="width:100%;">
                        </div>
                    </div>
                </div>
            </div>
        {{ end }}
    </div>
</div>

<style>
.itinerary-container {
    width: 100%;
    margin: 0 auto 2rem;
    padding: 20px 0;
}

.section-header {
    text-align: center;
    margin-bottom: 1.5rem;
    color: rgba(0, 0, 0, 0.3);
}

.section-header h2 {
    font-family: 'Inter', sans-serif;
    font-weight: 600;
    color: rgba(0, 0, 0, 0.3);
    margin-bottom: 0.5rem;
    font-size: clamp(1.5rem, 2.5vw, 2rem);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
}

.section-header i {
    color: #ffd700;
    font-size: 1.2em;
}

.section-subtitle {
    font-family: 'Inter', sans-serif;
    color: rgba(0, 0, 0, 0.3);
    font-size: clamp(0.9rem, 1.3vw, 1.1rem);
    margin: 0;
    font-weight: 400;
}

/* Responsive adjustments */
@media (max-width: 768px) {
    .section-header {
        margin-bottom: 1rem;
    }
    
    .section-header h2 {
        flex-direction: column;
        gap: 0.25rem;
    }
}

@media (max-width: 480px) {
    .itinerary-container {
        padding: 15px 0;
    }
}

.itinerary-form-section {
    background: rgba(0,0,0,0.03);
    border-radius: 12px;
    padding: 1.5rem;
    margin-bottom: 1rem;
}
.itinerary-form h2 {
    font-size: 2rem;
    margin-bottom: 1rem;
}
.date-row {
    display: flex;
    gap: 1rem;
    align-items: flex-end;
    margin-bottom: 1rem;
}
.categories-row {
    margin-bottom: 1rem;
}
.categories-checkboxes label {
    margin-right: 1rem;
}
.options-grid {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
}
.options-column {
    border: 1px solid #aaa;
    border-radius: 4px;
    padding: 0.5rem;
    flex: 1;
    background: #fafafa;
}
.options-column label {
    display: block;
    color: #888;
    margin-bottom: 0.5rem;
}
.map-preview {
    margin-top: 1rem;
}
</style>