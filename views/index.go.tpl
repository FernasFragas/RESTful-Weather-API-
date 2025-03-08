<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Weather Service</title>

    <!-- Bootstrap CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/css/bootstrap.min.css" rel="stylesheet"
          integrity="sha384-rbsA2VBKQhggwzxH7pPCaAqO46MgnOM80zW1RWuH61DGLwZJEdK2Kadq2F9CUG65" crossorigin="anonymous">

    <!-- Favicon -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.3.0/font/bootstrap-icons.css">


    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"/>
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>

    <!-- Google Fonts -->
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@100..900&display=swap" rel="stylesheet">

    <!-- Custom CSS -->
    <link rel="stylesheet" href="/styles.css">

    <!-- Flag Icons -->
    <link
            rel="stylesheet"
            href="https://cdn.jsdelivr.net/gh/lipis/flag-icons@7.0.0/css/flag-icons.min.css"
    />

    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.4/css/all.min.css"
     rel="stylesheet"
     >

     <link rel="stylesheet" 
     href="https://cdn.jsdelivr.net/npm/bootstrap-icons/font/bootstrap-icons.css">

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/js/bootstrap.bundle.min.js"
            integrity="sha384-kenU1KFdBIe4zVF0s0G1M5b4hcpxyD9F7jL+jjXkk+Q2h455rYXK/7HAuoJl+0I4"
            crossorigin="anonymous"></script>

    <!-- HTMX Library -->
    <script src="https://unpkg.com/htmx.org@1.9.11" integrity="sha384-0gxUXCCR8yv9FM2b+U3FDbsKthCI66oH5IA9fHppQq9DDMHuMauqq1ZHBpJxQ0J0" crossorigin="anonymous"></script>

</head>
<body>
    <div class="mb-3" style="width:50%; height:250px">
        <form action="/process-form/:CityName" method="POST" class="d-flex gap-2 p-3" style="background-color: #ffffff; border-radius: 12px; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);">
            <label class="form-label" style="color: #333;">City Name</label>
            <input type="text" name="city_name" id="city_name" class="form-control" style="border-radius: 8px; border: 1px solid #ddd;">
            <input type="submit" value="Submit" class="btn btn-primary" style="border-radius: 8px; background-color: #b1b2b3; border: none; color: #fff;">
        </form>

        <!-- Flex Container for Weather Display and Iframe -->
        <div class="d-flex" style="width:100%; align-items: stretch;">
            <!-- Weather Display -->
            {{ if .GeneralInfo }}
            <div class="weather-card p-4 shadow-sm" style="width:25%; margin-right: 6px; background-size: cover; background-position: center; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); overflow-y: auto;">
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

        <!-- Iframe for Windy -->
        <iframe 
            src="{{ .GeneralInfo.EmbedURL }}" 
            style="width:75%; background-color: #f5f5f5; border:none; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
        </iframe>
    </div>

    <script>
        document.addEventListener("DOMContentLoaded", function() {
            let condition = "{{ .GeneralInfo.Weather.Condition }}".toLowerCase();
            let weatherCard = document.querySelector(".weather-card");

            if (weatherCard) {
                if (condition.includes("sun")) {
                    weatherCard.style.backgroundImage = "url('/sunny.jpg')";
                } else if (condition.includes("cloud")) {
                    weatherCard.style.backgroundImage = "url('/cloudy.jpg')";
                } else if (condition.includes("rain")) {
                    weatherCard.style.backgroundImage = "url('/rainny.jpg')";
                } else if (condition.includes("storm")) {
                    weatherCard.style.backgroundImage = "url('/stormy.jpg')";
                } else {
                    weatherCard.style.backgroundImage = "url('/default.jpg')"; // Sky Blue (Default)
                }
                }
            });
        </script>

    <!-- Map Display -->
    <div id="map" style="width:100%; height:500px; background-color: #f5f5f5; border-radius: 10px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);"></div>

    <script>
        document.addEventListener("DOMContentLoaded", function() {
            const lat = "{{ .GeneralInfo.Lat }}";
            const lon = "{{ .GeneralInfo.Lon }}";
            console.log(lat, lon);

            const map = L.map('map').setView([lat, lon], 12);
        
            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                maxZoom: 19,
                attribution: '© OpenStreetMap'
            }).addTo(map);
        
            L.marker([lat, lon]).addTo(map)
                .bindPopup("{{ .GeneralInfo.City }}")
                .openPopup();
        });
    </script>
    </div>
</body>
</html>