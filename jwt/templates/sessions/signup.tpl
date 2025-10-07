{{ define "sessions/signup.tpl"}}
    {{ template "layouts/header.tpl" .}}
        <h2>Signup</h2>
        <form action="/signup" method="POST">
            <label for="email">Email:</label><br>
            <input type="email" id="email" name="email" required><br>
            <label for="password">Password:</label><br>
            <input type="password" id="password" name="password" required><br><br>
            <input type="submit" value="Signup">
        </form>
    {{ template "layouts/footer.tpl" .}}
{{end}}