{{ if .Hotels }}
    {{ range .Hotels }}
    <div class="hotel-card p-4 shadow-sm" style="width:100%; margin-right: 6px; background-size: cover; background-position: center; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
        <h2 style="font-family: 'Inter', sans-serif; font-weight: 500; color: #dddddd;">
            {{ .HotelName }}
        </h2>
        <p style="font-family: 'Inter', sans-serif; color: #dddddd;"> <strong>{{ .HotelName }}</strong></p>
        
        <!-- Add input fields for start date and end date -->
        <div class="date-inputs">
            <label for="start-date" style="color: #dddddd;">Start Date:</label>
            <input type="date" id="start-date" name="start-date" class="form-control mb-2">
            
            <label for="end-date" style="color: #dddddd;">End Date:</label>
            <input type="date" id="end-date" name="end-date" class="form-control">
        </div>
    </div>
    {{ end }}
{{ end }}