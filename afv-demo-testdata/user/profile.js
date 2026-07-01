const api = {
    getProfile() {
        return axios.get("/profile");
    },

    saveProfile(data) {
        return axios.post("/profile", data);
    }
};
