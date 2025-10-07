{{ define "sessions/login.tpl"}}
    {{ template "layouts/header.tpl" .}}
        <h2>Login</h2>
        <form action="/login" method="POST">
            <label for="email">Email:</label><br>
            <input type="email" id="email" name="email" required><br>
            <label for="password">Password:</label><br>
            <input type="password" id="password" name="password" required><br><br>
            <input type="submit" value="Login">
        </form>
    {{ template "layouts/footer.tpl" .}}
{{end}}