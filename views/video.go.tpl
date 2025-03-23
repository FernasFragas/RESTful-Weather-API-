{{if .Videos}}
<div class="video-container">
    {{range .Videos}}
        <div class="video">
            <iframe src="https://www.youtube.com/embed/{{.VideoID}}" 
            frameborder="0" allowfullscreen></iframe>
            <p>{{.Title}}</p>
        </div>
    {{end}}
</div>
{{else if .Query}}
    <p>No results found for "{{.Query}}"</p>
{{end}}