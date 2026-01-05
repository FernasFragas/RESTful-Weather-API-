<input type="hidden" name="city_itenary" value="{{ .GeneralInfo.City }}">
<div class="date-row">
    <div class="row g-3">
        <div class="col-12 col-sm-6 col-md-4">
            <div class="form-group">
                <label class="form-label">Start Date</label>
                <input type="date" name="start_date" class="form-control" required>
            </div>
        </div>
        <div class="col-12 col-sm-6 col-md-4">
            <div class="form-group">
                <label class="form-label">End Date</label>
                <input type="date" name="end_date" class="form-control" required>
            </div>
        </div>
        <div class="col-12 col-md-4">
            <div class="form-group d-flex align-items-end">
                <button type="submit" class="btn btn-primary w-100">Generate Itinerary</button>
            </div>
    </div>
    </div>
</div>