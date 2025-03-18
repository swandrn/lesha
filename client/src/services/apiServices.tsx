import axios from "axios";

export const getUser = async () => {
    const token = localStorage.getItem('token');
    const response = await axios.get('http://localhost:8080/get-user', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    return response.data;
}

