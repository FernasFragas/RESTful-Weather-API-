<div class="hotels-container">
    {{ if .Hotels }}
        {{ range .Hotels }}
        <div class="hotel-card">
            <div class="hotel-header">
                <h2>{{ .HotelName }}</h2>

                {{ if .HotelURL }}
                <a href="{{ .HotelURL }}" target="_blank" class="hotel-link">
                    <i class="fas fa-external-link-alt"></i> Visit Hotel
                </a>
                {{ end }}

                {{ if .HotelRating }}
                <div class="rating">
                    <span>{{ .HotelRating }} ★</span>
                </div>
                {{ end }}
            </div>

            <div class="hotel-images image-scroll">
                {{ range .HotelPhotos }}
                <img src="data:image/jpeg;base64,{{ . }}" alt="Hotel Image">
                {{ end }}
            </div>

            <div class="hotel-info">
                {{ if .HotelAddress }}
                <div class="info-item">
                    <i class="fas fa-map-marker-alt"></i>
                    <span>{{ .HotelAddress }}</span>
                </div>
                {{ end }}
                
                {{ if .ContactPhone }}
                <div class="info-item">
                    <i class="fas fa-phone"></i>
                    <span>{{ .ContactPhone }}</span>
                </div>
                {{ end }}
            </div>

            {{ if .HotelReviews }}
            <div class="hotel-reviews">
                <h3>Reviews</h3>
                <div class="reviews-container">
                    {{ range .HotelReviews }}
                    <div class="review">{{ . }}</div>
                    {{ end }}
                </div>
            </div>
            {{ end }}
        </div>
        {{ end }}
    {{ end }}
</div>
