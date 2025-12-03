import { defineStore } from 'pinia';
import axios from 'axios';

const API_URL = 'http://localhost:8080/v1/auth';

export const useUserStore = defineStore('user', {
    state: () => ({
        token: localStorage.getItem('token') || '',
        user: JSON.parse(localStorage.getItem('user') || 'null'),
    }),
    getters: {
        isAuthenticated: (state) => !!state.token,
    },
    actions: {
        async login(email, password) {
            try {
                const response = await axios.post(`${API_URL}/login`, { email, password });
                this.token = response.data.token;
                this.user = { username: response.data.username, id: response.data.id };
                localStorage.setItem('token', this.token);
                localStorage.setItem('user', JSON.stringify(this.user));
                return true;
            } catch (error) {
                console.error(error);
                return false;
            }
        },
        async register(username, email, password) {
            try {
                await axios.post(`${API_URL}/register`, { username, email, password });
                return true;
            } catch (error) {
                console.error(error);
                return false;
            }
        },
        logout() {
            this.token = '';
            this.user = null;
            localStorage.removeItem('token');
            localStorage.removeItem('user');
        }
    }
});
