<div class="hotels-container">
    {{ if .Hotels }}
        {{ range .Hotels }}
        <div class="hotel-card">
            <!-- Hotel Header Section -->
            <div class="hotel-header">
                <h2>{{ .HotelName }}</h2>
                {{ if .HotelRating }}
                <div class="rating">
                    <span>{{ .HotelRating }} ★</span>
                </div>
                {{ end }}
            </div>

            <!-- Hotel Images Section -->
            <div class="hotel-images">
                {{ range .HotelPhotos }}
                <img src="data:image/jpeg;base64,{{ . }}" alt="Hotel Image">
                {{ end }}
            </div>

            <!-- Hotel Info Section -->
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

            <!-- Hotel Reviews Section -->
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

            <!-- Hotel URL -->
            {{ if .HotelURL }}
            <a href="{{ .HotelURL }}" target="_blank" class="hotel-link">
                <i class="fas fa-external-link-alt"></i>
                Visit Website
            </a>
            {{ end }}
        </div>
        {{ end }}
    {{ end }}
</div>

<style>
.hotels-container {
    height: calc(100vh - 200px); /* Adjust based on your header/footer height */
    overflow-y: auto;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    margin: 0 auto;
    width: 100%;
    max-width: 1200px;
}

.hotel-card {
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    border-radius: 16px;
    padding: 2rem;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    border: 1px solid rgba(255, 255, 255, 0.1);
    min-height: fit-content;
}

.hotel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 1rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.hotel-header h2 {
    font-size: 1.25rem;
    font-weight: 500;
    color: white;
    margin: 0;
}

.rating {
    background: rgba(255, 255, 255, 0.15);
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    font-size: 0.9rem;
    color: white;
}

.hotel-images {
    width: 100%;
    overflow-x: auto;
    white-space: nowrap;
    padding: 1rem 0;
    margin: 0 -0.5rem;
}

.hotel-images img {
    display: inline-block;
    width: 250px;
    height: 180px;
    object-fit: cover;
    border-radius: 12px;
    margin: 0 0.5rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.hotel-info {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
    padding: 1rem 0;
}

.info-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.9rem;
}

.info-item i {
    opacity: 0.7;
}

.hotel-reviews {
    max-height: 300px;
    overflow-y: auto;
}

.hotel-reviews h3 {
    font-size: 1rem;
    color: white;
    margin: 0;
}

.reviews-container {
    display: grid;
    gap: 1rem;
    padding: 1rem 0;
}

.review {
    background: rgba(255, 255, 255, 0.05);
    padding: 1rem;
    border-radius: 12px;
    line-height: 1.6;
}

.hotel-link {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    color: white;
    text-decoration: none;
    font-size: 0.9rem;
    padding: 0.5rem 1rem;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    width: fit-content;
    transition: all 0.2s ease;
}

.hotel-link:hover {
    background: rgba(255, 255, 255, 0.15);
    transform: translateY(-1px);
}

/* Scrollbar styling */
.hotels-container::-webkit-scrollbar,
.hotel-reviews::-webkit-scrollbar,
.hotel-images::-webkit-scrollbar {
    width: 6px;
    height: 6px;
}

.hotels-container::-webkit-scrollbar-track,
.hotel-reviews::-webkit-scrollbar-track,
.hotel-images::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 3px;
}

.hotels-container::-webkit-scrollbar-thumb,
.hotel-reviews::-webkit-scrollbar-thumb,
.hotel-images::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 3px;
}

/* Smooth scrolling */
.hotels-container,
.hotel-reviews,
.hotel-images {
    scroll-behavior: smooth;
    -webkit-overflow-scrolling: touch;
}

/* Responsive adjustments */
@media (max-width: 768px) {
    .hotels-container {
        height: calc(100vh - 150px);
        padding: 0.5rem;
    }

    .hotel-card {
        padding: 1.5rem;
        gap: 1rem;
    }

    .hotel-images img {
        width: 200px;
        height: 150px;
    }

    .hotel-info {
        grid-template-columns: 1fr;
    }
}

@media (max-width: 480px) {
    .hotels-container {
        height: calc(100vh - 100px);
    }

    .hotel-card {
        padding: 1rem;
    }

    .hotel-images img {
        width: 180px;
        height: 135px;
    }

    .hotel-header {
        flex-direction: column;
        align-items: flex-start;
        gap: 0.5rem;
    }
}

@media (max-width: 600px) {
    .hotel-card {
        width: 90%; /* Adjust for smaller screens */
        height: 150px;
    }
}

@media (max-width: 400px) {
    .hotel-card {
        width: 100%; /* Adjust for very small screens */
        height: 100xp;
    }
}
</style>