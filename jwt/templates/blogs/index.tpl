{{ define "blogs/index.tpl"}}
    {{ template "layouts/header.tpl" .}}
        <script>
            function sendDelete(event, href){
                var xhttp = new XMLHttpRequest();
                event.preventDefault();
                xhttp.onreadystatechange = function() {
                    // return if not ready state -> ~4
                    if(this.readyState !== 4) {
                        return;
                    }

                    if(this.readyState === 4) {
                        // Redirect the page
                        window.location.replace(this.responseURL);
                    }
                };
                xhttp.open("DELETE", href, true);
                xhttp.send();
            }
        </script>
        
        <div>
            <h1>Blogs</h1>

            <br/>
            <br/>

            <h3>Logged in as {{ .user.Email }}</h3>
            <a class="btn btn-outline-danger" 
                href="/logout" 
                onclick="sendDelete(event, this.href)" 
                role="button">
                Logout
            </a>

            <ul>
                {{ with .blogs }}
                    {{ range . }}
                        <li>
                            <div>
                                <a href="/blogs/{{.ID}}">
                                    <h5>{{ .Title }} </h5>
                                </a>
                                <p>{{ .Content }}</p>
                            </div>
                            <br/>
                        </li>
                    {{ end }}
                {{ end }}
            </ul>
        </div>
    {{ template "layouts/footer.tpl" .}}
{{end}}