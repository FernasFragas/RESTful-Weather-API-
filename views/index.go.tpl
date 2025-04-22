<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>✈🌤️TravelTab</title>

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
<body class="macos-bg">
    <div class="glass-container">
        <div class="container-fluid p-4">
            <div class="row justify-content-center">
                <img src="/traveltab.png" alt="TravelTab Logo" style="  width: 15%;  
                    max-height: 10%; 
                    object-fit: contain;"> 
                <!-- Search Form -->
                <div class="col-12 mb-4">
                    <div class="form-container">
                        <form action="/process-form/:CityName" method="POST">
                            <input type="text" name="city_name" placeholder="Search for City..." id="city_name" class="form-control" required>
                            <input type="submit" value="Search" class="btn btn-primary">
                        </form>
                    </div>
                </div>

                <!-- Weather Display -->
                <div class="col-12 mb-4">
                    {{ template "weather_display" . }}
                </div>

                <!-- Videos Section -->
                <div class="col-12 mb-4">
                    {{ template "video" . }}
                </div>

                <!-- Hotels Section -->
                <div class="col-12 mb-4">
                    {{ template "hotels_card" . }}
                </div>
            </div>
        </div>
    </div>

    <!-- Weather Card Background Script -->
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
                    weatherCard.style.backgroundImage = "url('/default.jpg')";
                }
            }
        });
    </script> 
</body>
</html>