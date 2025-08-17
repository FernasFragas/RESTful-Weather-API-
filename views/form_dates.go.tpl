<input type="hidden" name="city_itenary" value="{{ .GeneralInfo.City }}">
<div class="date-row">
    <div>
        <label>Start Date</label>
        <input type="date" name="start_date" required>
    </div>
    <div>
        <label>End Date</label>
        <input type="date" name="end_date" required>
    </div>
    <button type="submit">Generate</button>
</div>