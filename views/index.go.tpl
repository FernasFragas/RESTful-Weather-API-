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



    <!-- Google Fonts -->
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@100..900&display=swap" rel="stylesheet">

    <!-- Custom CSS -->
    <link rel="stylesheet" href="/public/styles.css">

    <!-- Flag Icons -->
    <link
            rel="stylesheet"
            href="https://cdn.jsdelivr.net/gh/lipis/flag-icons@7.0.0/css/flag-icons.min.css"
    />

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/js/bootstrap.bundle.min.js"
            integrity="sha384-kenU1KFdBIe4zVF0s0G1M5b4hcpxyD9F7jL+jjXkk+Q2h455rYXK/7HAuoJl+0I4"
            crossorigin="anonymous"></script>

    <!-- HTMX Library -->
    <script src="https://unpkg.com/htmx.org@1.9.11" integrity="sha384-0gxUXCCR8yv9FM2b+U3FDbsKthCI66oH5IA9fHppQq9DDMHuMauqq1ZHBpJxQ0J0" crossorigin="anonymous"></script>


</head>
<body>
<div class="mb-3">
    <form action="/process-form/:CityName" method="POST" class="d-flex gap-2">
        <label class="form-label">City Name</label>
            <input type="text" name="city_name" id="city_name">
            <input type="submit" value="Submit">
    </form>

    <!-- Weather Display -->
    {{ if .Weather }}
        <div class="weather-card p-4 rounded shadow text-white">
            <h2>
                {{ .Weather.Data.City }}
                <span class="fi fi-{{ .Weather.Data.Country }}"></span>
            </h2>
            <p><i class="bi bi-thermometer-half"></i> <strong>{{ .Weather.Data.Temperature }}°C</strong></p>
            <p><i class="bi bi-moisture"></i> Humidity: {{ .Weather.Data.Humidity }}%</p>
            <p><i class="bi bi-cloud"></i> {{ .Weather.Data.Condition }}</p>
        </div>
    {{ end }}

    <script>
        document.addEventListener("DOMContentLoaded", function() {
            let condition = "{{ .Weather.Data.Condition }}".toLowerCase();
            let weatherCard = document.querySelector(".weather-card");

            if (weatherCard) {
                if (condition.includes("sunny")) {
                    weatherCard.style.backgroundColor = "#FFD700"; // Gold
                } else if (condition.includes("cloudy")) {
                    weatherCard.style.backgroundColor = "#B0C4DE"; // Light Steel Blue
                } else if (condition.includes("rain")) {
                    weatherCard.style.backgroundColor = "#4682B4"; // Steel Blue
                } else if (condition.includes("storm")) {
                    weatherCard.style.backgroundColor = "#708090"; // Slate Gray
                } else {
                    weatherCard.style.backgroundColor = "#87CEEB"; // Sky Blue (Default)
                }
            }
        });
    </script>
</div>
</body>
</html>