function login(username, password) {
    validate(username, password);
    const token = getToken();

    return fetch("/login", {
        method: "POST",
        headers: {
            Authorization: token
        }
    });
}
