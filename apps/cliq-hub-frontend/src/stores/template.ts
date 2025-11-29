import { defineStore } from 'pinia';
import axios from 'axios';
import { useUserStore } from './user';

const API_URL = 'http://localhost:8080/v1/templates';

export const useTemplateStore = defineStore('template', {
    state: () => ({
        templates: [],
        currentTemplate: null,
    }),
    actions: {
        async fetchTemplates() {
            try {
                const response = await axios.get(API_URL);
                this.templates = response.data;
            } catch (error) {
                console.error(error);
            }
        },
        async fetchTemplate(id) {
            try {
                const response = await axios.get(`${API_URL}/${id}`);
                this.currentTemplate = response.data;
            } catch (error) {
                console.error(error);
            }
        },
        async createTemplate(template) {
            const userStore = useUserStore();
            try {
                await axios.post(API_URL, template, {
                    headers: {
                        Authorization: `Bearer ${userStore.token}`
                    }
                });
                return true;
            } catch (error) {
                console.error(error);
                return false;
            }
        }
    }
});
