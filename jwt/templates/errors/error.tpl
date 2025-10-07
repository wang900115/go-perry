{{ define "errors/error.tpl"}}
    {{ template "layouts/header.tpl" .}}
        <br/>
        <br/>
        <br/>
        <br/>
        <div style="color: white; background-color: red; padding: 20px; margin: 10px 0; border-radius: 5px;">
            {{ .error }}
        </div>
    {{ template "layouts/footer.tpl" .}}
{{end}}
