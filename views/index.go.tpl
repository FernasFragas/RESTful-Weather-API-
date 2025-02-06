<html lang="en">
<div class="mb-3">
    <form action="/process-form/:CityName" method="POST">
        <label class="form-label">City Name</label>
            <input type="text" name="city_name" id="city_name">
            <input type="submit" value="Submit">
    </form>
    <div>
        <tbody>
        {{.}}
        </tbody>
    </div>
</div>
</html>