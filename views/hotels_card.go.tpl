<div style="overflow-y: auto; max-height: 80vh;">
    {{ if .Hotels }}
        {{ range .Hotels }}
        <div class="hotel-card p-4 shadow-sm" style="width:100%; margin-right: 6px; background-size: cover; background-position: center; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
            <h2 style="font-family: 'Inter', sans-serif; font-weight: 500; color: #dddddd;">
                {{ .HotelName }}
            </h2>
            
            <!-- Display the image -->
            <div style="overflow-x: auto; white-space: nowrap;">
                {{ range .HotelPhotos }}
                <img src="data:image/jpeg;base64,{{ . }}" alt="Hotel Image" style="width: 30%; min-height: 200px; max-height: 200px; object-fit: cover; border-radius: 5px; margin-right: 10px; display: inline-block; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);">
                {{ end }}
            </div>
            
            <p style="font-family: 'Inter', sans-serif; color: #dddddd;"> <strong><a href="{{ .HotelURL }}" target="_blank">{{ .HotelURL }}</a></strong></p>
            
            <!-- Add input fields for start date and end date -->
            <!--
            <div class="date-inputs">
                <label for="start-date" style="color: #dddddd;">Start Date:</label>
                <input type="date" id="start-date" name="start-date" class="form-control mb-2">
                
                <label for="end-date" style="color: #dddddd;">End Date:</label>
                <input type="date" id="end-date" name="end-date" class="form-control">
            </div>
            -->
        </div>
        {{ end }}
    {{ end }}
</div>

<style>
    @media (max-width: 600px) {
        .hotel-card img {
            width: 45%; /* Adjust for smaller screens */
        }
    }
    @media (max-width: 400px) {
        .hotel-card img {
            width: 100%; /* Adjust for very small screens */
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